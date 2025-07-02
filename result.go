package jman

import "fmt"

// Result represents the result of a JSON Get operation from either an Obj or Arr.
// It contains the data retrieved and any error that occurred during the operation.
type Result struct {
	data any
	err  error
}

// Data returns the data retrieved from the JSON Get operation as any
func (r Result) Data() any {
	return r.data
}

// Err returns the error that occurred during the JSON Get operation, if any.
func (r Result) Error() error {
	return r.err
}

// IsError checks if the Result contains an error.
func (r Result) IsError() bool {
	return r.err != nil
}

// String returns the data as a string if it is of type string, otherwise it returns an error.
func (r Result) String() (string, error) {
	if r.err != nil {
		return "", r.err
	}

	if str, ok := r.data.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("result is not a string, got %T", r.data)
}

// Number returns the data as a float64 if it is of type float64, otherwise it returns an error.
func (r Result) Number() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if num, ok := r.data.(float64); ok {
		return num, nil
	}

	return 0, fmt.Errorf("result is not a number, got %T", r.data)
}

// Bool returns the data as a boolean if it is of type bool, otherwise it returns an error.
func (r Result) Bool() (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	if b, ok := r.data.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("result is not a boolean, got %T", r.data)
}

// Obj returns the data as an Obj if it is of type Obj, otherwise it returns an error.
func (r Result) Obj() (Obj, error) {
	if r.err != nil {
		return nil, r.err
	}

	if obj, ok := r.data.(Obj); ok {
		return obj, nil
	}

	return nil, fmt.Errorf("result is not an object, got %T", r.data)
}

// Arr returns the data as an Arr if it is of type Arr, otherwise it returns an error.
func (r Result) Arr() (Arr, error) {
	if r.err != nil {
		return nil, r.err
	}

	if arr, ok := r.data.(Arr); ok {
		return arr, nil
	}

	return nil, fmt.Errorf("result is not an array, got %T", r.data)
}
