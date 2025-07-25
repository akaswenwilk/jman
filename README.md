<pre>
                     ██╗███╗   ███╗ █████╗ ███╗   ██╗
                     ██║████╗ ████║██╔══██╗████╗  ██║
                     ██║██╔████╔██║███████║██╔██╗ ██║
                ██   ██║██║╚██╔╝██║██╔══██║██║╚██╗██║
                ╚█████╔╝██║ ╚═╝ ██║██║  ██║██║ ╚████║
                ╚════╝  ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝
</pre>

# jman (json manipulator)

A lightweight Go library for working with and asserting JSON in tests.  
Minimal, fast, idiomatic.

---

# Table of Contents

- [Context](#context)
- [Quick Start](#quick-start)
- [API](#api)
  - [Basic Models](#basic-models)
  - [Basic Equal](#basic-equal)
    - [Error Messages](#error-messages)
  - [Options](#options)
  - [Helper Methods](#helper-methods)

---

## Context

This package is intended for testing.  There are many bits of syntactical sugar that make working with it more flexible and simpler to write.  It gets around much of the type safety and can panic with certain methods.  Therefore it is really not intended for production where these safety measures help ensure correct code!  

However in testing, these safety measures can be hindrances.  The Design Philosophy of Jman is:

- concise
- flexible
- clear

## Quick Start

```go
import "github.com/akaswenwilk/jman"

func TestExample(t *testing.T) {
    expected := jman.Obj{
        "foo": "bar",
        "baz": "$ANY",
        "id": "$UUID",
    }
    actual := `{"bas":123, "foo":"bar", "id": "d47422ff-683a-4077-b3eb-e06a99bc9b55"}`

    err := expected.Equal(actual, jman.WithMatcher(
            jman.NotEmpty("ANY$"),
            jman.IsUUID("$UUID"),
        ),
    )

    if err != nil {
        t.FatalF("expected no comparison error, got %w", err)
    }
}
```

## API


### Basic Models

jman is built around two structs: `jman.Obj` and `jman.Arr`, corresponding to `map[string]any` and `[]any`.  They then have helper methods to manipulate and compare with each other.  Whenever using any of the methods to generate, edit, fetch, or compare values, they will be "normalized".  Since go by default will unmarshal into specific value (I.E. any number will be unmarshalled into a float), it can get pretty annoying to ensure that you have specified the correct value in an array or object.  Therefore the normalization will ensure correctness of value.  So instead of having to do this:

```go
expected := jman.Obj{
    "numberKey": float64(123),
}
actual := `{"numberKey":123}`

err := expected.Equal(actual)
```

you can simply write:
```go
expected := jman.Obj{
    "numberKey": 123,
}
actual := `{"numberKey":123}`

err := expected.Equal(actual)
```

and the conversion will happen via normalization.

### Basic Equal

the main method for comparison is the `Equal` method present on both `jman.Arr` and `jman.Obj`

`Equal` can compare against other objects/arrays, but also against strings or byte slices.  No need to add any extra steps to ensure that the expected and actual are both in the same format just to compare the values!


#### Error Messages.

Error messages are given in dot notation, always preceded by the base character of `$`

Each difference between two json models is listed with its full path and a hopefully clear and concise message of the differences.  This can be seen with the following test:


```go 
func TestObj_Equal_Unequal_ManyObjectErrors(t *testing.T) {
	expected := jman.Obj{
		"key1": "value1",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue1",
			"nestedKey2": 42,
			"nestedKey4": "unexpectedKey", // unexpected key
		},
	}
	actual := jman.Obj{
		"key3": "value3",
		"key2": jman.Obj{
			"nestedKey1": "nestedValue2", // different value1
			"nestedKey2": "notANumber",  // different TestNewObj_UnsupportedType
			"nestedKey3": true,           // unexpected key3
		},
	}

	err := expected.Equal(actual)
	require.Error(t, err)
	assert.Equal(t, `expected not equal to actual:
$.key1 not found in actual
$.key3 unexpected key
$.key2.nestedKey4 not found in actual
$.key2.nestedKey3 unexpected key
$.key2.nestedKey1 expected "nestedValue1" - actual "nestedValue2"
$.key2.nestedKey2 expected 42 - actual "notANumber"
`, err.Error())
}
```

### Options
In addition to giving flexibility on types for comparison, jman also aims to give flexibility on actual comparisons. For example, if a value can be dynamic, and it's important only that it exists and not specifically what value  is present:

```go 
expected := jman.Obj{
    "id": "$SOME_ID"
}
actual := `{"id":"21493933-3338-450c-8113-62a35d2c1820"}`
err := expected.Equal(actual, jman.WithMatchers(
    jman.NotEmpty("$SOME_ID")
))
```

jman provides several matchers that can be used out of the box:
`NotEmpty` - checks that whatever is present is not a zero value or nil.
`IsUUID` - checks that something matches a uuid regex
`EqaulMatcher` - checks if a value deep equals what is received. useful for tests with dynamic values that need to be matched exactly
`Custom` - for passing in a custom matcher

the placeholder tells what value in the expected will be checked with the corresponding function from the matcher with the value from the actual.

In addition to the matchers, there is also an option to ignore array ordering:
```go 
	expected := jman.Obj{
		"items": jman.Arr{"item1", "item2", "item3"},
	}
	actual := jman.Obj{
		"items": jman.Arr{"item3", "item1", "item2"},
	}

	assert.NoError(t, expected.Equal(actual,
		jman.WithIgnoreArrayOrder("$.items"),
	))

```
there can be multiple paths passed in if multiple arrays should ignore order. Each path must follow the syntax of leading with `$`.

### Helper Methods

Additionally, both `jman.Arr` and `jman.Obj` have a set of helper methods to make json manipulation easier:

`Set(path string, value any) error`
sets the value in whatever path is given. The path is separated by dots and must begin with `$`.  For example:
```go 
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	err := data.Set("$.key1.nestedKey1.1.deepKey", "newDeepValue")
```
if a key is not present in the path, it will create a new key, otherwise it will overwrite the existing value of that key.

in arrays, it will replace the item at index, but it will not add new items.

`Get(path string) Result`
This will return a `Result` of the path.  Result will then either have the error or the data:
```go 
	data := jman.Obj{
		"key1": jman.Obj{
			"nestedKey1": jman.Arr{"item1", jman.Obj{"deepKey": "deepValue"}},
		},
	}
	value, err := data.Get("$.key1.nestedKey1.1.deepKey").String()
```

result can return its data using one of the following methods:
`Data` (`any`)
`String`
`Number` (`float64`)
`Bool`
`Obj` (`jman.Obj`)
`Arr` (`jman.Arr`)


`MustString() string`
`MustBytes() []byte`

both `jman.Obj` and `jman.Arr` can marshal json into outputs of either a string or byte slice.  These are convenience helpers that will either successfully return the marshalled value or panic.  Again, not for production code.
