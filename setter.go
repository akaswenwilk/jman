package jman

import (
	"fmt"
	"maps"
	"strconv"
)

func setByPath(data any, path string, value any) error {
	paths, err := getPathParts(path)
	if err != nil {
		return err
	}

	current := data

	for i, p := range paths { // skip the base part
		if p == base {
			continue
		}
		isLast := i == len(paths)-1

		switch curr := current.(type) {
		case Obj:
			if isLast {
				curr[p] = value
				return nil
			}
			next, ok := curr[p]
			if !ok {
				if isIndex(paths[i+1]) {
					curr[p] = Arr{}
				} else {
					curr[p] = Obj{}
				}
				next = curr[p]
			}
			current = next
		case Arr:
			idx, err := strconv.Atoi(p)
			if err != nil || idx < 0 {
				return fmt.Errorf("invalid array index: %s", p)
			}
			for len(curr) <= idx {
				curr = append(curr, nil)
			}
			if isLast {
				curr[idx] = value
				return nil
			}
			if curr[idx] == nil {
				if isIndex(paths[i+1]) {
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
			return fmt.Errorf("unexpected type at segment %s", p)
		}
	}

	switch v := data.(type) {
	case Obj:
		normalized, err := normalize(v)
		if err != nil {
			return err
		}
		maps.Copy(v, normalized)
	case Arr:
		normalized, err := normalize(v)
		if err != nil {
			return err
		}
		copy(v, normalized)
	}

	return nil
}
