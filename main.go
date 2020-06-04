package main

import (
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/elastic/go-elasticsearch/v8"
)

var ES01IP = flag.String("ES01IP", "http://172.20.0.3:9200", "ES01 IP Address")

func main() {
	log.Print("Starting the Database Server")
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
			log.Println("Connected to Elasticsearch :", r["name"], " :", *ES1IPAddress)
			res.Body.Close()
			break
		} else {
			log.Println("Waiting for connection...")
			time.Sleep(5 * time.Second)
		}
	}

	log.Print("Starting the Server")
	router := gin.Default()
	router.GET("/reviews", getReviewsByKeyword)
	router.GET("/reviews/:id", getReviewsByID)
	router.PUT("/reviews/:id", editReviewsByID)
	router.Run()
}

func getReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID != "" {
		c.SecureJSON(http.StatusOK, gin.H{
			"reviewID": reviewID,
		})
	}
}

func getReviewsByKeyword (c *gin.Context) {
	keyword := c.DefaultQuery("query", "")
	if keyword != "" {
		c.SecureJSON(http.StatusOK, gin.H{
			"Keyword": keyword,
		})
	}
}

func editReviewsByID (c *gin.Context) {
	body := c.Request.Body
	if body != nil {
		reviewText, _ := ioutil.ReadAll(body)
		c.SecureJSON(http.StatusOK, gin.H{
			"reviewText": string(reviewText),
		})
	}
}