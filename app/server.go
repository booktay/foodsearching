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