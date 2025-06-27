package jman

import "reflect"

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

func NotEmpty(placeholder string) Matcher {
	return Matcher{
		Placeholder: placeholder,
		MatcherFunc: isEmpty,
	}
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return val.Len() == 0
	case reflect.Ptr:
		if val.IsNil() {
			return true
		}

		return isEmpty(val.Elem().Interface())
	default:
		return reflect.DeepEqual(v, reflect.Zero(val.Type()).Interface())
	}
}
