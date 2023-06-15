package main

import (
	"strconv"
	"strings"
)

// A CommandFunc is the type of a Redis command function.
type CommandFunc func(args []string) string

// CommandMap stores the Redis command functions.
func getCommandMap() map[string]CommandFunc {
	return map[string]CommandFunc{
		"PING":     pingCommand,
		"ECHO":     echoCommand,
		"SET":      setCommand,
		"SETEX":    setexCommand,
		"GET":      getCommand,
		"GETSET":   getsetCommand,
		"GETDEL":   getdelCommand,
		"MSET":     msetCommand,
		"MSETNX":   msetnxCommand,
		"MGET":     mgetCommand,
		"DEL":      delCommand,
		"INCR":     incrCommand,
		"INCRBY":   incrbyCommand,
		"DECR":     decrCommand,
		"DECRBY":   decrbyCommand,
		"EXPIRE":   expireCommand,
		"TTL":      ttlCommand,
		"PERSIST":  persistCommand,
		"EXISTS":   existsCommand,
		"KEYS":     keysCommand,
		"SAVE":     saveCommand,
		"LOAD":     loadCommand,
		"SELECT":   selectCommand,
		"FLUSHDB":  flushdbCommand,
		"FLUSHALL": flushallCommand,
	}
}

// checkNumberOfArguments checks if the number of arguments is as expected.
func checkNumberOfArguments(args []string, expectedNumberOfArguments int) bool {
	return len(args) >= expectedNumberOfArguments
}

// returnWrongNumberOfArgumentsError returns an error message for wrong number of arguments.
func returnWrongNumberOfArgumentsError(command string) string {
	return returnError("wrong number of arguments for '" + command + "' command")
}

// pingCommand returns PONG if called with no arguments, otherwise it returns the first argument.
func pingCommand(args []string) string {
	if len(args) > 0 && args[0] != "" {
		return returnBulkString(args[0])
	}

	return returnSimpleString("PONG")
}

// echoCommand returns the first argument.
func echoCommand(args []string) string {
	return returnBulkString(args[0])
}

// setCommand sets the value at key to value.
// If key already holds a value, it is overwritten.
// If PX or EX is specified, the value is set with the specified expiration.
func setCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("SET")
	}

	if len(args) >= 3 {
		optionCommand := args[2]

		switch strings.ToUpper(optionCommand) {
		case "PX":
			milliseconds, err := strconv.Atoi(args[3])
			if err != nil {
				return returnError("value is not an integer or out of range")
			}

			redis.databases[redis.selectedDB].Setpx(args[0], milliseconds, args[1])
		case "EX":
			seconds, err := strconv.Atoi(args[3])
			if err != nil {
				return returnError("value is not an integer or out of range")
			}

			redis.databases[redis.selectedDB].Setpx(args[0], seconds*1000, args[1])
		default:
			return returnError("unknown command '" + optionCommand + "'")
		}
	} else {
		redis.databases[redis.selectedDB].Set(args[0], args[1])
	}

	return returnSimpleString("OK")
}

// setexCommand sets the value and expiration in seconds of a key.
func setexCommand(args []string) string {
	validate := checkNumberOfArguments(args, 3)
	if !validate {
		return returnWrongNumberOfArgumentsError("SETEX")
	}

	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return returnError("value is not an integer or out of range")
	}

	redis.databases[redis.selectedDB].Setpx(args[0], seconds*1000, args[2])

	return returnSimpleString("OK")
}

// getCommand returns the value at key.
func getCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("GET")
	}

	value := redis.databases[redis.selectedDB].Get(args[0])
	if value == "" {
		return returnNullBulkString()
	}

	return returnBulkString(value)
}

// getsetCommand sets the value at key to value and returns the old value at key.
func getsetCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("GETSET")
	}

	value := redis.databases[redis.selectedDB].GetSet(args[0], args[1])
	if value == "" {
		return returnNullBulkString()
	}

	return returnBulkString(value)
}

// getdelCommand deletes the key and returns the value at key.
func getdelCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("GETDEL")
	}

	value := redis.databases[redis.selectedDB].GetDel(args[0])
	if value == "" {
		return returnNullBulkString()
	}

	return returnBulkString(value)
}

// msetCommand sets the given keys to their respective values.
func msetCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("MSET")
	}

	if len(args)%2 != 0 {
		return returnError("wrong number of arguments for 'MSET' command")
	}

	redis.databases[redis.selectedDB].MSet(args...)

	return returnSimpleString("OK")
}

// msetnxCommand sets the given keys to their respective values if none of the keys already exist.
func msetnxCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("MSETNX")
	}

	if len(args)%2 != 0 {
		return returnError("wrong number of arguments for 'MSETNX' command")
	}

	if redis.databases[redis.selectedDB].MSetNX(args...) {
		return returnInteger(1)
	}

	return returnInteger(0)
}

// mgetCommand returns the values of all specified keys.
func mgetCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("MGET")
	}

	values := redis.databases[redis.selectedDB].MGet(args...)
	return returnArray(values)
}

// delCommand deletes the specified keys and returns the number of keys deleted.
func delCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("DEL")
	}

	numberOfKeysDeleted := redis.databases[redis.selectedDB].Del(args...)
	return returnInteger(numberOfKeysDeleted)
}

// incrCommand increments the number stored at key by one.
func incrCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("INCR")
	}

	return returnInteger(redis.databases[redis.selectedDB].Incr(args[0]))
}

// incrbyCommand increments the number stored at key by increment.
func incrbyCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("INCRBY")
	}

	increment, err := strconv.Atoi(args[1])
	if err != nil {
		return returnError("value is not an integer or out of range")
	}

	return returnInteger(redis.databases[redis.selectedDB].IncrBy(args[0], increment))
}

// decrCommand decrements the number stored at key by one.
func decrCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("DECR")
	}

	return returnInteger(redis.databases[redis.selectedDB].Decr(args[0]))
}

// decrbyCommand decrements the number stored at key by decrement.
func decrbyCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("DECRBY")
	}

	decrement, err := strconv.Atoi(args[1])
	if err != nil {
		return returnError("value is not an integer or out of range")
	}

	return returnInteger(redis.databases[redis.selectedDB].DecrBy(args[0], decrement))
}

// expireCommand sets a timeout on key.
func expireCommand(args []string) string {
	validate := checkNumberOfArguments(args, 2)
	if !validate {
		return returnWrongNumberOfArgumentsError("EXPIRE")
	}

	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return returnError("value is not an integer or out of range")
	}

	if redis.databases[redis.selectedDB].Expire(args[0], seconds) {
		return returnInteger(1)
	}

	return returnInteger(0)
}

// ttlCommand returns the remaining time to live of a key that has a timeout.
func ttlCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("TTL")
	}

	seconds := redis.databases[redis.selectedDB].TTL(args[0])

	return returnInteger(seconds)
}

// persistCommand removes the existing timeout on key.
func persistCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("PERSIST")
	}

	if redis.databases[redis.selectedDB].Persist(args[0]) {
		return returnInteger(1)
	}

	return returnInteger(0)
}

// existsCommand returns if key exists.
func existsCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("EXISTS")
	}

	numberOfKeysExisting := redis.databases[redis.selectedDB].Exists(args...)

	return returnInteger(numberOfKeysExisting)
}

// keysCommand returns all keys matching pattern.
func keysCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("KEYS")
	}

	keys := redis.databases[redis.selectedDB].Keys(args[0])
	return returnArray(keys)
}

// saveCommand saves the current database on disk.
func saveCommand(_ []string) string {
	redis.databases[redis.selectedDB].Save()
	return returnSimpleString("OK")
}

// loadCommand loads the current database from disk.
func loadCommand(args []string) string {
	validate := checkNumberOfArguments(args, 1)
	if !validate {
		return returnWrongNumberOfArgumentsError("LOAD")
	}

	if len(args) > 0 {
		redis.databases[redis.selectedDB].Load(args[0])
	} else {
		redis.databases[redis.selectedDB].Load("")
	}

	return returnSimpleString("OK")
}

// selectCommand selects the database having the specified zero-based numeric index.
func selectCommand(args []string) string {
	if len(args) != 1 {
		return returnWrongNumberOfArgumentsError("SELECT")
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		return returnError("value is not an integer")
	}

	if index < 0 || index >= 16 {
		return returnError("value is out of range or invalid DB index")
	}

	redis.SelectDB(index)
	redis.logger.Println("Switched to database id:", index)

	return returnSimpleString("OK")
}

// flushdbCommand deletes all keys from the current database.
func flushdbCommand(_ []string) string {
	redis.databases[redis.selectedDB].Flush()
	return returnSimpleString("OK")
}

// flushallCommand deletes all keys from all databases.
func flushallCommand(_ []string) string {
	for _, database := range redis.databases {
		database.Flush()
	}

	return returnSimpleString("OK")
}
