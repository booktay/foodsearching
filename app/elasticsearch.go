package main

import (
	"log"
	"encoding/json"
	"time"
	"flag"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES01IP = flag.String("ES01IP", "http://172.21.0.2:9200", "ES01 IP Address")

func startElasticsearchConnection() {
	log.Print("Starting the Elasticsearch Client")
	var (
		r map[string]interface{}
	)

	cfg := elasticsearch.Config{
		Addresses: []string {
			*ES01IP,
		},
	}

	elasticClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	for {
		res, err := elasticClient.Info()
		if err == nil {
			json.NewDecoder(res.Body).Decode(&r)
			log.Println("Connected to Elasticsearch :", r["name"])
			log.Println("IP Address :", *ES01IP)
			res.Body.Close()
			break
		} else {
			log.Println("Waiting for Elasticsearch connection...")
			time.Sleep(5 * time.Second)
		}
	}
}