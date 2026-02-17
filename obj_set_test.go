package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

func TestObj_Set(t *testing.T) {
	testCases := []struct {
		name  string
		data  jman.Obj
		path  string
		value any
		want  jman.Obj
	}{
		{
			name:  "NewKey",
			data:  jman.Obj{"key1": "value1"},
			path:  "$.key2",
			value: "newValue",
			want: jman.Obj{
				"key1": "value1",
				"key2": "newValue",
			},
		},
		{
			name:  "ExistingKey",
			data:  jman.Obj{"key1": "value1"},
			path:  "$.key1",
			value: "updatedValue",
			want: jman.Obj{
				"key1": "updatedValue",
			},
		},
		{
			name: "NestedKey",
			data: jman.Obj{
				"key1": jman.Obj{
					"nestedKey": "nestedValue",
				},
			},
			path:  "$.key1.nestedKey",
			value: "newNestedValue",
			want: jman.Obj{
				"key1": jman.Obj{
					"nestedKey": "newNestedValue",
				},
			},
		},
		{
			name: "ComplexPath",
			data: jman.Obj{
				"key1": jman.Obj{
					"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
				},
			},
			path:  "$.key1.nestedKey1.1.deepKey",
			value: "newDeepValue",
			want: jman.Obj{
				"key1": jman.Obj{
					"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "newDeepValue"}},
				},
			},
		},
		{
			name:  "Obj",
			data:  jman.Obj{},
			path:  "$.key1",
			value: jman.Obj{"nestedkey": "nestedValue"},
			want: jman.Obj{
				"key1": jman.Obj{"nestedkey": "nestedValue"},
			},
		},
		{
			name:  "GenericMap",
			data:  jman.Obj{},
			path:  "$.key1",
			value: map[string]any{"nestedkey": "nestedValue"},
			want: jman.Obj{
				"key1": jman.Obj{"nestedkey": "nestedValue"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.data
			data.Set(t, tc.path, tc.value)
			assert.Equal(t, tc.want, data)
		})
	}
}
