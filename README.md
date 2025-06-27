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

[[_TOC_]]

---

## Context

- integrates with testing libraries
- easily work with json

## Quick Start

import "github.com/akaswenwilk/jman"

func TestExample(t *testing.T) {
    expected := `{"foo": "bar", "baz": 123}`
    actual := `{"baz":123, "foo":"bar"}`

    jman.Equal(t, expected, actual)
}

## API

### Basic Equal
### Options
#### Matchers
#### IgnoreArrayOrder
### Helper Methods
### json models (Arr, Obj)

## Error Messages
