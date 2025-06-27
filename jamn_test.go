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

func TestEqual_NestedObj(t *testing.T) {}

func TestEqual_SliceSimple(t *testing.T) {}

func TestEqual_SliceObjects(t *testing.T) {}

func TestNew_UnexpectedType(t *testing.T) {}
