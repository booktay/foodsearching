package main

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	log.Print("Starting the Server")
	r := gin.Default()
	r.GET("/reviews", getReviewsByKeyword)
	r.GET("/reviews/:id", getReviewsByID)
	r.PUT("/reviews/:id", editReviewsByID)
	r.Run()

	log.Print("Starting the Database Server")
	var (
		r map[string]interface{}
	)

	elasticClient, err := elasticsearch.NewClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	} 
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