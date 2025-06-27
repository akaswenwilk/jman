package jman_test

import (
	"testing"
	"time"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

// TestEqual Basic Functionality

func TestEqual_Strings(t *testing.T) {
	expected := `{"hello": "world", "foo": "bar", "baz": 123}`
	actual := `{"foo": "bar", "hello": "world", "baz": 123}`

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_Bytes(t *testing.T) {
	expected := `{"hello": "world", "foo": "bar", "baz": 123}`
	actual := `{"foo": "bar", "hello": "world", "baz": 123}`

	assert.NoError(t, jman.Equal([]byte(expected), []byte(actual)))
}

func TestEqual_Obj(t *testing.T) {
	expected := jman.Obj{
		"hello": "world",
		"foo":   "bar",
		"baz":   123,
	}
	actual := jman.Obj{
		"foo":   "bar",
		"hello": "world",
		"baz":   123,
	}
	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_Mix(t *testing.T) {
	expected := jman.Obj{
		"hello": "world",
		"foo":   "bar",
		"baz":   123,
	}
	actual := `{"foo": "bar", "hello": "world", "baz": 123}`

	assert.NoError(t, jman.Equal(expected, actual))
}

// should work with any marshallable type
func TestEqual_ObjTime(t *testing.T) {
	now := time.Now()
	expected := jman.Obj{
		"current_time": now,
	}
	actual := jman.Obj{
		"current_time": now,
	}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_NestedObj_Obj(t *testing.T) {
	expected := jman.Obj{
		"hello": "world",
		"foo": jman.Obj{
			"bar": "baz",
		},
	}
	actual := jman.Obj{
		"foo": jman.Obj{
			"bar": "baz",
		},
		"hello": "world",
	}
	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_NestedObj_String(t *testing.T) {
	expected := `{
			"hello": "world",
			"foo": {
				"bar": "baz"
			}
	}`
	actual := `{
			"foo": {
				"bar": "baz"
			},
			"hello": "world"
	}`

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_SliceSimple_String(t *testing.T) {
	expected := `["hello", "world", 123, false]`
	actual := `["hello", "world", 123, false]`

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_SliceSimple_Bytes(t *testing.T) {
	expected := []byte(`["hello", "world"]`)
	actual := []byte(`["hello", "world"]`)

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_SliceSimple_Array(t *testing.T) {
	expected := jman.Arr{"hello", "world", 123, false}
	actual := jman.Arr{"hello", "world", 123, false}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_SliceSimple_Mixed(t *testing.T) {
	expected := `["hello", "world", 123, false]`
	actual := jman.Arr{"hello", "world", 123, false}

	assert.NoError(t, jman.Equal(expected, actual))
}

// func TestEqual_SliceNotSameOrder(t *testing.T) {
// 	expected := jman.Arr{"hello", "world", 123, false}
// 	actual := jman.Arr{"world", "hello", 123, false}
//
// 	// This should fail because the order is different
// 	err := jman.Equal(expected, actual)
// 	assert.Error(t, err)
// 	assert.Equal(t, `expected not equal to actual:
// $.0 expected "hello" - actual "world"
// $.1 expected "world" - actual "hello"
// `, err.Error())
// }

func TestEqual_SliceNullValues(t *testing.T) {
	expected := `[null, "world", 123, false]`
	actual := jman.Arr{nil, "world", 123, false}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_SliceOfObjects(t *testing.T) {
	expected := jman.Arr{
		jman.Obj{"name": "Alice", "age": 30},
		jman.Obj{"name": "Bob", "age": 25},
	}
	actual := jman.Arr{
		jman.Obj{"name": "Alice", "age": 30},
		jman.Obj{"name": "Bob", "age": 25},
	}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_ObjectOfSlices(t *testing.T) {
	expected := jman.Obj{
		"names": jman.Arr{"Alice", "Bob"},
		"ages":  jman.Arr{30, 25},
	}
	actual := jman.Obj{
		"ages":  jman.Arr{30, 25},
		"names": jman.Arr{"Alice", "Bob"},
	}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_ComplexObj(t *testing.T) {
	expected := jman.Obj{
		"users": jman.Arr{
			jman.Obj{"name": "Alice", "age": 30, "active": true},
			jman.Obj{"name": "Bob", "age": 25, "active": false},
		},
		"meta": jman.Obj{
			"count":  2,
			"status": "ok",
		},
	}
	actual := jman.Obj{
		"meta": jman.Obj{
			"status": "ok",
			"count":  2,
		},
		"users": jman.Arr{
			jman.Obj{"name": "Alice", "age": 30, "active": true},
			jman.Obj{"name": "Bob", "age": 25, "active": false},
		},
	}

	assert.NoError(t, jman.Equal(expected, actual))
}

func TestEqual_DifferentTypes_ExpectObj(t *testing.T) {
	expected := jman.Obj{
		"key": "value",
	}
	actual := jman.Arr{"key", "value"}
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "expected a json object, got []interface {}", err.Error())
}

func TestEqual_DifferentTypes_ExpectArr(t *testing.T) {
	expected := jman.Arr{"key", "value"}
	actual := jman.Obj{
		"key": "value",
	}
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "expected a json array, got map[string]interface {}", err.Error())
}

func TestEqual_UnexpectedType(t *testing.T) {
	expected := jman.Obj{
		"key": "value",
	}
	actual := 123 // This is an unexpected type

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "invalid actual: type %!s(int=123) not supported. use either string, []byte, jman.Obj, or jman.Arr", err.Error())
}

func TestEqual_InvalidJSON_Expected(t *testing.T) {
	expected := `{"key": "value"`
	actual := `{"key": "value"}`
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "invalid expected: error unmarshalling JSON string {\"key\": \"value\": unexpected end of JSON input", err.Error())
}

func TestEqual_InvalidJSON_Actual(t *testing.T) {
	expected := `{"key": "value"}`
	actual := `{"key": "value"`
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "invalid actual: error unmarshalling JSON string {\"key\": \"value\": unexpected end of JSON input", err.Error())
}

// TestEqual Options

func TestEqual_Matcher_SimpleObj(t *testing.T) {
	expected := `{"hello": "$ANY"}`
	actual := `{"hello": "world"}`

	// Using a placeholder for the 'baz' field
	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithMatchers(
			jman.NotEmpty("$ANY"),
		),
	))
}

func TestEqual_Matcher_SimpleArray(t *testing.T) {
	expected := `["$ANY", "world"]`
	actual := `["hello", "world"]`

	// Using a placeholder for the 'baz' field
	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithMatchers(
			jman.NotEmpty("$ANY"),
		),
	))
}

func TestEqual_Matcher_ComplexObj(t *testing.T) {
	expected := `{
		"users": [
			{"name": "$ANY", "age": 30},
			{"name": "Bob", "age": 25}
		],
		"meta": {"count": 2, "status": "ok"}
	}`
	actual := `{
		"users": [
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25}
		],
		"meta": {"status": "ok", "count": 2}
	}`

	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithMatchers(
			jman.NotEmpty("$ANY"),
		),
	))
}

func TestEqual_Matcher_SimpleEqualsMatcher(t *testing.T) {
	expected := `{"hello": "$MY_DYNAMIC_VALUE"}`
	actual := `{"hello": "world"}`

	// Using a placeholder for the 'baz' field
	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithMatchers(
			jman.EqualMatcher("$MY_DYNAMIC_VALUE", "world"),
		),
	))
}

func TestEqual_Matcher_SimpleUUIDMatcher(t *testing.T) {
	expected := `{"id": "$MY_UUID"}`
	actual := `{"id": "123e4567-e89b-12d3-a456-426614174000"}`

	// Using a placeholder for the 'baz' field
	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithMatchers(
			jman.IsUUID("$MY_UUID"),
		),
	))
}

func TestEqual_IgnoreArrayOrder(t *testing.T) {
	expected := `{"items": ["apple", "banana", "cherry"]}`
	actual := `{"items": ["cherry", "banana", "apple"]}`

	// This should pass because we ignore the order of items in the array
	assert.NoError(t, jman.Equal(expected, actual,
		jman.WithIgnoreArrayOrder("$.items"),
	))
}

/*
 TODO:
--- Options ---
- add support for options on equal to add/subtract/edit from actual
- add support for ignore order option for array

--- Error Display ---
- test error display for mismatched types for each type
- test error display for different values for each type
- test for nice error display
- add subdifferences plus indented error messages to show differences in nested fields
- color output

--- Helper Methods ---
- add support for EqualValue (dot notation plus value)
- add has key method
- method to take a json payload and return a modified json payload
- add a marshal json and a must marshal json to Obj/Arr


--- Docu ---
- add godocs
- update readme
- add note that whatever is passed in as expected must be convertible into valid json as an array or object
- fix the TOC
- note that the dot notation must start with $
*/
