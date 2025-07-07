package jman_test

import (
	"strings"
	"testing"
	"time"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockT struct {
	mock.Mock
	gotMsg  string
	wantMsg string
}

func (m *MockT) Fatalf(format string, args ...any) {
	m.Called(format)
	m.gotMsg = format
	panic(format) // panic to simulate Fatalf behavior
}

func (m *MockT) AssertExpectations(t *testing.T) {
	m.Mock.AssertExpectations(t)
	gotMsgLines := strings.Split(m.gotMsg, "\n")
	wantMsgLines := strings.Split(m.wantMsg, "\n")
	assert.Equal(t, len(wantMsgLines), len(gotMsgLines), "number of lines in error message should match")
	for _, wantMsgLine := range wantMsgLines {
		assert.Contains(t, m.gotMsg, wantMsgLine, "got message should contain expected line")
	}
}

func newMockT(expectedMsg string) *MockT {
	mockT := &MockT{
		wantMsg: expectedMsg,
	}
	mockT.On("Fatalf", mock.Anything).Return()
	return mockT
}

func TestNewObj_String(t *testing.T) {
	data := `{"key1": "value1", "key2": 2, "key3": true}`
	obj := jman.New[jman.Obj](t, data)

	assert.Equal(t, jman.Obj{
		"key1": "value1",
		"key2": float64(2),
		"key3": true,
	}, obj)
}

func TestNewObj_Bytes(t *testing.T) {
	data := []byte(`{"key1": "value1", "key2": 2, "key3": true}`)
	obj := jman.New[jman.Obj](t, data)

	assert.Equal(t, jman.Obj{
		"key1": "value1",
		"key2": float64(2),
		"key3": true,
	}, obj)
}

func TestNewObj_UnsupportedType(t *testing.T) {
	mt := newMockT("int unsupported type for JSON data, use either string or []byte")
	defer mt.AssertExpectations(t)
	data := 12345 // unsupported type
	assert.Panics(t, func() { jman.New[jman.Obj](mt, data) })
}

func TestNewObj_InvalidJSON(t *testing.T) {
	mt := newMockT("error parsing JSON data {\"key1\": \"value1\", \"key2\": 2, \"key3\": true,: unexpected end of JSON input")
	defer mt.AssertExpectations(t)
	data := `{"key1": "value1", "key2": 2, "key3": true,` // invalid JSON
	assert.Panics(t, func() { jman.New[jman.Obj](mt, data) })
}

func TestObj_Equal_Strings(t *testing.T) {
	expected := jman.Obj{
		"foo":   "bar",
		"hello": "world",
		"baz":   123,
	}
	actual := `{"foo": "bar", "hello": "world", "baz": 123}`
	expected.Equal(t, actual)
}

func TestObj_Equal_Bytes(t *testing.T) {
	expected := jman.Obj{
		"foo":   "bar",
		"hello": "world",
		"baz":   123,
	}
	actual := []byte(`{"foo": "bar", "hello": "world", "baz": 123}`)
	expected.Equal(t, actual)
}

func TestObj_Equal_Objects(t *testing.T) {
	expected := jman.Obj{
		"foo":   "bar",
		"hello": "world",
		"baz":   123,
	}
	actual := jman.Obj{
		"foo":   "bar",
		"hello": "world",
		"baz":   123,
	}
	expected.Equal(t, actual)
}

// should work with any marshallable type
func TestObj_Equal_ObjTime(t *testing.T) {
	now := time.Now()
	expected := jman.Obj{
		"current_time": now,
	}
	actual := jman.Obj{
		"current_time": now,
	}
	expected.Equal(t, actual)
}

func TestObj_Equal_NestedObjects(t *testing.T) {
	expected := jman.Obj{
		"outer": jman.Obj{
			"inner": "value",
		},
	}
	actual := jman.Obj{
		"outer": jman.Obj{
			"inner": "value",
		},
	}
	expected.Equal(t, actual)
}

func TestObj_Equal_NestedSlices(t *testing.T) {
	expected := jman.Obj{
		"outer": jman.Arr{
			jman.Obj{"inner": "value1"},
			jman.Obj{"inner": "value2"},
		},
	}
	actual := jman.Obj{
		"outer": jman.Arr{
			jman.Obj{"inner": "value1"},
			jman.Obj{"inner": "value2"},
		},
	}
	expected.Equal(t, actual)
}

func TestObj_Equal_NullValues(t *testing.T) {
	expected := jman.Obj{
		"key1": nil,
		"key2": "value",
	}
	actual := `{"key1": null, "key2": "value"}`
	expected.Equal(t, actual)
}

func TestObj_Equal_ComplexObjects(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
		"key2": 42,
		"key3": true,
		"nested": jman.Obj{
			"innerKey1": jman.Arr{"innverValue1", jman.Obj{"innerKey2": 3.14}},
			"innerKey2": 3.14,
		},
	}
	actual := `{"key1": "value1", "key2": 42, "key3": true, "nested": {"innerKey1": ["innverValue1", {"innerKey2": 3.14}], "innerKey2": 3.14}}`
	expected.Equal(t, actual)
}

func TestObj_Equal_ComparedArray(t *testing.T) {
	mt := newMockT("can't compare json object with array")
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"key1": "value1",
		"key2": jman.Arr{"item1", "item2"},
	}
	actual := jman.Arr{
		"item1",
		"item2",
	}
	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Matcher_SimpleObj(t *testing.T) {
	expected := jman.Obj{
		"hello": "$ANY",
	}
	actual := `{"hello": "world"}`
	expected.Equal(t, actual,
		jman.WithMatchers(jman.NotEmpty("$ANY")),
	)
}

func TestObj_Equal_Matcher_ComplexObj(t *testing.T) {
	expected := jman.Obj{
		"users": jman.Arr{
			jman.Obj{
				"id":   "$UUID",
				"name": "$ANY",
				"age":  30,
			},
			jman.Obj{
				"id":   "$UUID",
				"name": "$BOBS_NAME",
				"age":  25,
			},
		},
		"meta": jman.Obj{
			"count":  2,
			"status": "active",
		},
	}
	actual := jman.Obj{
		"users": jman.Arr{
			jman.Obj{
				"id":   "6bd8f7c1-a528-4829-8a98-2003066697b0",
				"name": "Alice",
				"age":  30,
			},
			jman.Obj{
				"id":   "5facaa2a-77c3-40b9-9fa7-9f7b3823bdac",
				"name": "Bob",
				"age":  25,
			},
		},
		"meta": jman.Obj{
			"count":  2,
			"status": "active",
		},
	}

	expected.Equal(t, actual,
		jman.WithMatchers(
			jman.IsUUID("$UUID"),
			jman.NotEmpty("$ANY"),
			jman.EqualMatcher("$BOBS_NAME", "Bob"),
		),
	)
}

func TestObj_Equal_IgnoreArrayOrder(t *testing.T) {
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item3", "item1", "item2"},
	}

	expected.Equal(t, actual,
		jman.WithIgnoreArrayOrder("$.items"),
	)
}

func TestObj_Equal_IgnoreArrayOrder_IncorrectPathStart(t *testing.T) {
	mt := newMockT("invalid options: path must start with $")
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item3", "item1", "item2"},
	}
	assert.Panics(t, func() {
		expected.Equal(mt, actual,
			jman.WithIgnoreArrayOrder("items"),
		)
	})
}

func TestObj_Equal_Unequal_DifferentTypes(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.key expected "value" - actual 123
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{"key": "value"}
	actual := jman.Obj{"key": 123} // different type

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_DifferentTypesSubObjArr(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.key1 expected object - got jman.Arr ([subValue1])
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"key1": jman.Obj{"subKey1": "subValue1"},
	}
	actual := jman.Obj{
		"key1": jman.Arr{"subValue1"}, // different type (array instead of object)
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_MissingKeyFromActual(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.key2 not found in actual
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"key1": "value1",
		"key2": "value2",
	}
	actual := jman.Obj{
		"key1": "value1",
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_UnexpectedKeyInActual(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.key2 unexpected key
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"key1": "value1",
	}
	actual := jman.Obj{
		"key1": "value1",
		"key2": "value2", // unexpected key
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_TooFewItemsInArray(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.items expected 3 items - got 2 items
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item1", "item2"}, // too few items
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_TooManyItemsInArray(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.items expected 2 items - got 3 items
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"}, // too many items
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_DisplayPath(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.foo.bar.1.baz expected "quux" - actual "quant"
`)
	defer mt.AssertExpectations(t)
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

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_ManyDisplayPaths(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.0 expected "hello" - actual "HELLO"
$.1.baz.0.key expected "value" - actual "VALUE"
$.1.baz.1 expected "quux" - actual "QUUX"
$.1.foo expected "bar" - actual "BAR"
$.2 expected "world" - actual "WORLD"
$.3.0 expected 1 - actual <nil>
$.3.2 expected false - actual 1
$.3.3 expected <nil> - actual false
`)
	defer mt.AssertExpectations(t)
	expected := jman.Arr{
		"hello",
		jman.Obj{
			"foo": "bar",
			"baz": jman.Arr{
				jman.Obj{"key": "value"},
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
				jman.Obj{"key": "VALUE"},
				"QUUX",
			},
		},
		"WORLD",
		jman.Arr{nil, 2, 1, false},
	}
	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Equal_Unequal_ManyObjectErrors(t *testing.T) {
	mt := newMockT(`expected not equal to actual:
$.key1 not found in actual
$.key3 unexpected key
$.key2.nestedKey4 not found in actual
$.key2.nestedKey3 unexpected key
$.key2.nestedKey1 expected "nestedValue1" - actual "nestedValue2"
$.key2.nestedKey2 expected 42 - actual "notANumber"
`)
	defer mt.AssertExpectations(t)
	expected := jman.Obj{
		"key1": "value1",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue1",
			"nestedKey2": 42,
			"nestedKey4": "unexpectedKey", // unexpected key
		},
	}
	actual := jman.Obj{
		"key3": "value3",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue2", // different value1
			"nestedKey2": "notANumber",   // different TestNewObj_UnsupportedType
			"nestedKey3": true,           // unexpected key3
		},
	}

	assert.Panics(t, func() { expected.Equal(mt, actual) })
}

func TestObj_Get_String(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	value := data.Get(t, "$.key1")
	assert.Equal(t, "value1", value)
}

func TestObj_Get_Number(t *testing.T) {

	data := jman.Obj{"key1": 42}
	value := data.Get(t, "$.key1")
	assert.Equal(t, float64(42), value)
}

func TestObj_Get_Bool(t *testing.T) {
	data := jman.Obj{"key1": true}
	value := data.Get(t, "$.key1")
	assert.Equal(t, true, value)
}

func TestObj_Get_Obj(t *testing.T) {
	data := jman.Obj{"key1": jman.Obj{"nestedKey": "nestedValue"}}
	value := data.Get(t, "$.key1")
	assert.Equal(t, jman.Obj{"nestedKey": "nestedValue"}, value)
}

func TestObj_Get_Arr(t *testing.T) {
	data := jman.Obj{"key1": jman.Arr{"item1", "item2"}}
	value := data.Get(t, "$.key1")
	assert.Equal(t, jman.Arr{"item1", "item2"}, value)
}

func TestObj_Get_NestedKey(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey": "nestedValue",
		},
	}
	value := data.Get(t, "$.key1.nestedKey")
	assert.Equal(t, "nestedValue", value)
}

func TestObj_Get_NestedArrItem(t *testing.T) {

	data := jman.Obj{
		"key1": jman.Arr{
			"item1",
			jman.Obj{"nestedKey": "nestedValue"},
			"item3",
		},
	}
	value := data.Get(t, "$.key1.1.nestedKey")
	assert.Equal(t, "nestedValue", value)
}

func TestObj_Get_ComplexPath(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	value := data.Get(t, "$.key1.nestedKey1.1.deepKey")
	assert.Equal(t, "deepValue", value)
}

func TestObj_Get_PathNotFound(t *testing.T) {
	mt := newMockT("key 'key2' not found in object")
	defer mt.AssertExpectations(t)
	data := jman.Obj{"key1": "value1"}
	assert.Panics(t, func() { _ = data.Get(mt, "$.key2") })
}

func TestObj_Get_IncorrectConversion(t *testing.T) {
	mt := newMockT("expected number at path '$.key1', got string")
	defer mt.AssertExpectations(t)

	data := jman.Obj{"key1": "value1"}
	assert.Panics(t, func() { _ = data.GetNumber(mt, "$.key1") })
}

func TestObj_Set_NewKey(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	data.Set(t, "$.key2", "newValue")

	expected := jman.Obj{
		"key1": "value1",
		"key2": "newValue",
	}
	assert.Equal(t, expected, data)
}

func TestObj_Set_ExistingKey(t *testing.T) {

	data := jman.Obj{"key1": "value1"}
	data.Set(t, "$.key1", "updatedValue")

	expected := jman.Obj{
		"key1": "updatedValue",
	}
	assert.Equal(t, expected, data)
}

func TestObj_Set_NestedKey(t *testing.T) {

	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey": "nestedValue",
		},
	}
	data.Set(t, "$.key1.nestedKey", "newNestedValue")

	expected := jman.Obj{
		"key1": jman.Obj{
			"nestedKey": "newNestedValue",
		},
	}
	assert.Equal(t, expected, data)
}

func TestObj_Set_ComplexPath(t *testing.T) {

	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	data.Set(t, "$.key1.nestedKey1.1.deepKey", "newDeepValue")

	expected := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "newDeepValue"}},
		},
	}

	assert.Equal(t, expected, data)
}

func TestObj_Set_Obj(t *testing.T) {

	data := jman.Obj{}
	data.Set(t, "$.key1", jman.Obj{"nestedkey": "nestedValue"})
	expected := jman.Obj{
		"key1": jman.Obj{"nestedkey": "nestedValue"},
	}
	assert.Equal(t, expected, data)
}

func TestObj_Set_GenericMap(t *testing.T) {
	data := jman.Obj{}
	data.Set(t, "$.key1", map[string]any{"nestedkey": "nestedValue"})
	expected := jman.Obj{
		"key1": jman.Obj{"nestedkey": "nestedValue"},
	}
	assert.Equal(t, expected, data)
}
