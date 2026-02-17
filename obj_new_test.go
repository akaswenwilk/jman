package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

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
	data := 12345 // unsupported type
	assertFatalf(t, "int unsupported type for JSON data, use either string or []byte", func(mt jman.T) {
		jman.New[jman.Obj](mt, data)
	})
}

func TestNewObj_InvalidJSON(t *testing.T) {
	data := `{"key1": "value1", "key2": 2, "key3": true,` // invalid JSON
	assertFatalf(t, `error parsing JSON data {"key1": "value1", "key2": 2, "key3": true,: unexpected end of JSON input`, func(mt jman.T) {
		jman.New[jman.Obj](mt, data)
	})
}
