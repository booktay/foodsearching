package main

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func startServer() {
	log.Println("Starting the Server")
	router := setupRouter()
	router.Run()
}

func setupCORSforFrontend() cors.Config {
	config := cors.DefaultConfig()
	
	config.AllowOrigins = []string {
		"http://localhost:3000",
	}
	config.AllowMethods = []string {"PUT", "GET"}
	config.AllowHeaders = []string {"Origin"}

	return config
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

	// Enable CORS for Frontend - API Testing
	config := setupCORSforFrontend()
	router.Use(cors.New(config))

	// Set API Route
	router.GET("/reviews", getReviewsByKeyword)
	router.GET("/reviews/:id", getReviewsByID)
	router.PUT("/reviews/:id", editReviewsByID)
	return router
}

func getReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	result := searchByMatchID(reviewID)
	c.SecureJSON(http.StatusOK, result)
}

func getReviewsByKeyword (c *gin.Context) {
	reviewtext := c.DefaultQuery("query", "")
	if reviewtext == "" {
		c.SecureJSON(http.StatusOK, "Reviews API")
	} else {
		result := searchByMatchKeyword(reviewtext)
		c.SecureJSON(http.StatusOK, result)
	}
}

func editReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	reviewText, _ := ioutil.ReadAll(c.Request.Body)
	result := editReviewsByMatchID(reviewID, reviewText)
	c.SecureJSON(http.StatusOK, result)
}