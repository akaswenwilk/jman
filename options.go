package jman

const base = "$"

type OptsFunc func(o *EqualOptions)

type EqualOptions struct {
	matchers         Matchers
	ignoreArrayOrder []string
}

func (o EqualOptions) valid() error {
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
