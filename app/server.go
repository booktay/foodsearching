package main

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func startServer() {
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
		result := searchByMatchID(reviewID)
		c.SecureJSON(http.StatusOK, result)
	}
}

func getReviewsByKeyword (c *gin.Context) {
	reviewtext := c.DefaultQuery("query", "")
	if reviewtext != "" {
		result := searchByMatchKeyword(reviewtext)
		c.SecureJSON(http.StatusOK, result)
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