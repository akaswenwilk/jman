package jman

type OptsFunc func(o *EqualOptions)

type EqualOptions struct {
	matchers Matchers
}

func WithMatchers(matchers ...Matcher) OptsFunc {
	return func(o *EqualOptions) {
		for _, m := range matchers {
			o.matchers = append(o.matchers, m)
		}
	}
}
