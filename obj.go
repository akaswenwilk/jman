package jman

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
)

// Obj represents a JSON object. It implements the Equaler interface for deep equality checks.
type Obj map[string]any

// UnmarshalJSON implements the json.Unmarshaler interface for Obj. Any nested fields with values of type []any or map[string]any
// will be converted to Arr or Obj, and other simple types will be converted to their
// corresponding types (bool, string, float64).
func (o *Obj) UnmarshalJSON(data []byte) error {
	raw := map[string]any{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*o = make(Obj)
	for k, v := range raw {
		(*o)[k] = convert(v)
	}
	return nil
}

func convert(v any) any {
	switch vv := v.(type) {
	case map[string]any:
		sub := Obj{}
		for k, val := range vv {
			sub[k] = convert(val)
		}
		return sub
	case []any:
		sub := make(Arr, len(vv))
		for i, el := range vv {
			sub[i] = convert(el)
		}
		return sub
	default:
		return v
	}
}

// Equal checks if the Obj is equal to another value, which can be either a JSON string, byte slice, or another Obj.
func (ob Obj) Equal(other any, optFuncs ...optsFunc) error {
	opts := equalOptions{}

	for _, o := range optFuncs {
		o(&opts)
	}

	if err := opts.valid(); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}

	act, err := New[Obj](other)
	if err != nil {
		if errors.Is(err, ErrUnsupportedType) {
			return fmt.Errorf("can't compare json object with %T", other)
		}
		return fmt.Errorf("invalid actual: %w", err)
	}

	ob, err = normalize(ob)
	if err != nil {
		return fmt.Errorf("expected is invalid json: %w", err)
	}

	differences := compareObjects(base, ob, act, opts)
	if len(differences) > 0 {
		return fmt.Errorf("expected not equal to actual:\n%s", differences.report())
	}

	return nil
}

func compareObjects(path string, expected, actual Obj, opts equalOptions) differences {
	var diffs differences
	for k := range maps.Keys(expected) {
		_, exists := actual[k]
		if !exists {
			diffs = append(diffs, difference{
				diff: "not found in actual",
				path: pathAndKey(path, k),
			})
		}
	}

	for k := range maps.Keys(actual) {
		_, exists := expected[k]
		if !exists {
			diffs = append(diffs, difference{
				diff: "unexpected key",
				path: pathAndKey(path, k),
			})
		}
	}

	for key, expectedValue := range expected {
		// we know that this key is not present on actual
		// so we can skip
		if diffs.hasPath(pathAndKey(path, key)) {
			continue
		}

		actualValue := actual[key]

		equal, diff := compareValues(pathAndKey(path, key), expectedValue, actualValue, opts)
		if equal {
			continue
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func pathAndKey(path, key string) string {
	return fmt.Sprintf("%s.%s", path, key)
}

// Get retrieves a value from the Obj at the specified path.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be retrieved, it returns a Result with an error.
// If the path is valid, it returns a Result containing the value.
func (ob Obj) Get(path string) Result {
	ob, err := normalize(ob)
	if err != nil {
		return Result{err: fmt.Errorf("invalid json object: %w", err)}
	}

	val, err := getValue(path, ob)
	return Result{
		data: val,
		err:  err,
	}
}

// MustString returns the JSON representation of the Obj as a string.
// It panics if the marshaling fails.
func (ob Obj) MustString() string {
	data, err := json.Marshal(ob)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return string(data)
}

// MustBytes returns the JSON representation of the Obj as a byte slice.
// It panics if the marshaling fails.
func (ob Obj) MustBytes() []byte {
	data, err := json.Marshal(ob)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return data
}

// Set sets a value at the specified path in the Obj.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be set, it returns an error.
// After setting the value, it normalizes the Arr to ensure it is valid JSON.
// If normalization fails, it returns an error indicating the invalid JSON array.
// If successful, it updates the Arr in place.
// The value can be of any type, but it will be normalized to one of the supported JSON types:
// bool, string, float64, Obj, or Arr.
func (o Obj) Set(path string, value any) error {
	if err := setByPath(o, path, value); err != nil {
		return err
	}
	normalizedO, err := normalize(o)
	if err != nil {
		return fmt.Errorf("invalid json object: %w", err)
	}
	maps.Copy(o, normalizedO)
	return nil
}
