package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"bytes"
	
	"github.com/stretchr/testify/assert"
)

func TestGetReviews(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reviews", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)
}

func TestGetReviewsByCorrectID(t *testing.T) {
	r := setupRouter()

	testCases := [] struct {
		Input string
		Output string
	} {
		{
			Input: "/reviews/1",
			Output: "1",
		},
		{
			Input: "/reviews/6200",
			Output: "6200",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Output, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.Input, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotNil(t, w.Body)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)

			// assert field of response body
			assert.Nil(t, err)
			assert.Equal(t, tc.Output, resp["_id"].(string))
		})
	}
}

func TestGetReviewsByIncorrectID(t *testing.T) {
	r := setupRouter()

		testCases := [] struct {
		Input string
		Output string
	} {
		{
			Input: "/reviews/10000",
			Output: "ReviewID is not found",
		},
		{
			Input: "/reviews/1 or 1=1",
			Output: "ReviewID is not found",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Output, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.Input, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotNil(t, w.Body)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)

			// assert field of response body
			assert.Nil(t, err)
			assert.Equal(t, tc.Output, resp["Message"].(string))
		})
	}

}

func TestGetReviewByKeyword(t *testing.T) {
	r := setupRouter()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "/reviews?query=ข้าวผัดกากหมู",
			Output: map[string]interface{} {
				"relation": "eq",
				"value": float64(2789),
			},
		},
		{
			Input: "/reviews?query=blueberry icecream crepe",
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

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.Input, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotNil(t, w.Body)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)

			// assert field of response body
			assert.Nil(t, err)
			assert.Equal(t, tc.Output, resp["total"])
		})
	}
}

func TestGetReviewByNonkeyword(t *testing.T) {
	r := setupRouter()

	testCases := [] struct {
		Input string
		Output map[string]interface{}
	} {
		{
			Input: "/reviews?query=เค้กปริ้นเซส & โรส บัทเทอร์วานิลลา",
			Output: map[string]interface{} {
				"message": "Food keyword isn't in 20,000 keywords",
			},
		},
		{
			Input: "/reviews?query=soft pancake rolls christmas set",
			Output: map[string]interface{} {
				"message": "Food keyword isn't in 20,000 keywords",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.Input, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotNil(t, w.Body)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)

			// assert field of response body
			assert.Nil(t, err)
			assert.Equal(t, tc.Output, resp)
		})
	}
}

func TestEditReviewsByID(t *testing.T) {
	r := setupRouter()

	testCases := [] struct {
		ID string
		Text []byte
		Output map[string] interface{}
	} {
		{
			ID: "/reviews/5000",
			Text: []byte(`ข้าวผัดคอหมูย่าง\n ข้าวผัดมากลิ่นหอม\nแต่ให้มาน้อยชิ้นไปหน่อย\n 'Piglet Go`),
			Output: map[string] interface{} {
				"id": "5000",
				"result": "updated",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.ID, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", tc.ID, bytes.NewBuffer(tc.Text))
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotNil(t, w.Body)

			var resp map[string]interface{}
			err := json.Unmarshal([]byte(w.Body.String()), &resp)

			// assert field of response body
			assert.Nil(t, err)
			assert.Equal(t, tc.Output, resp)
		})
	}
}