package jman




import (
	"reflect"
	"regexp"
)

// Matchers is a collection of Matcher objects.
type Matchers []Matcher

// FindByPlaceholder searches for a Matcher by its Placeholder.
func (m Matchers) FindByPlaceholder(p string) (Matcher, bool) {
	for _, matcher := range m {
		if matcher.Placeholder == p {
			return matcher, true
		}
	}

	return Matcher{}, false
}

// Matcher represents a matcher that can be used to validate values against certain conditions.
// It contains a placeholder for identification and a function that defines the matching logic.
// It will be checked whenever a string value is encountered in the expected value.
// If a matcher with the same placeholder is found, it will be used to validate the actual value.
type Matcher struct {
	Placeholder string
	MatcherFunc
}

// MatcherFunc is a function type that defines the matching logic.
// It takes a value of any type and returns true if the value matches the expected condition, false otherwise.
// the value passed into the MatcherFunc can be of any type, including nil and is taken from the actual value in the JSON.
type MatcherFunc func(v any) bool

// NotEmpty creates a matcher that checks if the value is not empty.
// a string must not be empty, cannot be nil.  booleans can be either true or false.
// a number can be any number, including zero.
// an array must not be empty, cannot be nil.  
// an object must not be empty, cannot be nil.
func NotEmpty(placeholder string) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: func(v any) bool {
			switch typed := v.(type) {
			case nil:
				return false
			case bool, float64:
				return true
			case string:
				return typed != ""
			case []any:
				return len(typed) > 0
			case map[string]any:
				return len(typed) > 0
			}

			return false
		},
	}
}

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

// IsUUID creates a matcher that checks if the value is a valid UUID against a uuid regular expression.
func IsUUID(placeholder string) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: func(v any) bool {
			if str, ok := v.(string); ok {
				return uuidRegex.MatchString(str)
			}
			return false
		},
	}
}

// EqualMatcher checks if the value matches the expected value.
func EqualMatcher[T any](placeholder string, expected T) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: func(v any) bool {
			if typed, ok := v.(T); ok {
				return reflect.DeepEqual(typed, expected)
			}
			return false
		},
	}
}

// Custom creates a matcher with a custom matching function.
func Custom(placeholder string, matcherFunc MatcherFunc) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: matcherFunc,
	}
}
