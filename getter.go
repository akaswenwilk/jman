package jman

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var numberRegex = regexp.MustCompile(`^\d+$`)

func getValue[R JSONEqual](paths []string, data R) (any, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no path segments provided")
	}
	current := any(data)
	for _, path := range paths {
		if isIndex(path) {
			index, _ := strconv.Atoi(path) // converting after checking regex, no need for error handling here
			if arr, ok := current.(Arr); ok {
				if index < 0 || index >= len(arr) {
					return nil, fmt.Errorf("index %d out of bounds for array of length %d", index, len(arr))
				}
				current = arr[index]
			} else {
				return nil, fmt.Errorf("expected an array at path '%s', got %T", strings.Join(paths, "."), current)
			}
		} else {
			if converted, ok := current.(Obj); ok {
				if val, ok := converted[path]; ok {
					current = val
				} else {
					return nil, fmt.Errorf("key '%s' not found in object", path)
				}
			}

		}
	}

	return current, nil
}

func isIndex(segment string) bool {
	return numberRegex.MatchString(segment)
}
