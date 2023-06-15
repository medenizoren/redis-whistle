package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// A Type represents a Value type.
type Type byte

const (
	SimpleString Type = '+'
	BulkString   Type = '$'
	Array        Type = '*'
)

// A Value represents the data of a valid RESP type.
type Value struct {
	typ   Type
	bytes []byte
	array []Value
}

// String converts Value to a string.
// If Value cannot be converted, an empty string is returned.
func (v Value) String() string {
	if v.typ == BulkString || v.typ == SimpleString {
		return string(v.bytes)
	}

	return ""
}

// Array converts Value to an array.
// If Value cannot be converted, an empty array is returned.
func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}

	return []Value{}
}

// StringArray converts Value to a string array.
// If Value cannot be converted, an empty string slice is returned.
func (v Value) StringArray() []string {
	result := make([]string, 0, len(v.array))
	for _, value := range v.Array() {
		result = append(result, value.String())
	}

	return result
}

// DecodeRESP parses a RESP message and returns a RedisValue.
func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(byteStream)
	case "$":
		return decodeBulkString(byteStream)
	case "*":
		return decodeArray(byteStream)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

// decodeSimpleString parses a simple string and returns a RedisValue.
func decodeSimpleString(byteStream *bufio.Reader) (Value, error) {
	readBytes, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:   SimpleString,
		bytes: readBytes,
	}, nil
}

// decodeBulkString parses a bulk string and returns a RedisValue.
func decodeBulkString(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %w", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %w", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(byteStream, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %w", err)
	}

	return Value{
		typ:   BulkString,
		bytes: readBytes[:count],
	}, nil
}

// decodeArray parses an array and returns a RedisValue.
func decodeArray(byteStream *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(byteStream)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %w", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %w", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(byteStream)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

// readUntilCRLF reads bytes from a byte stream until it encounters a CRLF.
func readUntilCRLF(byteStream *bufio.Reader) ([]byte, error) {
	readBytes := []byte{}

	for {
		b, err := byteStream.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		readBytes = append(readBytes, b...)
		if len(readBytes) >= 2 && readBytes[len(readBytes)-2] == '\r' {
			break
		}
	}

	return readBytes[:len(readBytes)-2], nil
}

// returnError returns an RESP error string.
func returnError(s string) string {
	return "-ERR " + s + "\r\n"
}

// returnSimpleString returns a RESP simple string.
func returnSimpleString(s string) string {
	return "+" + s + "\r\n"
}

// returnNullBulkString returns a RESP null bulk string.
func returnNullBulkString() string {
	return "$-1\r\n"
}

// returnBulkString returns a RESP bulk string.
func returnBulkString(s string) string {
	return "$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n"
}

// returnInteger returns a RESP integer.
func returnInteger(i int) string {
	return ":" + strconv.Itoa(i) + "\r\n"
}

// returnArray returns a RESP array.
func returnArray(a []string) string {
	s := "*" + strconv.Itoa(len(a)) + "\r\n"

	for _, v := range a {
		if v == "" {
			s += returnNullBulkString()
		} else {
			s += returnBulkString(v)
		}
	}

	return s
}
