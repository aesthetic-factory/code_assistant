package util

import (
	"encoding/json"
	"fmt"
)

// Helper function
func ParseJsonObject[T any](input string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return result, fmt.Errorf("error parsing JSON %s", err)
	}
	return result, nil
}

// Helper function
func ParseJsonArray[T any](input string) ([]T, error) {

	// Parse the JSON data into a slice of structs
	var results []T
	if err := json.Unmarshal([]byte(input), &results); err != nil {
		return results, fmt.Errorf("error parsing JSON %s", err)
	}
	return results, nil
}
