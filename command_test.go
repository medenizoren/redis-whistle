package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	// Initialize database
	redis = &RedisServer{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		config: &config{},
	}

	redis.Init()
}

const (
	okReply   = "+OK\r\n"
	nullReply = "$-1\r\n"
	zeroReply = ":0\r\n"
	oneReply  = ":1\r\n"
)

func TestPingCommand(t *testing.T) {
	// Test with no arguments
	result := pingCommand([]string{})
	if result != "+PONG\r\n" {
		t.Errorf("pingCommand([]string{}) = %s; want +PONG\\r\\n", result)
	}

	// Test with one argument
	result = pingCommand([]string{"hello"})
	if result != "$5\r\nhello\r\n" {
		t.Errorf("pingCommand([]string{\"hello\"}) = %s; want $5\\r\\nhello\\r\\n", result)
	}
}

func TestEchoCommand(t *testing.T) {
	// Test with one argument
	result := echoCommand([]string{"hello"})
	if result != "$5\r\nhello\r\n" {
		t.Errorf("echoCommand([]string{\"hello\"}) = %s; want $5\\r\\nhello\\r\\n", result)
	}
}

func TestSetCommand(t *testing.T) {
	// Test with two arguments
	result := setCommand([]string{"key", "value"})
	if result != okReply {
		t.Errorf("setCommand([]string{\"key\", \"value\"}) = %s; want +OK\\r\\n", result)
	}

	// Test with three arguments and PX option
	result = setCommand([]string{"key", "value", "PX", "1000"})
	if result != okReply {
		t.Errorf("setCommand([]string{\"key\", \"value\", \"PX\", \"1000\"}) = %s; want +OK\\r\\n", result)
	}

	// Test with three arguments and EX option
	result = setCommand([]string{"key", "value", "EX", "1"})
	if result != okReply {
		t.Errorf("setCommand([]string{\"key\", \"value\", \"EX\", \"1\"}) = %s; want +OK\\r\\n", result)
	}

	// Test with three arguments and unknown option
	result = setCommand([]string{"key", "value", "FOO", "1"})
	if result != "-ERR unknown command 'FOO'\r\n" {
		t.Errorf("setCommand([]string{\"key\", \"value\", \"FOO\", \"1\"}) = %s; want -ERR unknown command 'FOO'\\r\\n", result)
	}
}

func TestSetexCommand(t *testing.T) {
	// Test with two arguments
	result := setexCommand([]string{"key", "1", "value"})
	if result != okReply {
		t.Errorf("setexCommand([]string{\"key\", \"1\", \"value\"}) = %s; want +OK\\r\\n", result)
	}
}

func TestGetCommand(t *testing.T) {
	// Test with existing key

	redis.databases[redis.selectedDB].Set("key", "value")
	result := getCommand([]string{"key"})
	if result != "$5\r\nvalue\r\n" {
		t.Errorf("getCommand([]string{\"key\"}) = %s; want $5\\r\\nvalue\\r\\n", result)
	}

	// Test with non-existing key
	result = getCommand([]string{"non-existing-key"})
	if result != nullReply {
		t.Errorf("getCommand([]string{\"non-existing-key\"}) = %s; want $-1\\r\\n", result)
	}
}

func TestGetSetCommand(t *testing.T) {
	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "value")
	result := getsetCommand([]string{"key", "new-value"})
	if result != "$5\r\nvalue\r\n" {
		t.Errorf("getsetCommand([]string{\"key\", \"new-value\"}) = %s; want $5\\r\\nvalue\\r\\n", result)
	}

	// Test with non-existing key
	result = getsetCommand([]string{"non-existing-key", "value"})
	if result != nullReply {
		t.Errorf("getsetCommand([]string{\"non-existing-key\", \"value\"}) = %s; want $-1\\r\\n", result)
	}
}

func TestGetDelCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "value")
	result := getdelCommand([]string{"key"})
	if result != "$5\r\nvalue\r\n" {
		t.Errorf("getdelCommand([]string{\"key\"}) = %s; want $5\\r\\nvalue\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key") != "" {
		t.Errorf("database.Get(\"key\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key"))
	}

	// Test with non-existing key
	result = getdelCommand([]string{"non-existing-key"})
	if result != nullReply {
		t.Errorf("getdelCommand([]string{\"non-existing-key\"}) = %s; want $-1\\r\\n", result)
	}
}

func TestMsetCommand(t *testing.T) {
	// Test with even number of arguments
	result := msetCommand([]string{"key1", "value1", "key2"})
	if result != "-ERR wrong number of arguments for 'MSET' command\r\n" {
		t.Errorf("msetCommand([]string{\"key1\", \"value1\", \"key2\"}) = %s; want -ERR wrong number of arguments for 'MSET' command\\r\\n", result)
	}

	// Test with odd number of arguments
	result = msetCommand([]string{"key1", "value1", "key2", "value2"})
	if result != okReply {
		t.Errorf("msetCommand([]string{\"key1\", \"value1\", \"key2\", \"value2\"}) = %s; want +OK\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key1") != "value1" {
		t.Errorf("database.Get(\"key1\") = %s; want \"value1\"", redis.databases[redis.selectedDB].Get("key1"))
	}
	if redis.databases[redis.selectedDB].Get("key2") != "value2" {
		t.Errorf("database.Get(\"key2\") = %s; want \"value2\"", redis.databases[redis.selectedDB].Get("key2"))
	}
}

func TestMsetnxCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with even number of arguments
	result := msetnxCommand([]string{"key1", "value1", "key2"})
	if result != "-ERR wrong number of arguments for 'MSETNX' command\r\n" {
		t.Errorf("msetnxCommand([]string{\"key1\", \"value1\", \"key2\"}) = %s; want -ERR wrong number of arguments for 'MSETNX' command\\r\\n", result)
	}

	// Test with non-existing keys
	result = msetnxCommand([]string{"key1", "value1", "key2", "value2"})
	if result != oneReply {
		t.Errorf("msetnxCommand([]string{\"key1\", \"value1\", \"key2\", \"value2\"}) = %s; want :1\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key1") != "value1" {
		t.Errorf("database.Get(\"key1\") = %s; want \"value1\"", redis.databases[redis.selectedDB].Get("key1"))
	}
	if redis.databases[redis.selectedDB].Get("key2") != "value2" {
		t.Errorf("database.Get(\"key2\") = %s; want \"value2\"", redis.databases[redis.selectedDB].Get("key2"))
	}

	// Test with existing keys
	redis.databases[redis.selectedDB].Flush()
	redis.databases[redis.selectedDB].Set("key1", "value1")
	result = msetnxCommand([]string{"key1", "new-value1", "key2", "value2"})
	if result != zeroReply {
		t.Errorf("msetnxCommand([]string{\"key1\", \"new-value1\", \"key2\", \"value2\"}) = %s; want :0\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key1") != "value1" {
		t.Errorf("database.Get(\"key1\") = %s; want \"value1\"", redis.databases[redis.selectedDB].Get("key1"))
	}
	if redis.databases[redis.selectedDB].Get("key2") != "" {
		t.Errorf("database.Get(\"key2\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key2"))
	}
}

func TestMgetCommand(t *testing.T) {
	// Test with non-existing keys
	result := mgetCommand([]string{"non-existing-key1", "non-existing-key2"})
	if result != "*2\r\n$-1\r\n$-1\r\n" {
		t.Errorf("mgetCommand([]string{\"non-existing-key1\", \"non-existing-key2\"}) = %s; want *2\\r\\n$-1\\r\\n$-1\\r\\n", result)
	}

	// Test with existing keys
	redis.databases[redis.selectedDB].Set("key1", "value1")
	redis.databases[redis.selectedDB].Set("key2", "value2")
	result = mgetCommand([]string{"key1", "key2"})
	if result != "*2\r\n$6\r\nvalue1\r\n$6\r\nvalue2\r\n" {
		t.Errorf("mgetCommand([]string{\"key1\", \"key2\"}) = %s; want *2\\r\\n$6\\r\\nvalue1\\r\\n$6\\r\\nvalue2\\r\\n", result)
	}
}

func TestDelCommand(t *testing.T) {
	// Test with non-existing key
	result := delCommand([]string{"non-existing-key"})
	if result != zeroReply {
		t.Errorf("delCommand([]string{\"non-existing-key\"}) = %s; want :0\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "value")
	result = delCommand([]string{"key"})
	if result != oneReply {
		t.Errorf("delCommand([]string{\"key\"}) = %s; want :1\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key") != "" {
		t.Errorf("database.Get(\"key\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key"))
	}
}

func TestIncrCommand(t *testing.T) {
	// Test with non-existing key
	result := incrCommand([]string{"non-existing-key"})
	if result != oneReply {
		t.Errorf("incrCommand([]string{\"non-existing-key\"}) = %s; want :1\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "10")
	result = incrCommand([]string{"key"})
	if result != ":11\r\n" {
		t.Errorf("incrCommand([]string{\"key\"}) = %s; want :11\\r\\n", result)
	}
}

func TestDecrCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with non-existing key
	result := decrCommand([]string{"non-existing-key"})
	if result != ":-1\r\n" {
		t.Errorf("decrCommand([]string{\"non-existing-key\"}) = %s; want :-1\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "10")
	result = decrCommand([]string{"key"})
	if result != ":9\r\n" {
		t.Errorf("decrCommand([]string{\"key\"}) = %s; want :9\\r\\n", result)
	}
}

func TestExpireCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with non-existing key
	result := expireCommand([]string{"non-existing-key", "10"})
	if result != zeroReply {
		t.Errorf("expireCommand([]string{\"non-existing-key\", \"10\"}) = %s; want :0\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "value")
	result = expireCommand([]string{"key", "1"})
	if result != oneReply {
		t.Errorf("expireCommand([]string{\"key\", \"1\"}) = %s; want :1\\r\\n", result)
	}
	time.Sleep(2 * time.Second)
	if redis.databases[redis.selectedDB].Get("key") != "" {
		t.Errorf("database.Get(\"key\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key"))
	}
}

func TestTtlCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with non-existing key
	result := ttlCommand([]string{"non-existing-key"})
	if result != ":-2\r\n" {
		t.Errorf("ttlCommand([]string{\"non-existing-key\"}) = %s; want :-2\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Flush()
	redis.databases[redis.selectedDB].Set("key", "value")
	result = ttlCommand([]string{"key"})
	if result != ":-1\r\n" {
		t.Errorf("ttlCommand([]string{\"key\"}) = %s; want :-1\\r\\n", result)
	}
	redis.databases[redis.selectedDB].Expire("key", 1)
	time.Sleep(2 * time.Second)
	result = ttlCommand([]string{"key"})
	if result != ":-2\r\n" {
		t.Errorf("ttlCommand([]string{\"key\"}) = %s; want :-2\\r\\n", result)
	}
}

func TestPersistCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with non-existing key
	result := persistCommand([]string{"non-existing-key"})
	if result != zeroReply {
		t.Errorf("persistCommand([]string{\"non-existing-key\"}) = %s; want :0\\r\\n", result)
	}

	// Test with existing key that has no expiration
	redis.databases[redis.selectedDB].Set("key", "value")
	result = persistCommand([]string{"key"})
	if result != zeroReply {
		t.Errorf("persistCommand([]string{\"key\"}) = %s; want :0\\r\\n", result)
	}

	// Test with existing key that has expiration
	redis.databases[redis.selectedDB].Expire("key", 1)
	result = persistCommand([]string{"key"})
	if result != oneReply {
		t.Errorf("persistCommand([]string{\"key\"}) = %s; want :1\\r\\n", result)
	}
	time.Sleep(2 * time.Second)
	if redis.databases[redis.selectedDB].Get("key") == "" {
		t.Errorf("database.Get(\"key\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key"))
	}
}

func TestExistsCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with non-existing key
	result := existsCommand([]string{"non-existing-key"})
	if result != zeroReply {
		t.Errorf("existsCommand([]string{\"non-existing-key\"}) = %s; want :0\\r\\n", result)
	}

	// Test with existing key
	redis.databases[redis.selectedDB].Set("key", "value")
	result = existsCommand([]string{"key"})
	if result != oneReply {
		t.Errorf("existsCommand([]string{\"key\"}) = %s; want :1\\r\\n", result)
	}
}

func TestKeysCommand(t *testing.T) {
	redis.databases[redis.selectedDB].Flush()

	// Test with no keys
	result := keysCommand([]string{"non-existing-pattern"})
	if result != "*0\r\n" {
		t.Errorf("keysCommand([]string{\"non-existing-pattern\"}) = %s; want *0\\r\\n", result)
	}

	// Test with one key
	redis.databases[redis.selectedDB].Set("key1", "value1")
	result = keysCommand([]string{"key1"})
	if result != "*1\r\n$4\r\nkey1\r\n" {
		t.Errorf("keysCommand([]string{\"key1\"}) = %s; want *1\\r\\n$4\\r\nkey1\\r\\n", result)
	}

	// Test with multiple keys
	redis.databases[redis.selectedDB].Set("key2", "value2")
	redis.databases[redis.selectedDB].Set("key3", "value3")
	result = keysCommand([]string{"key*"})
	if result != "*3\r\n$4\r\nkey1\r\n$4\r\nkey2\r\n$4\r\nkey3\r\n" {
		t.Errorf("keysCommand([]string{\"key*\"}) = %s; want *3\\r\\n$4\\r\nkey1\\r\\n$4\\r\nkey2\\r\\n$4\\r\nkey3\\r\\n", result)
	}
}

func TestSelectCommand(t *testing.T) {
	// Test selecting an existing database
	result := selectCommand([]string{"1"})
	if result != okReply {
		t.Errorf("selectCommand([]string{\"1\"}) = %s; want +OK\\r\\n", result)
	}

	// Test selecting a database that doesn't exist
	result = selectCommand([]string{"100"})
	if result != "-ERR value is out of range or invalid DB index\r\n" {
		t.Errorf("selectCommand([]string{\"2\"}) = %s; want -ERR value is out of range or invalid DB index\\r\\n", result)
	}

	// Test selecting a database with a non-integer argument
	result = selectCommand([]string{"non-integer"})
	if result != "-ERR value is not an integer\r\n" {
		t.Errorf("selectCommand([]string{\"non-integer\"}) = %s; want -ERR value is not an integer\\r\\n", result)
	}

	// Test selecting a database with no argument
	result = selectCommand([]string{})
	if result != "-ERR wrong number of arguments for 'SELECT' command\r\n" {
		t.Errorf("selectCommand([]string{}) = %s; want -ERR wrong number of arguments for 'SELECT' command\\r\\n", result)
	}

	// Test selecting a database with multiple arguments
	result = selectCommand([]string{"1", "2"})
	if result != "-ERR wrong number of arguments for 'SELECT' command\r\n" {
		t.Errorf("selectCommand([]string{\"1\", \"2\"}) = %s; want -ERR wrong number of arguments for 'SELECT' command\\r\\n", result)
	}

	// Test selecting a database with a negative argument
	result = selectCommand([]string{"-1"})
	if result != "-ERR value is out of range or invalid DB index\r\n" {
		t.Errorf("selectCommand([]string{\"-1\"}) = %s; want -ERR value is out of range or invalid DB index\\r\\n", result)
	}

	// Test selecting a database with a zero argument
	result = selectCommand([]string{"0"})
	if result != okReply {
		t.Errorf("selectCommand([]string{\"0\"}) = %s; want +OK\\r\\n", result)
	}
}

func TestFlushDBCommand(t *testing.T) {
	// Test flushing an existing database
	redis.databases[redis.selectedDB].Set("key", "value")
	result := flushdbCommand([]string{})
	if result != okReply {
		t.Errorf("flushDBCommand([]string{}) = %s; want +OK\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key") != "" {
		t.Errorf("database.Get(\"key\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key"))
	}

	// Test flushing a non-existing database
	result = flushdbCommand([]string{})
	if result != okReply {
		t.Errorf("flushDBCommand([]string{}) = %s; want +OK\\r\\n", result)
	}
}

func TestFlushAllCommand(t *testing.T) {
	// Test flushing all databases
	redis.databases[redis.selectedDB].Set("key1", "value1")
	selectCommand([]string{"1"})
	redis.databases[redis.selectedDB].Set("key2", "value2")
	result := flushallCommand([]string{})
	if result != okReply {
		t.Errorf("flushAllCommand([]string{}) = %s; want +OK\\r\\n", result)
	}
	if redis.databases[redis.selectedDB].Get("key1") != "" {
		t.Errorf("database.Get(\"key1\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key1"))
	}
	if redis.databases[redis.selectedDB].Get("key2") != "" {
		t.Errorf("database.Get(\"key2\") = %s; want \"\"", redis.databases[redis.selectedDB].Get("key2"))
	}
}
