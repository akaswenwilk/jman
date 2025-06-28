package jman

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
)

const base = "$"

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

	if exp.IsArr() {
		if !act.IsArr() {
			return fmt.Errorf("expected a json array, got %s", act.UnderlyingType())
		}

		expArr, _ := exp.Arr()
		actArr, _ := act.Arr()
		return handleArrays(expArr, actArr, opts)
	}

	if exp.IsObj() {
		if !act.IsObj() {
			return fmt.Errorf("expected a json object, got %s", act.UnderlyingType())
		}

		expObj, _ := exp.Obj()
		actObj, _ := act.Obj()
		return handleObjects(expObj, actObj, opts)
	}

	return fmt.Errorf("unexpected type for expected %+v.  Use a json object/array string/byte, or jman.Obj or jman.Arr", expected)
}

func handleObjects(expected, actual Obj, opts EqualOptions) error {
	differences := compareObjects(base, expected, actual, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func handleArrays(expected, actual Arr, opts EqualOptions) error {
	differences := compareArrays(base, expected, actual, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
}

func compareObjects(path string, expected, actual Obj, opts EqualOptions) Differences {
	var diffs Differences
	for k := range maps.Keys(expected) {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, Difference{
				prefix: Expected,
				diff:   "missing key",
				path:   k,
			})
		}
	}

	for k := range maps.Keys(actual) {
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

		equal, diff := compareValues(fmt.Sprintf("%s.%s", path, key), expectedValue, actualValue, opts)
		if equal {
			continue
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func compareArrays(path string, expected, actual Arr, opts EqualOptions) Differences {
	var diffs Differences
	if len(expected) != len(actual) {
		diffs = append(diffs, Difference{
			prefix: Expected,
			diff:   fmt.Sprintf("expected %d items - got %d items", len(expected), len(actual)),
		})
	}

	if slices.Contains(opts.ignoreArrayOrder, path) {
		comparedDiffs := compareArraysIgnoreOrder(path, expected, actual, opts)
		diffs = append(diffs, comparedDiffs...)
	} else {
		comparedDiffs := compareArraysStrictOrder(path, expected, actual, opts)
		diffs = append(diffs, comparedDiffs...)
	}
	return diffs
}

func compareArraysIgnoreOrder(path string, expected, actual Arr, opts EqualOptions) Differences {
	var diffs Differences
	for i, item := range expected {
		found := slices.ContainsFunc(actual, func(v any) bool {
			equal, _ := compareValues(fmt.Sprintf("%s.%d", path, i), item, v, opts)
			return equal
		})
		if found {
			continue
		}

		d := Difference{
			prefix: Expected,
			path:   fmt.Sprintf("%s.%d", path, i),
			diff:   "not found in actual",
		}
		diffs = append(diffs, d)
	}
	return diffs
}

func compareArraysStrictOrder(path string, expected, actual Arr, opts EqualOptions) Differences {
	var diffs Differences
	for i, item := range expected {
		if i >= len(actual) {
			continue
		}
		equal, diff := compareValues(fmt.Sprintf("%s.%d", path, i), item, actual[i], opts)
		if equal {
			continue
		}
		diffs = append(diffs, diff)
	}
	return diffs
}

func compareValues(path string, expected, actual any, opts EqualOptions) (bool, Difference) {
	var (
		diff = Difference{
			path: path,
		}
		equal = true
	)
	switch expectedTyped := expected.(type) {
	case nil:
		if actual != nil {
			diff.diff = unequalMessage(expectedTyped, actual)
			equal = false
		}
	case bool, float64:
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
		diffs := compareArrays(path, expectedArr, actualArr, opts)
		if len(diffs) > 0 {
			diff.subDiffs = diffs
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
		diffs := compareObjects(path, expectedObj, actualObj, opts)
		if len(diffs) > 0 {
			diff.subDiffs = diffs
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
		return errors.New(unequalMessage(expected, actual))
	}

	return nil
}

func unequalMessage(expected, actual any) string {
	msg := `expected ` + formatterFor(expected) + ` - actual ` + formatterFor(actual)
	return fmt.Sprintf(msg, expected, actual)
}

func formatterFor(value any) string {
	if value == nil {
		return "%v"
	}
	if _, ok := value.(string); ok {
		return "%q"
	}
	if _, ok := value.(bool); ok {
		return "%t"
	}
	return "%v"
}
