package jman

import (
	"errors"
	"fmt"
	"strings"
)

func getPathParts(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	parts := strings.Split(path, ".")
	if parts[0] != base {
		return nil, fmt.Errorf("path must start with %s", base)
	}

	return parts, nil
}
