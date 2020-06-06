package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputgetReviewData(t *testing.T) {

	reviewDatas, err := getReviewData()

	assert.Equal(t, nil, err, "Equal")
	assert.Equal(t, 6203, len(reviewDatas), "Equal")
}

func TestOutputgetFoodKeyword(t *testing.T) {

	foodKeywords, err := getFoodKeyword()

	assert.Equal(t, nil, err, "Equal")
	assert.Equal(t, 20000, len(foodKeywords), "Equal")
}