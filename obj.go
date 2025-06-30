package jman

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"
)

type Obj map[string]any

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

func (ob Obj) Equal(other any, optFuncs ...OptsFunc) error {
	opts := EqualOptions{}

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
				diff:   "not found in actual",
				path:   pathAndKey(path, k),
			})
		}
	}

	for k := range maps.Keys(actual) {
		_, exists := expected[k]
		if !exists {
			diffs = append(diffs, Difference{
				diff:   "unexpected key",
				path:   pathAndKey(path, k),
			})
		}
	}

	for key, expectedValue := range expected {
		// we know that this key is not present on actual
		// so we can skip
		if diffs.HasPath(pathAndKey(path, key)) {
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

func (ob Obj) Get(path string) Result {
	if path == "" {
		return Result{err: fmt.Errorf("path cannot be empty")}
	}
	paths := strings.Split(path, ".")
	if paths[0] != "$" {
		return Result{err: fmt.Errorf("path must start with '$', got '%s'", paths[0])}
	}

	ob, err := normalize(ob)
	if err != nil {
		return Result{err: fmt.Errorf("invalid json object: %w", err)}
	}

	val, err := getValue(paths[1:], ob)
	return Result{
		data: val,
		err:  err,
	}
}

func (ob Obj) MustString() string {
	data, err := json.Marshal(ob)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return string(data)
}

func (ob Obj) MustBytes() []byte {
	data, err := json.Marshal(ob)
	if err != nil {
		panic(fmt.Sprintf("error marshaling JSON object: %v", err))
	}
	return data
}

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
