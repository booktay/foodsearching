package main

import (
	"bytes"
	"encoding/json"
)

func checkValidJson(text string) string {
	// Check for JSON errors
	if json.Valid([]byte(text)) {
		return text
	} else {
		return "{}"
	}
}