package jman_test

import (
	"testing"
	"time"

	"github.com/akaswenwilk/jman"
)

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

func TestObj_Equal_MarshalableType_Time(t *testing.T) {
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
			"innerKey1": jman.Arr{"innerValue1", jman.Obj{"innerKey2": 3.14}},
			"innerKey2": 3.14,
		},
	}
	actual := `{"key1": "value1", "key2": 42, "key3": true, "nested": {"innerKey1": ["innerValue1", {"innerKey2": 3.14}], "innerKey2": 3.14}}`
	expected.Equal(t, actual)
}

func TestObj_Equal_ComparedArray(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
		"key2": jman.Arr{"item1", "item2"},
	}
	actual := jman.Arr{"item1", "item2"}

	assertFatalf(t, "can't compare json object with array", func(mt jman.T) {
		expected.Equal(mt, actual)
	})
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
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item3", "item1", "item2"},
	}
	assertFatalf(t, "invalid options: path must start with $", func(mt jman.T) {
		expected.Equal(mt, actual,
			jman.WithIgnoreArrayOrder("items"),
		)
	})
}

func TestObj_Equal_Unequal_DifferentTypes(t *testing.T) {
	expected := jman.Obj{"key": "value"}
	actual := jman.Obj{"key": 123} // different type

	assertFatalf(t, `expected not equal to actual:
expected {"key":"value"}
actual {"key":123}

$.key expected "value" - actual 123
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_DifferentTypesSubObjArr(t *testing.T) {
	expected := jman.Obj{
		"key1": jman.Obj{"subKey1": "subValue1"},
	}
	actual := jman.Obj{
		"key1": jman.Arr{"subValue1"}, // different type (array instead of object)
	}

	assertFatalf(t, `expected not equal to actual:
expected {"key1":{"subKey1":"subValue1"}}
actual {"key1":["subValue1"]}

$.key1 expected object - got jman.Arr ([subValue1])
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_MissingKeyFromActual(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
		"key2": "value2",
	}
	actual := jman.Obj{
		"key1": "value1",
	}

	assertFatalf(t, `expected not equal to actual:
expected {"key1":"value1","key2":"value2"}
actual {"key1":"value1"}

$.key2 not found in actual
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_UnexpectedKeyInActual(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
	}
	actual := jman.Obj{
		"key1": "value1",
		"key2": "value2", // unexpected key
	}

	assertFatalf(t, `expected not equal to actual:
expected {"key1":"value1"}
actual {"key1":"value1","key2":"value2"}

$.key2 unexpected key
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_TooFewItemsInArray(t *testing.T) {
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item1", "item2"}, // too few items
	}

	assertFatalf(t, `expected not equal to actual:
expected {"items":["item1","item2","item3"]}
actual {"items":["item1","item2"]}

$.items expected 3 items - got 2 items
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_TooManyItemsInArray(t *testing.T) {
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"}, // too many items
	}

	assertFatalf(t, `expected not equal to actual:
expected {"items":["item1","item2"]}
actual {"items":["item1","item2","item3"]}

$.items expected 2 items - got 3 items
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_DisplayPath(t *testing.T) {
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

	assertFatalf(t, `expected not equal to actual:
expected {"foo":{"bar":["hello",{"baz":"quux"}]}}
actual {"foo":{"bar":["hello",{"baz":"quant"}]}}

$.foo.bar.1.baz expected "quux" - actual "quant"
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_ManyDisplayPaths(t *testing.T) {
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

	assertFatalf(t, `expected not equal to actual:
expected ["hello",{"baz":[{"key":"value"},"quux"],"foo":"bar"},"world",[1,2,false,null]]
actual ["HELLO",{"baz":[{"key":"VALUE"},"QUUX"],"foo":"BAR"},"WORLD",[null,2,1,false]]

$.0 expected "hello" - actual "HELLO"
$.1.baz.0.key expected "value" - actual "VALUE"
$.1.baz.1 expected "quux" - actual "QUUX"
$.1.foo expected "bar" - actual "BAR"
$.2 expected "world" - actual "WORLD"
$.3.0 expected 1 - actual <nil>
$.3.2 expected false - actual 1
$.3.3 expected <nil> - actual false
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}

func TestObj_Equal_Unequal_ManyObjectErrors(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue1",
			"nestedKey2": 42,
			"nestedKey4": "unexpectedKey",
		},
	}
	actual := jman.Obj{
		"key3": "value3",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue2",
			"nestedKey2": "notANumber",
			"nestedKey3": true,
		},
	}

	assertFatalf(t, `expected not equal to actual:
expected {"key1":"value1","key2":{"nestedKey1":"nestedValue1","nestedKey2":42,"nestedKey4":"unexpectedKey"}}
actual {"key2":{"nestedKey1":"nestedValue2","nestedKey2":"notANumber","nestedKey3":true},"key3":"value3"}

$.key1 not found in actual
$.key3 unexpected key
$.key2.nestedKey4 not found in actual
$.key2.nestedKey3 unexpected key
$.key2.nestedKey1 expected "nestedValue1" - actual "nestedValue2"
$.key2.nestedKey2 expected 42 - actual "notANumber"
`, func(mt jman.T) {
		expected.Equal(mt, actual)
	})
}
