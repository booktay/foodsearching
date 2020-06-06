package main

import (
	"log"
)

func init() {
	log.Println("Starting the Container")
	loadEnvironment()
	startElasticsearchConnection()
}

func main() {
	startServer()
}