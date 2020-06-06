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

func TestInputResultFoodInDict(t *testing.T) {
	startElasticsearchConnection()

	// resultEdit := editReviewsByMatchID("10", "")	
	// test.AssertEquals(t, map[string]interface{}, reflect.TypeOf(resultEdit))
	
	resultFoodInDict := searchFoodInDictionary("assorted coffee")
	resultFoodInDictMock := map[string]interface{} {
		"_index": "foods",
		"_type": "_doc",
		"_id": "assorted coffee",
		"_score": 1.0,
		"_source": map[string]interface{} {
			"keywordid": "19964",
			"keyword": "assorted coffee",
		},
	}
	assert.Equal(t, resultFoodInDictMock, resultFoodInDict, "Equal")
}