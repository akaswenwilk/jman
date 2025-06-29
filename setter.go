package jman

import (
	"fmt"
	"strconv"
	"strings"
)

func setByPath(data any, path string, value any) error {
	if !strings.HasPrefix(path, "$.") {
		return fmt.Errorf("invalid path: must start with $.")
	}
	segments := strings.Split(path[2:], ".")
	current := data

	for i, segment := range segments {
		isLast := i == len(segments)-1

		switch curr := current.(type) {
		case Obj:
			if isLast {
				curr[segment] = value
				return nil
			}
			next, ok := curr[segment]
			if !ok {
				if isIndex(segments[i+1]) {
					curr[segment] = Arr{}
				} else {
					curr[segment] = Obj{}
				}
				next = curr[segment]
			}
			current = next
		case Arr:
			idx, err := strconv.Atoi(segment)
			if err != nil || idx < 0 {
				return fmt.Errorf("invalid array index: %s", segment)
			}
			for len(curr) <= idx {
				curr = append(curr, nil)
			}
			if isLast {
				curr[idx] = value
				return nil
			}
			if curr[idx] == nil {
				if isIndex(segments[i+1]) {
					curr[idx] = Arr{}
				} else {
					curr[idx] = Obj{}
				}
			}
			current = curr[idx]
			// must assign back in case the array grew
			if parent, ok := current.([]any); ok {
				current = Arr(parent)
			}
		default:
			return fmt.Errorf("unexpected type at segment %s", segment)
		}
	}
	return nil
}
