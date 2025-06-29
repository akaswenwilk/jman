package jman

import "fmt"

type Result struct {
	data any
	err  error
}

func (r Result) Data() any {
	return r.data
}

func (r Result) Error() error {
	return r.err
}

func (r Result) IsError() bool {
	return r.err != nil
}

func (r Result) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}

	if str, ok := r.data.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("result is not a string, got %T", r.data)
}

func (r Result) Number() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if num, ok := r.data.(float64); ok {
		return num, nil
	}

	return 0, fmt.Errorf("result is not a number, got %T", r.data)
}

func (r Result) Bool() (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	if b, ok := r.data.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("result is not a boolean, got %T", r.data)
}

func (r Result) Obj() (Obj, error) {
	if r.err != nil {
		return nil, r.err
	}

	if obj, ok := r.data.(Obj); ok {
		return obj, nil
	}

	return nil, fmt.Errorf("result is not an object, got %T", r.data)
}

func (r Result) Arr() (Arr, error) {
	if r.err != nil {
		return nil, r.err
	}

	if arr, ok := r.data.(Arr); ok {
		return arr, nil
	}

	return nil, fmt.Errorf("result is not an array, got %T", r.data)
}
