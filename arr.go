package jman

import (
	"encoding/json"
	"fmt"
	"slices"
)

// Arr represents a JSON array. It implements the Equaler interface for deep equality checks.
type Arr []any

// UnmarshalJSON implements the json.Unmarshaler interface for Arr. Any items in the array of type []any or map[string]any
// will be converted to Arr or Obj, and other simple types will be converted to their
// corresponding types (bool, string, float64).
func (a *Arr) UnmarshalJSON(data []byte) error {
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for _, v := range raw {
		*a = append(*a, convert(v))
	}
	return nil
}

// Equal checks if the Arr is equal to another value, which can be either a JSON string, byte slice, or another Arr.
// It uses the provided options to customize the equality check, such as ignoring array order.
// If the actual value is not a valid JSON array, it returns an error.
func (a Arr) Equal(other any, optFuncs ...optsFunc) error {
	opts := equalOptions{}

	for _, o := range optFuncs {
		o(&opts)
	}

	if err := opts.valid(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	act, err := New[Arr](other)
	if err != nil {
		return fmt.Errorf("invalid actual: %w", err)
	}

	a, err = normalize(a)
	if err != nil {
		return fmt.Errorf("expected is invalid json: %w", err)
	}

	differences := compareArrays(base, a, act, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.report())
	}

	return nil
}

func compareArrays(path string, expected, actual Arr, opts equalOptions) differences {
	var diffs differences
	if len(expected) != len(actual) {
		diffs = append(diffs, difference{
			path: path,
			diff: fmt.Sprintf("expected %d items - got %d items", len(expected), len(actual)),
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

func compareArraysIgnoreOrder(path string, expected, actual Arr, opts equalOptions) differences {
	var diffs differences
	for i, item := range expected {
		found := slices.ContainsFunc(actual, func(v any) bool {
			equal, _ := compareValues(fmt.Sprintf("%s.%d", path, i), item, v, opts)
			return equal
		})
		if found {
			continue
		}

		d := difference{
			path: fmt.Sprintf("%s.%d", path, i),
			diff: "not found in actual",
		}
		diffs = append(diffs, d)
	}
	return diffs
}

func compareArraysStrictOrder(path string, expected, actual Arr, opts equalOptions) differences {
	var diffs differences
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

// Get retrieves a value from the Arr at the specified path.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be retrieved, it returns a Result with an error.
// If the path is valid, it returns a Result containing the value.
func (a Arr) Get(path string) Result {
	a, err := normalize(a)
	if err != nil {
		return Result{err: fmt.Errorf("expected is invalid json: %w", err)}
	}

	val, err := getValue(path, a)
	return Result{
		data: val,
		err:  err,
	}
}

// MustString returns the JSON representation of the Arr as a string.
// It panics if there is an error during marshaling.
func (a Arr) MustString() string {
	data, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return string(data)
}

// MustBytes returns the JSON representation of the Arr as a byte slice.
// It panics if there is an error during marshaling.
func (a Arr) MustBytes() []byte {
	data, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return data
}

// Set sets a value at the specified path in the Arr.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be set, it returns an error.
// After setting the value, it normalizes the Arr to ensure it is valid JSON.
// If normalization fails, it returns an error indicating the invalid JSON array.
// If successful, it updates the Arr in place.
// The value can be of any type, but it will be normalized to one of the supported JSON types:
// bool, string, float64, Obj, or Arr.
func (a Arr) Set(path string, value any) error {
	if err := setByPath(a, path, value); err != nil {
		return err
	}

	normalA, err := normalize(a)
	if err != nil {
		return fmt.Errorf("invalid json array: %w", err)
	}
	copy(a, normalA)

	return nil
}
