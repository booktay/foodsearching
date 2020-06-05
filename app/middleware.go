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