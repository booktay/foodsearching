package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestcheckHaveFoodKeyword(t *testing.T) {
	startElasticsearchConnection()

	results := checkHaveFoodKeyword("assorted coffee")
	assert.Equal(t, true, results)
}