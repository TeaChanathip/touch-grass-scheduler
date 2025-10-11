package common

import (
	"encoding/json"
)

// StructToSnakeMap converts a struct to map[string]any using json tags for snake_case keys
func StructToSnakeMap(s any) (map[string]any, error) {
	jsonBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
