package jman

import (
	"errors"
	"regexp"
)

const base = "$"

type OptsFunc func(o *EqualOptions)

type EqualOptions struct {
	matchers         Matchers
	ignoreArrayOrder []string
}

var pathStartRegex = regexp.MustCompile(`^\$\.`)

func (o EqualOptions) valid() error {
	if len(o.ignoreArrayOrder) == 0 {
		return nil
	}

	for _, key := range o.ignoreArrayOrder {
		if key == "" || !pathStartRegex.MatchString(key) {
			return errors.New("key for ignoring array order must start with " + base)
		}
	}
	return nil
}

func WithMatchers(matchers ...Matcher) OptsFunc {
	return func(o *EqualOptions) {
		for _, m := range matchers {
			o.matchers = append(o.matchers, m)
		}
	}
}

func WithDefaultMatchers(matchers Matchers) OptsFunc {
	return func(o *EqualOptions) {
		o.matchers = matchers
	}
}

func WithIgnoreArrayOrder(keys ...string) OptsFunc {
	return func(o *EqualOptions) {
		o.ignoreArrayOrder = append(o.ignoreArrayOrder, keys...)
	}
}
