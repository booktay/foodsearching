package main

import (
	"log"
)

func main() {
	log.Println("Starting the Container")
	foodKeywordsData, _ := getFoodKeyword()
	log.Println(foodKeywordsData)
	// startElasticsearchConnection()
	// startServer()
}