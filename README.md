# RedisWhistle

RedisWhistle is a mini Redis clone developed in Go, designed to whistle away your data storage worries. With RedisWhistle, you can store, retrieve, and manipulate your data with a touch of humor and a pinch of magic.

## Download Binaries

You can download pre-built binaries for RedisWhistle for different operating systems from the GitHub Releases page. Choose the binary that matches your system and get ready to whistle!

- [Linux (x64)](https://github.com/medenizoren/redis-whistle/releases/download/v1.0.0/RedisWhistle-linux-amd64)
- [macOS](https://github.com/medenizoren/redis-whistle/releases/download/v1.0.0/RedisWhistle-darwin-amd64)
- [Windows](https://github.com/medenizoren/redis-whistle/releases/download/v1.0.0/RedisWhistle-windows-amd64.exe)

## Installation

To install RedisWhistle using the provided binaries, follow these steps:

### Linux

1. Open a terminal.

2. Download the Linux binary:

```bash
$ wget https://github.com/medenizoren/redis-whistle/releases/download/v1.0.0/RedisWhistle-linux-amd64
```

3. Make the binary executable:

```bash
$ chmod +x RedisWhistle-linux-amd64
```

4. Run RedisWhistle:

```bash
$ ./RedisWhistle-linux-amd64
```

### macOS

1. Open a terminal.

2. Download the macOS binary:

```bash
$ curl -LO https://github.com/medenizoren/redis-whistle/releases/download/v1.0.0/RedisWhistle-darwin-amd64
```

3. Make the binary executable:

```bash
$ chmod +x RedisWhistle-darwin-amd64
```

4. Run RedisWhistle:

```bash
$ ./RedisWhistle-darwin-amd64
```

### Windows

1. Download the Windows binary from the provided link.

2. Double-click the downloaded `RedisWhistle-windows-amd64.exe` file to run RedisWhistle.

That's it! RedisWhistle is now up and running on your machine. You can start using it by connecting to it using a Redis client.

## Building from Source

If you want to build RedisWhistle, follow these simple steps:

1. Clone the repository to your local machine:
   ```shell
   git clone https://github.com/medenizoren/redis-whistle.git
   ```

2. Navigate to the project directory:
   ```shell
   cd redis-whistle
   ```

3. Build the project:
   ```shell
   go build
   ```

4. Run RedisWhistle:
   ```shell
   ./redis-whistle
   ```

5. RedisWhistle is now up and running, ready to respond to your commands.

## Usage

RedisWhistle supports the following command-line flags:

- `port`: The port number for the Redis server to listen on. By default, it is set to `6379`. You can specify a different port using the `-port` flag. For example:

```bash
$ ./redis-whistle -port 8080
```

- `load`: If you have a Redis database dump file, you can use the `-load` flag to load the data into RedisWhistle. Provide the file name as the value for the `-load` flag. For example:

```bash
$ ./redis-whistle -load dump.db
```

## Supported Commands

RedisWhistle supports the following commands:

- `PING`: Test if RedisWhistle is listening. Expect a "PONG" response.

- `ECHO [message]`: Returns the message you provide. RedisWhistle will echo it back with its own twist of humor.

- `SET [key] [value]`: Set a key-value pair in the RedisWhistle store.

- `SETEX [key] [seconds] [value]`: Set a key-value pair in the RedisWhistle store with an expiration time in seconds.

- `GET [key]`: Retrieve the value associated with the given key.

- `GETSET [key] [value]`: Set the value of the key and return its previous value.

- `GETDEL [key]`: Get the value associated with the key and delete the key from RedisWhistle.

- `MSET [key1] [value1] [key2] [value2] ...`: Set multiple key-value pairs simultaneously.

- `MSETNX [key1] [value1] [key2] [value2] ...`: Set multiple key-value pairs if none of the keys exist.

- `MGET [key1] [key2] ...`: Retrieve the values associated with multiple keys.

- `DEL [key1] [key2] ...`: Delete one or more keys from RedisWhistle.

- `INCR [key]`: Increment the integer value stored at the given key by 1.

- `INCRBY [key] [increment]`: Increment the integer value stored at the given key by the provided increment.

- `DECR [key]`: Decrement the integer value stored at the given key by 1.

- `DECRBY [key] [decrement]`: Decrement the integer value stored at the given key by the provided decrement.

- `EXPIRE [key] [seconds]`: Set an expiration time (in seconds) for the given key.

- `TTL [key]`: Get the remaining time to live (in seconds) for the given key.

- `PERSIST [key]`: Remove the expiration time for the given key, making it persist.

- `EXISTS [key]`: Check if the given key exists in RedisWhistle.

- `KEYS [pattern]`: Return all the keys matching the provided pattern.

- `SAVE`: Save the current state of RedisWhistle to disk. 

- `LOAD`: Load the previously saved state of RedisWhistle. RedisWhistle never forgets, just like an elephant!

- `SELECT [database]`: Select the specified Redis Whistle database. RedisWhistle loves a good conversation, even when it involves multiple databases.

- `FLUSHDB`: Clear the currently selected database. RedisWhistle isn't afraid to start fresh when needed.

- `FLUSHALL`: Clear all databases in RedisWhistle. RedisWhistle knows how to make a clean sweep.

RedisWhistle will respond to your commands promptly and entertain you with witty replies along the way. Enjoy the RedisWhistle experience!

## Contributing
RedisWhistle is an open-source project, and we welcome contributions from the community. Feel free to create issues, submit pull requests, or suggest improvements. Together, we can make RedisWhistle an even funnier and more powerful Redis clone.

## License
RedisWhistle is released under the [MIT License](LICENSE). RedisWhistle knows the importance of freedom and spreading joy through open-source software.
