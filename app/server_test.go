package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

func TestInput1IDOutDocument(t *testing.T) {
	startElasticsearchConnection()
	time.Sleep(10 * time.Second)

	router := gin.Default()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews/0", nil)
	router.ServeHTTP(w, req)

	template := `{
		"Message": "ReviewID is not found"
	}`

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, template, w.Body.String())
}