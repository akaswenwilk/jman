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
func (a Arr) Equal(t T, other any, optFuncs ...optsFunc) {
	opts := equalOptions{}

	for _, o := range optFuncs {
		o(&opts)
	}

	if err := opts.valid(); err != nil {
		t.Fatalf(fmt.Sprintf("invalid options: %v", err))
	}

	// check if other value is an array
	switch v := other.(type) {
	case string:
		if v[0] == '{' {
			t.Fatalf("can't compare array with json object")
		}
	case []byte:
		if v[0] == '{' {
			t.Fatalf("can't compare array with json object")
		}
	case Obj:
		t.Fatalf("can't compare array with json object")
	}

	act := New[Arr](t, other)

	a, err := normalize(a)
	if err != nil {
		t.Fatalf(fmt.Sprintf("expected is invalid json: %v", err))
	}

	diffs := compareArrays(base, a, act, opts)
	if len(diffs) > 0 {
		t.Fatalf(fmt.Sprintf("expected not equal to actual:\nexpected %s\nactual %s\n\n%s", a.String(t), act.String(t), diffs.report()))
	}
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

// String returns the JSON representation of the Arr as a string.
// It fails if there is an error during marshaling.
func (a Arr) String(t T) string {
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return string(data)
}

// Bytes returns the JSON representation of the Arr as a byte slice.
// It fails if there is an error during marshaling.
func (a Arr) MustBytes(t T) []byte {
	data, err := json.Marshal(a)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return data
}

// Get retrieves a value from the Arr at the specified path.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be retrieved, it calls t.Fatal
func (a Arr) Get(t T, path string) any {
	normed, err := normalize(a)
	if err != nil {
		t.Fatalf(fmt.Sprintf("err normalizing arr: %v", err))
	}
	val, err := getValue(path, normed)
	if err != nil {
		t.Fatalf(fmt.Sprintf("failed to get value at path '%s': %v", path, err))
	}

	return val
}

// GetString functions like Get but attempts to convert to string. Fails if the value at the path is not a string.
func (a Arr) GetString(t T, path string) string {
	val := a.Get(t, path)
	if str, ok := val.(string); ok {
		return str
	}
	t.Fatalf(fmt.Sprintf("expected string at path '%s', got %T", path, val))
	return ""
}

// GetNumber functions like Get but attempts to convert to float64. Fails if the value at the path is not a number.
func (a Arr) GetNumber(t T, path string) float64 {
	val := a.Get(t, path)
	if num, ok := val.(float64); ok {
		return num
	}
	t.Fatalf(fmt.Sprintf("expected number at path '%s', got %T", path, val))
	return 0
}

// GetBool functions like Get but attempts to convert to bool. Fails if the value at the path is not a boolean.
func (a Arr) GetBool(t T, path string) bool {
	val := a.Get(t, path)
	if b, ok := val.(bool); ok {
		return b
	}
	t.Fatalf(fmt.Sprintf("expected bool at path '%s', got %T", path, val))
	return false
}

// GetArray functions like Get but attempts to convert to Arr. Fails if the value at the path is not an array.
func (a Arr) GetArray(t T, path string) Arr {
	val := a.Get(t, path)
	if arr, ok := val.(Arr); ok {
		return arr
	}
	t.Fatalf(fmt.Sprintf("expected array at path '%s', got %T", path, val))
	return nil
}

// GetObject functions like Get but attempts to convert to Obj. Fails if the value at the path is not an object.
func (a Arr) GetObject(t T, path string) Obj {
	val := a.Get(t, path)
	if obj, ok := val.(Obj); ok {
		return obj
	}
	t.Fatalf(fmt.Sprintf("expected object at path '%s', got %T", path, val))
	return nil
}

// Set sets a value at the specified path in the Arr.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be set, it fails.
// After setting the value, it normalizes the Arr to ensure it is valid JSON.
// If normalization fails, it calls Fatalf on T.
// If successful, it updates the Arr in place.
// The value can be of any type, but it will be normalized to one of the supported JSON types:
// bool, string, float64, Obj, or Arr.
func (a Arr) Set(t T, path string, value any) {
	if err := setByPath(a, path, value); err != nil {
		t.Fatalf(err.Error())
	}

	normalA, err := normalize(a)
	if err != nil {
		t.Fatalf(fmt.Sprintf("invalid json array: %v", err))
	}
	copy(a, normalA)
}
