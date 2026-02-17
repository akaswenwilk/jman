package jman_test

import (
	"testing"

	"github.com/akaswenwilk/jman"
)

func TestEqual_ObjExpected_ActualString(t *testing.T) {
	expected := jman.Obj{
		"foo": "bar",
		"num": 42,
	}
	actual := `{"foo":"bar","num":42}`

	jman.Equal(t, expected, actual)
}

func TestEqual_ArrExpected_ActualBytes_IgnoreOrder(t *testing.T) {
	expected := jman.Arr{"a", "b", "c"}
	actual := []byte(`["c","a","b"]`)

	jman.Equal(t, expected, actual, jman.WithIgnoreArrayOrder("$"))
}

func TestEqual_InstantiableObjInputs_MapAndStruct(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	expected := map[string]any{
		"name": "alice",
		"age":  30,
	}
	actual := person{
		Name: "alice",
		Age:  30,
	}

	jman.Equal(t, expected, actual)
}

func TestEqual_InstantiableArrInputs_Slices(t *testing.T) {
	expected := []int{1, 2, 3}
	actual := []any{1, 2, 3}

	jman.Equal(t, expected, actual)
}

func TestEqual_UnequalKinds_ExpectedObjActualArr(t *testing.T) {
	expected := jman.Obj{"foo": "bar"}
	actual := jman.Arr{"foo", "bar"}

	assertFatalf(t, "can't compare json object with array", func(mt jman.T) {
		jman.Equal(mt, expected, actual)
	})
}

func TestEqual_UnequalKinds_ExpectedArrActualObj(t *testing.T) {
	expected := jman.Arr{"foo", "bar"}
	actual := jman.Obj{"foo": "bar"}

	assertFatalf(t, "can't compare array with json object", func(mt jman.T) {
		jman.Equal(mt, expected, actual)
	})
}

func TestEqual_UnsupportedType(t *testing.T) {
	assertFatalf(t, "int unsupported type for JSON data, use either string or []byte", func(mt jman.T) {
		jman.Equal(mt, 123, jman.Obj{"foo": "bar"})
	})
}
