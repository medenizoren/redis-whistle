package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// A Database is a Redis database.
// It contains two maps: StringKeys and ExpireKeys.
// StringKeys stores the string values.
// ExpireKeys stores the expiration times of the keys.
// It also contains a stopSignal channel and a mutex.
// The stopSignal channel is used to stop the ExpireChecker.
type Database struct {
	id         int
	StringKeys map[string]string
	ExpireKeys map[string]time.Time
	stopSignal chan bool
	mutex      sync.RWMutex
}

// NewDatabase returns a pointer to a new database.
func NewDatabase(id int) *Database {
	return &Database{
		id:         id,
		StringKeys: make(map[string]string),
		ExpireKeys: make(map[string]time.Time),
		stopSignal: make(chan bool),
	}
}

// Init initializes the database.
// If fileName is not empty, it loads the database from the file.
// It also starts the ExpireChecker.
func (db *Database) Init(fileName string) {
	if fileName != "" {
		db.Load(fileName)
	}

	db.startExpireChecker()
}

// Flush deletes all the keys in the database.
func (db *Database) Flush() {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.StringKeys = make(map[string]string)
	db.ExpireKeys = make(map[string]time.Time)
}

// Close stops the ExpireChecker and saves the database.
func (db *Database) Close() {
	db.StopExpireChecker()
}

// Save saves the database to a file.
// The file name is "database_" + id + "_dump" + ".db".
func (db *Database) Save() {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	file, err := os.Create("database_" + strconv.Itoa(db.id) + "_dump" + ".db")
	if err != nil {
		redis.logger.Println(err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(db)
	if err != nil {
		redis.logger.Println(err)
	}
}

// Load loads the database from a file.
// The file name is "database_" + id + "_dump" + ".db".
func (db *Database) Load(fileName string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if fileName == "" {
		fileName = "database_" + strconv.Itoa(db.id) + "_dump" + ".db"
	}

	file, err := os.Open(fileName)
	if err != nil {
		redis.logger.Println(err)
		return
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(db)
	if err != nil {
		redis.logger.Println(err)
		return
	}
}

// startExpireChecker starts the ExpireChecker.
// It checks if a key has expired every second.
func (db *Database) startExpireChecker() {
	db.stopSignal = make(chan bool)

	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				db.checkAndRemoveExpiredKeys()
			case <-db.stopSignal:
				ticker.Stop()
				return
			}
		}
	}()
}

// checkAndRemoveExpiredKeys checks if a key has expired.
// If a key has expired, it removes the key.
func (db *Database) checkAndRemoveExpiredKeys() {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	for key, expireTime := range db.ExpireKeys {
		if time.Now().After(expireTime) {
			delete(db.StringKeys, key)
			delete(db.ExpireKeys, key)
		}
	}
}

// checkAndRemoveExpiredKey checks if a key has expired.
// If a key has expired, it removes the key.
func (db *Database) checkAndRemoveExpiredKey(key string) bool {
	db.mutex.RLock()
	expire, ok := db.ExpireKeys[key]
	db.mutex.RUnlock()

	if !ok {
		return false
	}

	if time.Now().After(expire) {
		db.mutex.Lock()
		delete(db.StringKeys, key)
		delete(db.ExpireKeys, key)
		db.mutex.Unlock()

		return true
	}

	return false
}

// StopExpireChecker stops the ExpireChecker.
func (db *Database) StopExpireChecker() {
	db.stopSignal <- true
}

func (db *Database) GetExpire(key string) time.Time {
	db.mutex.RLock()
	expire, ok := db.ExpireKeys[key]
	db.mutex.RUnlock()

	if !ok {
		return time.Time{}
	}

	return expire
}

// Get returns the value of the given key.
// If the key does not exist, it returns an empty string.
// If the key has expired, it returns an empty string.
func (db *Database) Get(key string) string {
	db.mutex.RLock()
	storage, ok := db.StringKeys[key]
	db.mutex.RUnlock()
	if !ok {
		return ""
	}

	if db.checkAndRemoveExpiredKey(key) {
		return ""
	}

	return storage
}

// Set sets the value of the given key.
func (db *Database) Set(key string, value string) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	db.StringKeys[key] = value
}

// Del deletes the given keys.
func (db *Database) Del(keys ...string) int {
	numberOfKeysDeleted := 0

	for _, key := range keys {
		value := db.Get(key)
		if value != "" {
			db.mutex.Lock()
			delete(db.StringKeys, key)
			db.mutex.Unlock()
			numberOfKeysDeleted++
		}
	}

	return numberOfKeysDeleted
}

// GetSet sets the value of the given key and returns the old value.
// If the key does not exist, it creates a new key.
// If the key has expired, it creates a new key.
func (db *Database) GetSet(key string, value string) string {
	storage := db.Get(key)

	if storage == "" {
		db.StringKeys[key] = value
		return ""
	}

	if db.checkAndRemoveExpiredKey(key) {
		db.Set(key, value)
		return ""
	}

	oldValue := storage
	db.Set(key, value)

	return oldValue
}

// GetDel returns the value of the given key and deletes the key.
// If the key does not exist, it returns an empty string.
// If the key has expired, it returns an empty string.
func (db *Database) GetDel(key string) string {
	storage := db.Get(key)
	if storage == "" {
		return ""
	}

	db.Del(key)

	return storage
}

// Setpx sets the value of the given key with the given milliseconds.
func (db *Database) Setpx(key string, milliseconds int, value string) {
	db.Set(key, value)
	db.mutex.Lock()
	db.ExpireKeys[key] = time.Now().Add(time.Millisecond * time.Duration(milliseconds))
	db.mutex.Unlock()
}

// MSet sets the values of the given keys.
func (db *Database) MSet(args ...string) {
	for i := 0; i < len(args); i += 2 {
		db.Set(args[i], args[i+1])
	}
}

// MSetNX sets the values of the given keys if the keys do not exist.
func (db *Database) MSetNX(args ...string) bool {
	for i := 0; i < len(args); i += 2 {
		if storage := db.Get(args[i]); storage != "" {
			return false
		}
	}

	for i := 0; i < len(args); i += 2 {
		db.Set(args[i], args[i+1])
	}

	return true
}

// MGet returns the values of the given keys.
func (db *Database) MGet(args ...string) []string {
	argsLen := len(args)
	values := make([]string, argsLen)

	for i := 0; i < argsLen; i++ {
		values[i] = db.Get(args[i])
	}

	return values
}

// Incr increments the value of the given key by 1.
// If the key does not exist, it creates a new key with the value 1.
// If value of the key is not an integer, it returns 0.
func (db *Database) Incr(key string) int {
	storage := db.Get(key)
	if storage == "" {
		db.Set(key, "1")
		return 1
	}

	if db.checkAndRemoveExpiredKey(key) {
		return 0
	}

	value, err := strconv.Atoi(storage)
	if err != nil {
		return 0
	}

	value++
	db.Set(key, strconv.Itoa(value))

	return value
}

// Incrby increments the value of the given key by the given increment.
// If the key does not exist, it creates a new key with the value increment.
// If value of the key is not an integer, it returns 0.
func (db *Database) IncrBy(key string, increment int) int {
	storage := db.Get(key)
	if storage == "" {
		db.Set(key, strconv.Itoa(increment))
		return increment
	}

	if db.checkAndRemoveExpiredKey(key) {
		return 0
	}

	value, err := strconv.Atoi(storage)
	if err != nil {
		return 0
	}

	value += increment
	db.Set(key, strconv.Itoa(value))

	return value
}

// Decr decrements the value of the given key by 1.
// If the key does not exist, it creates a new key with the value -1.
// If value of the key is not an integer, it returns 0.
func (db *Database) Decr(key string) int {
	storage := db.Get(key)
	if storage == "" {
		db.Set(key, "-1")
		return -1
	}

	if db.checkAndRemoveExpiredKey(key) {
		return 0
	}

	value, err := strconv.Atoi(storage)
	if err != nil {
		return 0
	}

	value--
	db.Set(key, strconv.Itoa(value))

	return value
}

// Decrby decrements the value of the given key by the given decrement.
// If the key does not exist, it creates a new key with the value -decrement.
// If value of the key is not an integer, it returns 0.
func (db *Database) DecrBy(key string, decrement int) int {
	storage := db.Get(key)
	if storage == "" {
		db.Set(key, strconv.Itoa(-decrement))
		return -decrement
	}

	if db.checkAndRemoveExpiredKey(key) {
		return 0
	}

	value, err := strconv.Atoi(storage)
	if err != nil {
		return 0
	}

	value -= decrement
	db.Set(key, strconv.Itoa(value))

	return value
}

// Expire sets the expire time of the given key.
// If the key does not exist, it returns false.
func (db *Database) Expire(key string, seconds int) bool {
	storage := db.Get(key)
	if storage == "" {
		return false
	}

	db.mutex.Lock()
	db.ExpireKeys[key] = time.Now().Add(time.Second * time.Duration(seconds))
	db.mutex.Unlock()

	return true
}

// TTL returns the remaining time to live of the given key.
// If the key does not exist, it returns -2.
// If the key exists but has no associated expire, it returns -1.
func (db *Database) TTL(key string) int {
	// db.mutex.Lock()
	// _, ok := db.StringKeys[key]
	// db.mutex.Unlock()
	storage := db.Get(key)
	if storage == "" {
		return -2
	}

	expire := db.GetExpire(key)
	if expire == (time.Time{}) {
		return -1
	}

	return int(time.Until(expire).Seconds())
}

// Persist removes the expire time of the given key.
// If the key exists but has no associated expire, it returns false.
// If the key does not exist, it returns false.
func (db *Database) Persist(key string) bool {
	storage := db.Get(key)
	if storage == "" {
		return false
	}

	expire := db.GetExpire(key)
	if expire == (time.Time{}) {
		return false
	}

	db.mutex.Lock()
	delete(db.ExpireKeys, key)
	db.mutex.Unlock()

	return true
}

// Exists returns true if the given key exists.
func (db *Database) Exists(key ...string) int {
	numberOfKeysExisting := 0

	for _, key := range key {
		storage := db.Get(key)
		if storage != "" {
			numberOfKeysExisting++
		}
	}

	return numberOfKeysExisting
}

// Keys returns all keys matching the given pattern.
func (db *Database) Keys(pattern string) []string {
	keys := make([]string, 0, len(db.StringKeys))

	db.mutex.RLock()
	for key := range db.StringKeys {
		match, _ := filepath.Match(pattern, key)
		if match {
			keys = append(keys, key)
		}
	}
	db.mutex.RUnlock()

	return keys
}
