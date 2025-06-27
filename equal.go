package jman

import (
	"fmt"
	"reflect"
)

func Equal(expected, actual any, optFuncs ...OptsFunc) error {
	opts := EqualOptions{}

	for _, o := range optFuncs {
		o(&opts)
	}

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

		return handleArrays(expectedSlice, actualSlice, opts)
	}

	if expectedObj, ok := NewObjFromAny(exp); ok {
		actualObj, ok := NewObjFromAny(act)
		if !ok {
			return fmt.Errorf("expected a json object, got %T", act)
		}
		return handleObjects(expectedObj, actualObj, opts)
	}

	return fmt.Errorf("unexpected type for expected %+v.  Use a json object/array string/byte, or jman.Obj or jman.Arr", expected)
}

func handleObjects(expected, actual Obj, opts EqualOptions) error {
	differences := compareObjects(expected, actual, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func handleArrays(expected, actual Arr, opts EqualOptions) error {
	differences := compareArrays(expected, actual, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func compareObjects(expected, actual Obj, opts EqualOptions) Differences {
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

		equal, diff := compareValues(expectedValue, actualValue, opts)
		if equal {
			continue
		}
		if diff.prefix == "" {
			diff.prefix = Both
		}
		diff.path = key
		diffs = append(diffs, diff)
	}

	return diffs
}

func compareArrays(expected, actual Arr, opts EqualOptions) Differences {
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
		equal, diff := compareValues(item, actual[i], opts)
		if equal {
			continue
		}
		if diff.prefix == "" {
			diff.prefix = Both
		}
		diff.path = fmt.Sprintf("%d", i)
		diffs = append(diffs, diff)
	}
	return diffs
}

func compareValues(expected, actual any, opts EqualOptions) (bool, Difference) {
	var (
		diff  Difference
		equal = true
	)
	switch expectedTyped := expected.(type) {
	case nil:
		if actual != nil {
			diff.diff = fmt.Sprintf("expected nil - got %T (%v)", actual, actual)
			equal = false
		}
	case bool:
		if err := compareTyped(expectedTyped, actual); err != nil {
			diff.diff = err.Error()
			equal = false
		}
	case float64:
		if err := compareTyped(expectedTyped, actual); err != nil {
			diff.diff = err.Error()
			equal = false
		}
	case string:
		// matcher placeholders have to be strings, so we only need to search them here
		matcher, found := opts.matchers.FindByPlaceholder(expectedTyped)
		if found {
			if !matcher.MatcherFunc(actual) {
				diff.diff = fmt.Sprintf("expected value for placeholder %q does not match actual value %v", expectedTyped, actual)
				equal = false
			}
			break
		}
		if err := compareTyped(expectedTyped, actual); err != nil {
			diff.diff = err.Error()
			equal = false
		}
	case []any:
		expectedArr := Arr(expectedTyped)
		actualTyped, ok := actual.([]any)
		if !ok {
			diff.diff = fmt.Sprintf("expected array - got %T (%v)", actual, actual)
			equal = false
			break
		}
		actualArr := Arr(actualTyped)
		diffs := compareArrays(expectedArr, actualArr, opts)
		if len(diffs) > 0 {
			diff.diff = diffs.Report()
			equal = false
		}
	case map[string]any:
		expectedObj := Obj(expectedTyped)
		actualTyped, ok := actual.(map[string]any)
		if !ok {
			diff.diff = fmt.Sprintf("expected object - got %T (%v)", actual, actual)
			equal = false
			break
		}
		actualObj := Obj(actualTyped)
		diffs := compareObjects(expectedObj, actualObj, opts)
		if len(diffs) > 0 {
			diff.diff = diffs.Report()
			equal = false
		}
	default:
		diff.prefix = Expected
		diff.diff = fmt.Sprintf("unsupported type comparison for expected value %q", expectedTyped)
		equal = false
	}

	return equal, diff
}

func compareTyped[T any](expected T, actual any) error {
	actualType, ok := actual.(T)
	if !ok || !reflect.DeepEqual(expected, actualType) {
		return fmt.Errorf(`expected "%v" - actual "%v"`, expected, actual)
	}

	return nil
}
