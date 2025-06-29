package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObj_Get_String(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	value, err := jman.Get[string]("$.key1", data)
	require.NoError(t, err)
	assert.Equal(t, "value1", value)
}

func TestObj_Get_Int(t *testing.T) {
	data := jman.Obj{"key1": 123}
	value, err := jman.Get[int]("$.key1", data)
	require.NoError(t, err)
	assert.Equal(t, 123, value)
}

func TestObj_Get_Bool(t *testing.T) {
	data := jman.Obj{"key1": true}
	value, err := jman.Get[bool]("$.key1", data)
	require.NoError(t, err)
	assert.Equal(t, true, value)
}

func TestObj_Get_Obj(t *testing.T) {
	data := jman.Obj{"key1": jman.Obj{"nestedKey": "nestedValue"}}
	value, err := jman.Get[jman.Obj]("$.key1", data)
	require.NoError(t, err)
	assert.Equal(t, jman.Obj{"nestedKey": "nestedValue"}, value)
}

func TestObj_Get_ObjFromString(t *testing.T) {
	data := `{"key1": {"nestedKey": "nestedValue"}}`
	obj, err := jman.New[jman.Obj](data)
	require.NoError(t, err)
	value, err := jman.Get[jman.Obj]("$.key1", obj)
	require.NoError(t, err)
	assert.Equal(t, jman.Obj{"nestedKey": "nestedValue"}, value)
}

func TestObj_Get_Arr(t *testing.T) {
	data := jman.Obj{"key1": jman.Arr{"item1", "item2"}}
	value, err := jman.Get[jman.Arr]("$.key1", data)
	require.NoError(t, err)
	assert.Equal(t, jman.Arr{"item1", "item2"}, value)
}

func TestObj_Get_NestedKey(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey": "nestedValue",
		},
	}
	value, err := jman.Get[string]("$.key1.nestedKey", data)
	require.NoError(t, err)
	assert.Equal(t, "nestedValue", value)
}

func TestObj_Get_NestedKeyJSONString(t *testing.T) {
	data := `{"key1": {"nestedKey": "nestedValue"}}`
	obj, err := jman.New[jman.Obj](data)
	require.NoError(t, err)
	value, err := jman.Get[string]("$.key1.nestedKey", obj)
	require.NoError(t, err)
	assert.Equal(t, "nestedValue", value)
}

func TestObj_Get_ArrItem(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Arr{"item1", "item2", "item3"},
	}
	value, err := jman.Get[string]("$.key1.1", data)
	require.NoError(t, err)
	assert.Equal(t, "item2", value)
}

func TestObj_Get_ComplexPath(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	value, err := jman.Get[string]("$.key1.nestedKey1.1.deepKey", data)
	require.NoError(t, err)
	assert.Equal(t, "deepValue", value)
}

func TestObj_Get_PathNotFound(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	_, err := jman.Get[string]("$.key2", data)
	require.Error(t, err)
	assert.Equal(t, "key 'key2' not found in object", err.Error())
}

func TestObj_Get_ComplexPathNotFound(t *testing.T) {
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	_, err := jman.Get[string]("$.key1.nestedKey1.2.deepKey", data)
	require.Error(t, err)
	assert.Equal(t, "index 2 out of bounds for array of length 2", err.Error())
}

func TestObj_Get_IncorrectConversion(t *testing.T) {
	data := jman.Obj{"key1": "value1"}
	_, err := jman.Get[int]("$.key1", data)
	require.Error(t, err)
	assert.Equal(t, "value at path '$.key1' is not of type int", err.Error())
}
