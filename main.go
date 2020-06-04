package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("Starting the Server")
	r := gin.Default()
	r.GET("/reviews/:id", getReviewsByID)
	r.Run()
}

func getReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	c.SecureJSON(http.StatusOK, gin.H{
		"reviewID": reviewID,
	})
}