package main

import (
	"log"
)

func main() {
	log.Println("Starting the Container")
	reviewsData, _ := getReviewData()
	log.Println(reviewsData)
	// startElasticsearchConnection()
	// startServer()
}