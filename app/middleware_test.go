package main

import (
	"testing"

	"encoding/json"
)

func TestcheckHaveFoodKeyword() {
	startElasticsearchConnection()

	results := checkHaveFoodKeyword("assorted coffee")
	assert.Equal(t, true, results)
}