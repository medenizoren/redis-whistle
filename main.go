// RedisWhistle is a mini Redis clone developed in Go,
// designed to whistle away your data storage worries.
// With RedisWhistle, you can store, retrieve, and manipulate your data.

// Package main implements a simple Redis server,
// with support for the most used commands.
package main

import (
	"flag"
	"log"
	"os"
)

var redis *RedisServer

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 6379, "REDIS server port")
	flag.StringVar(&cfg.fileName, "load", "", "Load DB from a file")
	flag.Parse()

	redis = &RedisServer{
		config: &cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	redis.Init()
	redis.Run()
}
