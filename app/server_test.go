package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"

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