package common_unit_test

import (
	"testing"

	"github.com/TeaChanathip/touch-grass-scheduler/server/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestStructToSnakeMap_Success(t *testing.T) {
	// ------------------ Arrange ------------------
	type TestStruct struct {
		Field      string `json:"field"`
		SnakeField string `json:"snake_field"`
	}

	testStruct := &TestStruct{
		Field:      "normal field",
		SnakeField: "snake case field",
	}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(testStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.NotEmpty(t, resultMap)
	assert.Contains(t, resultMap, "field")
	assert.Equal(t, "normal field", resultMap["field"])
	assert.Contains(t, resultMap, "snake_field")
	assert.Equal(t, "snake case field", resultMap["snake_field"])
}

func TestStructToSnakeMap_EmptyStruct(t *testing.T) {
	// ------------------ Arrange ------------------
	type EmptyStruct struct{}
	emptyStruct := &EmptyStruct{}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(emptyStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Empty(t, resultMap)
}

func TestStructToSnakeMap_NilPointer(t *testing.T) {
	// ------------------ Arrange ------------------
	var nilStruct *struct{} = nil

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(nilStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Nil(t, resultMap["field"]) // JSON marshals nil pointer to null
}

func TestStructToSnakeMap_NestedStruct(t *testing.T) {
	// ------------------ Arrange ------------------
	type NestedStruct struct {
		Name string `json:"name"`
	}
	type TestStruct struct {
		Field  string       `json:"field"`
		Nested NestedStruct `json:"nested_data"`
	}

	testStruct := &TestStruct{
		Field:  "test",
		Nested: NestedStruct{Name: "nested"},
	}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(testStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Contains(t, resultMap, "nested_data")
	nestedMap, ok := resultMap["nested_data"].(map[string]any)
	assert.True(t, ok)
	assert.Contains(t, nestedMap, "name")
	assert.Equal(t, "nested", nestedMap["name"])
}

func TestStructToSnakeMap_WithIgnoredFields(t *testing.T) {
	// ------------------ Arrange ------------------
	type TestStruct struct {
		IncludeField string `json:"include_field"`
		IgnoreField  string `json:"-"` // Should be ignored
		EmptyTag     string `json:""`  // Should use field name
	}

	testStruct := &TestStruct{
		IncludeField: "included",
		IgnoreField:  "should not appear",
		EmptyTag:     "empty tag field",
	}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(testStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Contains(t, resultMap, "include_field")
	assert.NotContains(t, resultMap, "IgnoreField") // Should be ignored
	assert.Contains(t, resultMap, "EmptyTag")       // Uses field name when json tag is empty
}

func TestStructToSnakeMap_WithOmitEmptyFields(t *testing.T) {
	// ------------------ Arrange ------------------
	type TestStruct struct {
		OmitEmpty1 string `json:"omit1,omitempty"` // Should be included if not empty
		OmitEmpty2 string `json:"omit2,omitempty"`
	}

	testStruct := &TestStruct{
		OmitEmpty1: "not empty",
		OmitEmpty2: "",
	}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(testStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Contains(t, resultMap, "omit1")
	assert.NotContains(t, resultMap, "omit2")
}

func TestStructToSnakeMap_WithDifferentTypes(t *testing.T) {
	// ------------------ Arrange ------------------
	type TestStruct struct {
		StringField string  `json:"string_field"`
		IntField    int     `json:"int_field"`
		FloatField  float64 `json:"float_field"`
		BoolField   bool    `json:"bool_field"`
	}

	testStruct := &TestStruct{
		StringField: "test",
		IntField:    42,
		FloatField:  3.14,
		BoolField:   true,
	}

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(testStruct)

	// ------------------ Assert -------------------
	assert.NoError(t, err)
	assert.Equal(t, "test", resultMap["string_field"])
	assert.Equal(t, float64(42), resultMap["int_field"]) // JSON unmarshals numbers as float64
	assert.Equal(t, 3.14, resultMap["float_field"])
	assert.Equal(t, true, resultMap["bool_field"])
}

func TestStructToSnakeMap_InvalidInput(t *testing.T) {
	// ------------------ Arrange ------------------
	// Test with a type that can't be marshaled to JSON
	invalidInput := make(chan int) // channels can't be marshaled to JSON

	// ------------------ Act ----------------------
	resultMap, err := common.StructToSnakeMap(invalidInput)

	// ------------------ Assert -------------------
	assert.Error(t, err)
	assert.Nil(t, resultMap)
}
