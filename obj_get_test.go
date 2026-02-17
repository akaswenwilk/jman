package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

func TestObj_Get(t *testing.T) {
	testCases := []struct {
		name string
		data jman.Obj
		path string
		want any
	}{
		{
			name: "String",
			data: jman.Obj{"key1": "value1"},
			path: "$.key1",
			want: "value1",
		},
		{
			name: "Number",
			data: jman.Obj{"key1": 42},
			path: "$.key1",
			want: float64(42),
		},
		{
			name: "Bool",
			data: jman.Obj{"key1": true},
			path: "$.key1",
			want: true,
		},
		{
			name: "Obj",
			data: jman.Obj{"key1": jman.Obj{"nestedKey": "nestedValue"}},
			path: "$.key1",
			want: jman.Obj{"nestedKey": "nestedValue"},
		},
		{
			name: "Arr",
			data: jman.Obj{"key1": jman.Arr{"item1", "item2"}},
			path: "$.key1",
			want: jman.Arr{"item1", "item2"},
		},
		{
			name: "NestedKey",
			data: jman.Obj{
				"key1": jman.Obj{
					"nestedKey": "nestedValue",
				},
			},
			path: "$.key1.nestedKey",
			want: "nestedValue",
		},
		{
			name: "NestedArrItem",
			data: jman.Obj{
				"key1": jman.Arr{
					"item1",
					jman.Obj{"nestedKey": "nestedValue"},
					"item3",
				},
			},
			path: "$.key1.1.nestedKey",
			want: "nestedValue",
		},
		{
			name: "ComplexPath",
			data: jman.Obj{
				"key1": jman.Obj{
					"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
				},
			},
			path: "$.key1.nestedKey1.1.deepKey",
			want: "deepValue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := tc.data.Get(t, tc.path)
			assert.Equal(t, tc.want, value)
		})
	}
}

func TestObj_Get_PathNotFound(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	assertFatalf(t, "key 'key2' not found in object", func(mt jman.T) {
		_ = data.Get(mt, "$.key2")
	})
}

func TestObj_Get_IncorrectConversion(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	assertFatalf(t, "expected number at path '$.key1', got string", func(mt jman.T) {
		_ = data.GetNumber(mt, "$.key1")
	})
}
