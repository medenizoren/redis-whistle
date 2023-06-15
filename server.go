package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// A config represents the server configuration.
type config struct {
	port     int
	fileName string
}

// A RedisServer represents a Redis server.
type RedisServer struct {
	config     *config
	logger     *log.Logger
	databases  []*Database
	selectedDB int
	mu         sync.Mutex
}

// Init initializes the redis server.
func (server *RedisServer) Init() {
	for i := 0; i < 16; i++ {
		server.databases = append(server.databases, NewDatabase(i))
	}

	server.selectedDB = 0
	server.StartDB(server.config.fileName)
}

// StartDB starts the database.
func (server *RedisServer) StartDB(fileName string) {
	server.databases[server.selectedDB].checkAndRemoveExpiredKeys()
	server.databases[server.selectedDB].Init(fileName)
}

// SelectDB selects the database with the given index.
// It closes the current database and opens the new one.
// It also updates the selectedDB field.
func (server *RedisServer) SelectDB(index int) {
	server.databases[server.selectedDB].Close()

	server.mu.Lock()
	server.selectedDB = index
	server.mu.Unlock()

	server.StartDB("")
}

// Run runs the server.
// It listens for connections and handles them.
func (server *RedisServer) Run() {
	l, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(server.config.port))
	if err != nil {
		server.logger.Fatal(err)
	}

	defer l.Close()

	server.logger.Printf("Listening on port %d\n", server.config.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			server.logger.Fatal("Error accepting connection: ", err.Error())
		}

		go server.handleRequest(conn)
	}
}

// handleRequest handles a client request.
// It reads the request, parses it and sends the response.
func (server *RedisServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	commandMap := getCommandMap()

	for {
		value, err := DecodeRESP(reader)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			server.logger.Println("Error decoding RESP: ", err.Error())
			return // Ignore clients that we fail to read from
		}

		comingCommand := strings.ToUpper(value.Array()[0].String())
		args := value.StringArray()[1:]

		// If the command is in the command map, execute it
		// Otherwise, return an error
		if command, ok := commandMap[comingCommand]; ok {
			response := command(args)

			_, err := conn.Write([]byte(response))
			if err != nil {
				server.logger.Println("Error writing to connection: ", err.Error())
				return
			}
		} else {
			_, err := conn.Write([]byte(returnError(fmt.Sprintf("Unknown command '%s'", comingCommand))))
			if err != nil {
				server.logger.Println("Error writing to connection: ", err.Error())
				return
			}
		}
	}
}
