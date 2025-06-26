package jman

import (
	"fmt"
	"reflect"
	"testing"
)

type TB interface {
	Errorf(format string, args ...any)
}

func Equal(t testing.TB, expected, actual any) {
	expectedObj, err := New(expected)
	if err != nil {
		t.Errorf("invalid expected: %s", err.Error())
		return
	}

	actualObj, err := New(actual)
	if err != nil {
		t.Errorf("invalid actual: %s", err.Error())
		return
	}

	differences := compare(expectedObj, actualObj)
	if len(differences) > 0 {
		t.Errorf("expected not equal to actual: %s", differences.Report())
	}
}

type Differences []Difference

func (d Differences) Report() string {
	var report string
	for _, diff := range d {
		report += fmt.Sprintf("%s.%s %s", diff.expectedOrActual, diff.path, diff.diff)
	}
	return report
}

func (d Differences) HasKey(key string) bool {
	for _, d := range d {
		if d.path == key {
			return true
		}
	}

	return false
}

type Difference struct {
	expectedOrActual string
	diff             string
	path             string
}

func compare(expected, actual Obj) Differences {
	var diffs Differences
	for _, k := range expected.Keys() {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, Difference{
				expectedOrActual: "expected",
				diff:             "missing key",
				path:             k,
			})
		}
	}

	for _, k := range actual.Keys() {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, Difference{
				expectedOrActual: "actual",
				diff:             "unexpected key",
				path:             k,
			})
		}
	}

	for key, expectedValue := range expected {
		// we know that this key is not present on actual
		// so we can skip
		if diffs.HasKey(key) {
			continue
		}

		actualValue := actual[key]

		if err := compareValues(expectedValue, actualValue); err != nil {
			diffs = append(diffs, Difference{
				expectedOrActual: "both",
				diff:             err.Error(),
				path:             key,
			})
		}
	}

	return diffs
}

func compareValues(expected, actual any) error {
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected %v not equal to actual %v", expected, actual)
	}
	return nil
}
