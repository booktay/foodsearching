package main

import (
	"testing"
	// "reflect"

	"github.com/stretchr/testify/assert"
)

func TestInputDataloadReviewsAndKeyword(t *testing.T) {

	reviewDatas, foodKeyword, err := loadReviewsAndKeyword()

	assert.Equal(t, nil, err, "Equal")
	assert.Equal(t, 6203, len(reviewDatas), "Equal")
	assert.Equal(t, 20000, len(foodKeyword), "Equal")
}
