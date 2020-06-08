package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputFoodInDict(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output bool
	} {
		{
			Input: "ข้าวพันผักสาหร่ยเส้นแก้วไข่ดาว",
			Output: true,
		},
		{
			Input: "ผัดเห็ดสามอย่าง ผัดเผ็ดกบ สปาเกตตี้",
			Output: true,
		},
		{
			Input: "ข้าวขาหมู ก๋วยเตี๋ยวน้ำใส",
			Output: true,
		},
		{
			Input: "homemade marshmallow",
			Output: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			foodInDict := searchFoodInDictionary(tc.Input)
			results := false
			if _, ok := foodInDict["Value"]; ok {
				results = true
			}
			assert.Equal(t, tc.Output, results, "Equal")
		})
	}
}

func TestInputFoodNotInDict(t *testing.T) {
	startElasticsearchConnection()

	testCases := [] struct {
		Input string
		Output bool
	} {
		{
			Input: "soft pancake rolls christmas set",
			Output: false,
		},
		{
			Input: "banana chocolate crepe",
			Output: false,
		},
		{
			Input: "เค้กปริ้นเซส & โรส บัทเทอร์วานิลลา",
			Output: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()
			foodInDict := searchFoodInDictionary(tc.Input)
			results := false
			if _, ok := foodInDict["Value"]; ok {
				results = true
			}
			assert.Equal(t, tc.Output, results, "Equal")
		})
	}
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
			Input: "เค้กปริ้นเซส & โรส บัทเทอร์วานิลลา",
			Output: map[string]interface{} {
				"message": "Food keyword isn't in 20,000 keywords",
			},
		},
		{
			Input: "soft pancake rolls christmas set",
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

// Be careful to enable this function
// It will replace a real review with a mockup review.
// Please enable only for the test to edit the review text
//
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

// Be careful to enable this function
// It will replace real documents with mock documents.
// Please enable only for the test to create Elasticsearch documents
//
// func TestInsertBulkDocument(t *testing.T) {
// 	startElasticsearchConnection()

// 	message, err := insertBulkDocument()
// 	assert.Equal(t, "Sucessfuly indexed documents", message)
// 	assert.Equal(t, nil, err)
// }