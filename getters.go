package jman

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Get[T any, R JSONEqual](path string, data R) (T, error) {
	var result T
	if path == "" {
		return result, fmt.Errorf("path cannot be empty")
	}
	paths := strings.Split(path, ".")
	if paths[0] != "$" {
		return result, fmt.Errorf("path must start with '$', got '%s'", paths[0])
	}

	return getValue[T, R](paths[1:], data)
}

var numberRegex = regexp.MustCompile(`^\d+$`)

func getValue[T any, R JSONEqual](paths []string, data R) (T, error) {
	var zero T
	if len(paths) == 0 {
		return zero, fmt.Errorf("no path segments provided")
	}
	current := any(data)
	for _, path := range paths {
		if numberRegex.MatchString(path) {
			index, _ := strconv.Atoi(path) // converting after checking regex, no need for error handling here
			if arr, ok := current.(Arr); ok {
				if index < 0 || index >= len(arr) {
					return zero, fmt.Errorf("index %d out of bounds for array of length %d", index, len(arr))
				}
				current = arr[index]
			} else {
				return zero, fmt.Errorf("expected an array at path '%s', got %T", strings.Join(paths, "."), current)
			}
		} else {
			if converted, ok := current.(Obj); ok {
				if val, ok := converted[path]; ok {
					current = val
				} else {
					return zero, fmt.Errorf("key '%s' not found in object", path)
				}
			}

		}
	}

	if result, ok := current.(T); ok {
		return result, nil
	}

	return zero, fmt.Errorf("value at path '%s.%s' is not of type %T", base, strings.Join(paths, "."), zero)
}
