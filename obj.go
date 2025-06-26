package jman

import (
	"encoding/json"
	"fmt"
)

type Obj map[string]any

func New(givenJSON any) (Obj, error) {
	var o Obj
	switch v := givenJSON.(type) {
	case string:
		data := givenJSON.(string)
		err := json.Unmarshal([]byte(data), &o)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON string %s: %w", data, err)
		}
	case []byte:
		data := givenJSON.([]byte)
		err := json.Unmarshal(data, &o)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON byte string %s: %w", string(data), err)
		}
	case Obj:
		return givenJSON.(Obj), nil
	default:
		return nil, fmt.Errorf("type %s not supported. use either string, []byte, or jman.Obj", v)
	}

	return o, nil
}

func (o Obj) Keys() []string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	return keys
}
