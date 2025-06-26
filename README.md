
     ██╗ ███╗   ███╗ █████╗ ███╗   ██╗
     ██║████╗ ████║██╔══██╗████╗  ██║
     ██║██╔████╔██║███████║██╔██╗ ██║
██   ██║██║╚██╔╝██║██╔══██║██║╚██╗██║
╚█████╔╝██║ ╚═╝ ██║██║  ██║██║ ╚████║
╚════╝  ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝

# jman (json manipulator)

A lightweight Go library for working with and asserting JSON in tests.  
Minimal, fast, idiomatic.

---

## Table of Contents

[[_TOC_]]

---

## Features

- easily compare json objects using strings, bytes, or jman.Obj
- support for dynamic comparison for values
- support for ignoring order when comparing arrays
- integrates with any test suite supporting testing.T
- easily access different parts of a json object with dot notation: e.g. `{"hello": "world", "foo":[{"bar":"baz}]}` `"hello.foo.0.bar"`

---

## Quick Start

import "github.com/akaswenwilk/jman"

func TestExample(t *testing.T) {
    expected := `{"foo": "bar", "baz": 123}`
    actual := `{"baz":123, "foo":"bar"}`

    jman.Equal(t, expected, actual)
}

## API
