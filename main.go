package main

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("Starting the Server")
	r := gin.Default()
	r.GET("/reviews", getReviewsByKeyword)
	r.GET("/reviews/:id", getReviewsByID)
	r.PUT("/reviews/:id", editReviewsByID)
	r.Run()
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