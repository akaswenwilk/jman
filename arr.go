package jman

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type Arr []any

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

func (a Arr) Equal(other any, optFuncs ...OptsFunc) error {
	opts := EqualOptions{}

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
		return fmt.Errorf("expected not equal to actual:\n%s", differences.Report())
	}

	return nil
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

func (a Arr) Get(path string) Result {
	if path == "" {
		return Result{err: fmt.Errorf("path cannot be empty")}
	}

	paths := strings.Split(path, ".")
	if paths[0] != "$" {
		return Result{err: fmt.Errorf("path must start with '$', got '%s'", paths[0])}
	}

	a, err := normalize(a)
	if err != nil {
		return Result{err: fmt.Errorf("expected is invalid json: %w", err)}
	}

	val, err := getValue(paths[1:], a)
	return Result{
		data: val,
		err:  err,
	}
}
