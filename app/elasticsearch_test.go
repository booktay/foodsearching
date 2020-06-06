package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputDataloadReviewsAndKeyword(t *testing.T) {

	reviewDatas, foodKeyword, err := loadReviewsAndKeyword()

	assert.Equal(t, nil, err, "Equal")
	assert.Equal(t, 6203, len(reviewDatas), "Equal")
	assert.Equal(t, 20000, len(foodKeyword), "Equal")
}

func TestInputFoodInDict(t *testing.T) {
	startElasticsearchConnection()

	input := "assorted coffee"
	foodInDict := searchFoodInDictionary(input)
	results := foodInDict["keyword"].(interface{})
	assert.Equal(t, input, results, "Equal")
}

func TestInputFoodNotInDict(t *testing.T) {
	startElasticsearchConnection()

	input := "ไก่"
	output := map[string]interface{} {
		"Message" : "Message Keyword is not found",
	}
	foodInDict := searchFoodInDictionary(input)
	assert.Equal(t, output, foodInDict, "Equal")
}

func TestInputSearchCorrectID(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
	} {
		{
			Input: "1",
		},
		{
			Input: "6203",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultID := searchByMatchID(tc.Input)
			assert.Equal(t, resultID["_id"].(string), tc.Input)
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

func TestInputSearchKeywordInDict(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "ข้าวผัดกากหมู",
			Output: map[string]interface{} {
				"relation": "eq",
				"value": float64(2789),
			},
		},
		{
			Input: "blueberry icecream crepe",
			Output: map[string]interface{} {
				"relation": "eq",
				"value": float64(6),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultKeyword := searchByMatchKeyword(tc.Input)
			assert.Equal(t, resultKeyword["total"].(map[string]interface{}), tc.Output)
		})
	}
}

func TestInputSearchKeywordNotinDict(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "ข้าวหมูแดง",
			Output: map[string]interface{} {
				"message": "Food keyword isn't in 20,000 keywords",
			},
		},
		{
			Input: "คุกกี้ไก่กรอบราดซอส",
			Output: map[string]interface{} {
				"message": "Food keyword isn't in 20,000 keywords",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultKeyword := searchByMatchKeyword(tc.Input)
			assert.Equal(t, resultKeyword, tc.Output)
		})
	}
}

func TestInputGetNumOfDocsofExistIndex(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output float64
	} {
		{
			Input: "reviews",
			Output: 6203,
		},
		{
			Input: "foods",
			Output: 20000,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultNumberofDocument := getNumberOfDocumentInIndex(tc.Input)
			assert.Equal(t, resultNumberofDocument, tc.Output)
		})
	}
}

func TestInputGetNumOfDocsofNoExistIndex(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output float64
	} {
		{
			Input: "cars",
			Output: 0,
		},
		{
			Input: "thread",
			Output: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			resultNumberofDocument := getNumberOfDocumentInIndex(tc.Input)
			assert.Equal(t, resultNumberofDocument, tc.Output)
		})
	}
}

// Enable When run Go Test
// func TestEditData(t *testing.T) {
// 	testCases := [] struct {
// 		ID string
// 		Text string
// 		Output map[string] interface{}
// 	} {
// 		{
// 			ID: "6100",
// 			Text: "โต๊ะไม่ค่อยสะอาด วางแขนไปเหนียวหนึบเลย ราคาก็ไม่ถูกแล้ว ราคาในเมนูไม่ net นะ มีคิดพวก service charge เพิ่มอีก",
// 			Output: map[string] interface{} {
// 				"id": "6100",
// 				"result": "updated",
// 			},
// 		},
// 		{
// 			ID: "5000",
// 			Text: "ข้าวผัดคอหมูย่าง\n ข้าวผัดมากลิ่นหอม คอหมูย่างอร่อยดี แต่ให้มาน้อยชิ้นไปหน่อย\n 'Piglet Go` \"`",
// 			Output: map[string] interface{} {
// 				"Message": "Error when updated",
// 				"result": "Not updated",
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.ID, func(t *testing.T) {
// 			t.Parallel()
// 			resultEdit := editReviewsByMatchID(tc.ID, string(tc.Text))
// 			assert.Equal(t, tc.Output, resultEdit)
// 		})
// 	}
// }

// Enable Only First Time to create Elasticsearch documents
// func TestInsertBulkDocument(t *testing.T) {
// 	startElasticsearchConnection()

// 	message, err := insertBulkDocument()
// 	assert.Equal(t, "Sucessfuly indexed documents", message)
// 	assert.Equal(t, nil, err)
// }