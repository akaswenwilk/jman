package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
)

func TestEqual_Strings(t *testing.T) {
	expected := `{"hello": "world", "foo": "bar", "baz": 123}`
	actual := `{"foo": "bar", "hello": "world", "baz": 123}`

	jman.Equal(t, expected, actual)
}

func TestEqual_Bytes(t *testing.T) {
}

func TestEqual_Obj(t *testing.T) {}

func TestEqual_Mix(t *testing.T) {}

func TestNew_UnexpectedType(t *testing.T) {}
