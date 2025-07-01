package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
)

func TestArr_Equal_IgnoreOrder_Base(t *testing.T) {
	expected := jman.Arr{"a", "b", "c"}
	actual := jman.Arr{"c", "b", "a"}

	assert.NoError(t, expected.Equal(actual, jman.WithIgnoreArrayOrder("$")))
}
