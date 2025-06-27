package jman

import (
	"fmt"
	"reflect"
)

func Equal(expected, actual any) error {
	exp, err := New(expected)
	if err != nil {
		return fmt.Errorf("invalid expected: %w", err)
	}

	act, err := New(actual)
	if err != nil {
		return fmt.Errorf("invalid actual: %w", err)
	}

	if expectedSlice, ok := NewArrFromAny(exp); ok {
		actualSlice, ok := NewArrFromAny(act)
		if !ok {
			return fmt.Errorf("expected a json array, got %T", act)
		}

		return handleArrays(expectedSlice, actualSlice)
	}

	if expectedObj, ok := NewObjFromAny(exp); ok {
		actualObj, ok := NewObjFromAny(act)
		if !ok {
			return fmt.Errorf("expected a json object, got %T", act)
		}
		return handleObjects(expectedObj, actualObj)
	}

	return fmt.Errorf("unexpected type for expected %+v.  Use a json object/array string/byte, or jman.Obj or jman.Arr", expected)
}

func handleObjects(expected, actual Obj) error {
	differences := compareObjects(expected, actual)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func handleArrays(expected, actual Arr) error {
	differences := compareArrays(expected, actual)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func compareObjects(expected, actual Obj) Differences {
	var diffs Differences
	for _, k := range expected.Keys() {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, Difference{
				prefix: Expected,
				diff:   "missing key",
				path:   k,
			})
		}
	}

	for _, k := range actual.Keys() {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, Difference{
				prefix: Actual,
				diff:   "unexpected key",
				path:   k,
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
				prefix: Both,
				diff:   err.Error(),
				path:   key,
			})
		}
	}

	return diffs
}

func compareArrays(expected, actual Arr) Differences {
	var diffs Differences
	if len(expected) != len(actual) {
		diffs = append(diffs, Difference{
			prefix: Expected,
			diff:   fmt.Sprintf("expected %d items - got %d items", len(expected), len(actual)),
		})
	}

	for i, item := range expected {
		if i >= len(actual) {
			continue
		}
		if err := compareValues(item, actual[i]); err != nil {
			diffs = append(diffs, Difference{
				path:   fmt.Sprintf("%d", i),
				diff:   err.Error(),
				prefix: Both,
			})
		}
	}
	return diffs
}

func compareValues(expected, actual any) error {
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected %q - actual %q", expected, actual)
	}
	return nil
}
