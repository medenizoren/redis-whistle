// Description: This file contains the main server logic
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
