package jman

const base = "$"

type optsFunc func(o *equalOptions)

type equalOptions struct {
	matchers         Matchers
	ignoreArrayOrder []string
}

func (o equalOptions) valid() error {
	if len(o.ignoreArrayOrder) == 0 {
		return nil
	}

	for _, key := range o.ignoreArrayOrder {
		_, err := getPathParts(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithMatchers allows you to add matchers to the comparison options.
func WithMatchers(matchers ...Matcher) optsFunc {
	return func(o *equalOptions) {
		for _, m := range matchers {
			o.matchers = append(o.matchers, m)
		}
	}
}

// WithDefaultMatchers allows you to set the default matchers for the comparison options.
// Useful if defining a set of matchers that should be used in most comparisons.
func WithDefaultMatchers(matchers Matchers) optsFunc {
	return func(o *equalOptions) {
		o.matchers = matchers
	}
}

// WithIgnoreArrayOrder allows you to specify keys for which the order of array elements should be ignored during comparison.
// each key should be a valid JSON path, and the order of elements in arrays at those paths will not be considered during comparison.
// the path must start with $.  For ignoring order of the base array, use "$" as the key.
func WithIgnoreArrayOrder(keys ...string) optsFunc {
	return func(o *equalOptions) {
		o.ignoreArrayOrder = append(o.ignoreArrayOrder, keys...)
	}
}
