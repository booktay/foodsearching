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
	result := searchByMatchID(reviewID)
	c.SecureJSON(http.StatusOK, result)
}

func getReviewsByKeyword (c *gin.Context) {
	reviewtext := c.DefaultQuery("query", "")
	result := searchByMatchKeyword(reviewtext)
	c.SecureJSON(http.StatusOK, result)
}

func editReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	reviewText, _ := ioutil.ReadAll(c.Request.Body)
	result := editReviewsByMatchID(reviewID, string(reviewText))
	c.SecureJSON(http.StatusOK, result)
}