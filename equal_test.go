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

func TestEqual_SliceNotSameOrder(t *testing.T) {
	expected := jman.Arr{"hello", "world", 123, false}
	actual := jman.Arr{"world", "hello", 123, false}

	// This should fail because the order is different
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.0 expected "hello" - actual "world"
$.1 expected "world" - actual "hello"
`, err.Error())
}

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
	assert.Equal(t, "expected a json object, got array", err.Error())
}

func TestEqual_DifferentTypes_ExpectArr(t *testing.T) {
	expected := jman.Arr{"key", "value"}
	actual := jman.Obj{
		"key": "value",
	}
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, "expected a json array, got object", err.Error())
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

func TestEqual_Error_DifferentTypes_ValuesString(t *testing.T) {
	expected := `{"key": "value"}`
	actual := `{"key": 123}` // Different type for the value

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected "value" - actual 123
`, err.Error())
}

func TestEqual_Error_DifferentTypes_ValuesNumber(t *testing.T) {
	expected := `{"key": 123}` // Different type for the value
	actual := `{"key": "value"}`

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected 123 - actual "value"
`, err.Error())
}

func TestEqual_Error_DifferentTypes_ValuesBool(t *testing.T) {
	expected := `{"key": true}` // Different type for the value
	actual := `{"key": "value"}`

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected true - actual "value"
`, err.Error())
}

func TestEqual_Error_DifferentTypes_ValuesNull(t *testing.T) {
	expected := `{"key": null}` // Different type for the value
	actual := `{"key": "value"}`

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected <nil> - actual "value"
`, err.Error())
}

func TestEqual_Error_DifferentTypes_ValuesArr(t *testing.T) {
	expected := `{"key": ["value"]}`
	actual := `{"key": "value"}`

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected array - got string (value)
`, err.Error())
}

func TestEqual_Error_DifferentTypes_ValuesObj(t *testing.T) {
	expected := `{"key": {"value":"subvalue"}}`
	actual := `{"key": ["value"]}`

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key expected object - got []interface {} ([value])
`, err.Error())
}

func TestEqual_Error_DisplayPathObj(t *testing.T) {
	expected := jman.Obj{
		"foo": jman.Obj{
			"bar": jman.Arr{
				"hello",
				jman.Obj{
					"baz": "quux",
				},
			},
		},
	}
	actual := jman.Obj{
		"foo": jman.Obj{
			"bar": jman.Arr{
				"hello",
				jman.Obj{
					"baz": "quant",
				},
			},
		},
	}

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.foo.bar.1.baz expected "quux" - actual "quant"
`, err.Error())
}

func TestEqual_Error_DisplayPathArr(t *testing.T) {
	expected := jman.Arr{1, 2, 3}
	actual := jman.Arr{1, 2, 5}

	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.2 expected 3 - actual 5
`, err.Error())
}

func TestEqual_Error_ManyDisplayPaths(t *testing.T) {
	expected := jman.Arr{
		"hello",
		jman.Obj{
			"foo": "bar",
			"baz": jman.Arr{
				jman.Obj{"key":"value"},
				"quux",
			},
		},
		"world",
		jman.Arr{1, 2, false, nil},
	}
	actual := jman.Arr{
		"HELLO",
		jman.Obj{
			"foo": "BAR",
			"baz": jman.Arr{
				jman.Obj{"key":"VALUE"},
				"QUUX",
			},
		},
		"WORLD",
		jman.Arr{nil, 2, 1, false},
	}
	err := jman.Equal(expected, actual)
	assert.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.0 expected "hello" - actual "HELLO"
$.1.baz.0.key expected "value" - actual "VALUE"
$.1.baz.1 expected "quux" - actual "QUUX"
$.1.foo expected "bar" - actual "BAR"
$.2 expected "world" - actual "WORLD"
$.3.0 expected 1 - actual <nil>
$.3.2 expected false - actual 1
$.3.3 expected <nil> - actual false
`, err.Error())
}

/*
 TODO:

--- Helper Methods ---
- add support for EqualValue (dot notation plus value)
- -
- add has key method
- method to take a json payload and return a modified json payload
- add a marshal json and a must marshal json to Obj/Arr

-- add method to check equality from file
- if not written, fail the test but write the file


--- Docu ---
- add godocs
- update readme
- add note that whatever is passed in as expected must be convertible into valid json as an array or object
- fix the TOC
- note that the dot notation must start with $
*/
