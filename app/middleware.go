package main

import (
	"encoding/json"
)

func checkInvalidJson(text string) bool {
	// Check for JSON errors
	if json.Valid([]byte(text)) {
		return false
	} else {
		return true
	}
}

func checkHaveFoodKeyword(text string) bool {
	document := searchFoodInDictionary(text)
	if _, ok := document["Message"]; ok {
		return false
	} else if _, ok := document["Value"]; ok {
		return true
	}
	return false
}