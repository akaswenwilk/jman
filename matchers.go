package jman

import (
	"reflect"
	"regexp"
)

type Matchers []Matcher

func (m Matchers) FindByPlaceholder(p string) (Matcher, bool) {
	for _, matcher := range m {
		if matcher.Placeholder == p {
			return matcher, true
		}
	}

	return Matcher{}, false
}

type Matcher struct {
	Placeholder string
	MatcherFunc
}

type MatcherFunc func(v any) bool

// Equal checks if the value matches the expected value.
// a string must not be empty, cannot be nil.  booleans can be either true or false.
// a number can be any number, including zero.
// an array must not be empty, cannot be nil.  an object must not be empty, cannot be nil.
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

func Custom(placeholder string, matcherFunc MatcherFunc) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: matcherFunc,
	}
}
