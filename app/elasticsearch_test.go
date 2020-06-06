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
	input := "assorted coffee"
	foodInDict := searchFoodInDictionary(input)
	results := foodInDict["_id"].(string)
	assert.Equal(t, input, results, "Equal")
}

func TestInputSearchCorrectID(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "1",
			Output: map[string]interface{} {
				"_id": "1",
				"_index": "reviews",
				"_score": 8.327645,
				"_type": "_doc",
			},
		},
		{
			Input: "6203",
			Output: map[string]interface{} {
				"_id": "6203",
				"_index": "reviews",
				"_score": 8.327645,
				"_type": "_doc",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultID := searchByMatchID(tc.Input)
			delete(resultID, "_source")
			assert.Equal(t, resultID, tc.Output)
		})
	}
}

func TestInputSearchInCorrectID(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "0",
			Output: map[string]interface{} {
				"Message": "ReviewID is not found",
			},
		},
		{
			Input: "1 or 1=1",
			Output: map[string]interface{} {
				"Message": "ReviewID is not found",
			},
		},
		{
			Input: "Helloworld 1",
			Output: map[string]interface{} {
				"Message": "ReviewID is not found",
			},
		},
		{
			Input: "Helloworld",
			Output: map[string]interface{} {
				"Message": "ReviewID is not found",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultID := searchByMatchID(tc.Input)
			assert.Equal(t, resultID, tc.Output)
		})
	}
}