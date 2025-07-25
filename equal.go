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

// JSONEqual is an interface that defines a method for comparing JSON objects.
// It requires the implementation to provide an Equal method that checks interface
// equality between two JSON objects, allowing for custom options via OptsFunc.
type JSONEqual interface {
	Equal(t T, other any, optFuncs ...optsFunc)
}

// New creates a new instance of type T from the provided data.
// The data can be a JSON string, a byte slice, or an instance of type T.
// It returns an error if the data cannot be parsed or normalized into type T.
func New[E JSONEqual](t T, data any) E {
	var result E

	switch d := data.(type) {
	case string:
		if err := json.Unmarshal([]byte(d), &result); err != nil {
			t.Fatalf(fmt.Sprintf("%v %s: %v", ErrJSONParse, d, err))
		}
	case []byte:
		if err := json.Unmarshal(d, &result); err != nil {
			t.Fatalf(fmt.Sprintf("%v %s: %v", ErrJSONParse, string(d), err))
		}
	case E:
		// If the data is already of type T, we can normalize it and return it directly
		var err error
		result, err = normalize(d)
		if err != nil {
			t.Fatalf(fmt.Sprintf("%v %T: %v", ErrNormalize, d, err))
		}
	default:
		t.Fatalf(fmt.Sprintf("%T %s", data, ErrUnsupportedType))
	}

	return result
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

func compareValues(path string, expected, actual any, opts equalOptions) (bool, difference) {
	var (
		diff = difference{
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
