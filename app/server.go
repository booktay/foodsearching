package main

import (
	"log"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/contrib/static"
	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func startServer() {
	log.Println("Starting the Server")
	router := setupRouter()
	router.Run()
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	
	// config := cors.DefaultConfig()
	
	// config.AllowOrigins = []string {
	// 	"http://localhost:3000",
	// }
	// config.AllowMethods = []string {"PUT", "GET", "OPTIONS"}
	// config.AllowHeaders = []string {"Origin"}

	// Use CORS
	// router.Use(cors.New(config))

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

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
	result := searchByMatchKeyword(reviewtext)
	c.SecureJSON(http.StatusOK, result)
}

func editReviewsByID (c *gin.Context) {
	reviewID := c.Param("id")
	reviewText, _ := ioutil.ReadAll(c.Request.Body)
	result := editReviewsByMatchID(reviewID, string(reviewText))
	c.SecureJSON(http.StatusOK, result)
}