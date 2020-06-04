package main

import (
	"log"
)

func main() {
	log.Println("Starting the Container")
	startElasticsearchConnection()
	startServer()
}