package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
)

func TestArr_Equal_IgnoreOrder_Base(t *testing.T) {
	expected := jman.Arr{"a", "b", "c"}
	actual := jman.Arr{"c", "b", "a"}

	expected.Equal(t, actual, jman.WithIgnoreArrayOrder("$"))
}
