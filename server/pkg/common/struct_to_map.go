package common

import (
	"encoding/json"
	"fmt"
)

// StructToSnakeMap converts a struct to map[string]any using json tags for snake_case keys
func StructToSnakeMap(s any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed converting struct to json bytes: %w", err)
	}

	var result map[string]any
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed converting json bytes to map: %w", err)
	}

	return result, nil
}
