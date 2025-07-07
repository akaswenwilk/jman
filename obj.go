package jman

import (
	"encoding/json"
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
func (ob Obj) Equal(t T, other any, optFuncs ...optsFunc) {
	opts := equalOptions{}

	for _, o := range optFuncs {
		o(&opts)
	}

	if err := opts.valid(); err != nil {
		t.Fatalf(fmt.Sprintf("invalid options: %v", err))
		return
	}

	// check if other value is an array
	switch v := other.(type) {
	case string:
		if v[0] == '[' {
			t.Fatalf("can't compare json object with array")
		}
	case []byte:
		if v[0] == '[' {
			t.Fatalf("can't compare json object with array")
		}
	case Arr:
		t.Fatalf("can't compare json object with array")
	}

	act := New[Obj](t, other)

	ob, err := normalize(ob)
	if err != nil {
		t.Fatalf(fmt.Sprintf("expected is invalid json: %v", err))
	}

	diffs := compareObjects(base, ob, act, opts)
	if len(diffs) > 0 {
		t.Fatalf(fmt.Sprintf("expected not equal to actual:\n%s", diffs.report()))
	}
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

// String returns the JSON representation of the Obj as a string.
// It fails if the marshaling fails.
func (ob Obj) String(t T) string {
	data, err := json.Marshal(ob)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return string(data)
}

// Bytes returns the JSON representation of the Obj as a byte slice.
// It fails if the marshaling fails.
func (ob Obj) Bytes(t T) []byte {
	data, err := json.Marshal(ob)
	if err != nil {
		t.Fatalf(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return data
}

// Get retrieves a value from the Obj at the specified path.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be retrieved, it calls t.Fatal
func (o Obj) Get(t T, path string) any {
	normed, err := normalize(o)
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
func (o Obj) GetString(t T, path string) string {
	val := o.Get(t, path)
	if str, ok := val.(string); ok {
		return str
	}
	t.Fatalf(fmt.Sprintf("expected string at path '%s', got %T", path, val))
	return ""
}

// GetNumber functions like Get but attempts to convert to float64. Fails if the value at the path is not a number.
func (o Obj) GetNumber(t T, path string) float64 {
	val := o.Get(t, path)
	if num, ok := val.(float64); ok {
		return num
	}
	t.Fatalf(fmt.Sprintf("expected number at path '%s', got %T", path, val))
	return 0
}

// GetBool functions like Get but attempts to convert to bool. Fails if the value at the path is not a boolean.
func (o Obj) GetBool(t T, path string) bool {
	val := o.Get(t, path)
	if b, ok := val.(bool); ok {
		return b
	}
	t.Fatalf(fmt.Sprintf("expected bool at path '%s', got %T", path, val))
	return false
}

// GetArray functions like Get but attempts to convert to Arr. Fails if the value at the path is not an array.
func (o Obj) GetArray(t T, path string) Arr {
	val := o.Get(t, path)
	if arr, ok := val.(Arr); ok {
		return arr
	}
	t.Fatalf(fmt.Sprintf("expected array at path '%s', got %T", path, val))
	return nil
}

// GetObject functions like Get but attempts to convert to Obj. Fails if the value at the path is not an object.
func (o Obj) GetObject(t T, path string) Obj {
	val := o.Get(t, path)
	if obj, ok := val.(Obj); ok {
		return obj
	}
	t.Fatalf(fmt.Sprintf("expected object at path '%s', got %T", path, val))
	return nil
}

// Set sets a value at the specified path in the Obj.
// The path is a JSONPath-like string that specifies the location of the value.
// If the path is invalid or the value cannot be set, it fails.
// After setting the value, it normalizes the Arr to ensure it is valid JSON.
// If normalization fails, it calls Fatalf on T.
// If successful, it updates the Arr in place.
// The value can be of any type, but it will be normalized to one of the supported JSON types:
// bool, string, float64, Obj, or Arr.
func (o Obj) Set(t T, path string, value any) {
	if err := setByPath(o, path, value); err != nil {
		t.Fatalf(err.Error())
	}

	normalA, err := normalize(o)
	if err != nil {
		t.Fatalf(fmt.Sprintf("invalid json array: %v", err))
	}
	maps.Copy(o, normalA)
}
