package jman_test

import (
	"testing"
	"time"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

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

// func TestEqual_SliceSimple_String(t *testing.T) {
// 	expected := `["hello", "world"]`
// 	actual := `["hello", "world"]`
//
// 	assert.NoError(t, jman.Equal(expected, actual))
// }
func TestEqual_SliceSimple_Bytes(t *testing.T) {}
func TestEqual_SliceSimple_Array(t *testing.T) {}

func TestEqual_SliceNotSameOrder(t *testing.T) {}

func TestEqual_SliceObjects(t *testing.T) {}

func TestEqual_ComplexObj(t *testing.T) {}

func TestNew_UnexpectedType(t *testing.T) {}

func TestEqual_InvalidJSON(t *testing.T) {}

/*
 TODO:

--- Basic Functionality ---
- make sure works with null/nil values
- add test for trying to compare object against array

--- Options ---
- add support for options for placeholders (naming tbd)
- add support for options on equal to add/subtract/edit from actual
- add support for ignore order option for array

--- Error Display ---
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
*/
