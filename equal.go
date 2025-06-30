package jman

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrJSONParse = errors.New("error parsing JSON data")
	ErrNormalize = errors.New("error normalizing JSON data")
	ErrUnsupportedType = errors.New("unsupported type for JSON data, use either string or []byte")
)

type JSONEqual interface {
	Equal(other any, optFuncs ...OptsFunc) error
}

func New[T JSONEqual](data any) (T, error) {
	var result T

	switch data.(type) {
	case string:
		d := data.(string)
		if err := json.Unmarshal([]byte(d), &result); err != nil {
			return result, fmt.Errorf("%w %s: %w", ErrJSONParse, data, err)
		}
	case []byte:
		d := data.([]byte)
		if err := json.Unmarshal(d, &result); err != nil {
			return result, fmt.Errorf("%w %s: %w", ErrJSONParse, string(d), err)
		}
	case T:
		// If the data is already of type T, we can normalize it and return it directly
		var err error
		result, err = normalize(data.(T))
		if err != nil {
			return result, fmt.Errorf("%w: %w", ErrNormalize, err)
		}
	default:
		return result, fmt.Errorf("%T %w", data, ErrUnsupportedType)
	}

	return result, nil
}

func normalize[T JSONEqual](data T) (T, error) {
	// Marshal and unmarshal to normalize the JSON structure
	marshaled, err := json.Marshal(data)
	if err != nil {
		return data, fmt.Errorf("error marshalling JSON object: %w", err)
	}

	var normalized T
	if err := json.Unmarshal(marshaled, &normalized); err != nil {
		return data, fmt.Errorf("error unmarshalling JSON object: %w", err)
	}

	return normalized, nil
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
	case Arr:
		actualTyped, ok := actual.(Arr)
		if !ok {
			diff.diff = fmt.Sprintf("expected array - got %T (%v)", actual, actual)
			equal = false
			break
		}
		diffs := compareArrays(path, expectedTyped, actualTyped, opts)
		if len(diffs) > 0 {
			diff.subDiffs = diffs
			equal = false
		}
	case Obj:
		actualTyped, ok := actual.(Obj)
		if !ok {
			diff.diff = fmt.Sprintf("expected object - got %T (%v)", actual, actual)
			equal = false
			break
		}
		diffs := compareObjects(path, expectedTyped, actualTyped, opts)
		if len(diffs) > 0 {
			diff.subDiffs = diffs
			equal = false
		}
	default:
		diff.path = path
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
