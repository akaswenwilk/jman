package jman

import (
	"encoding/json"
	"fmt"
)

type Obj map[string]any

func NewObjFromAny(raw any) (Obj, bool) {
	if m, ok := raw.(map[string]any); ok {
		return Obj(m), true
	}

	return nil, false
}

func (o Obj) Keys() []string {
	keys := make([]string, 0, len(o))
	for k := range o {
		keys = append(keys, k)
	}
	return keys
}

type Arr []any

func NewArrFromAny(raw any) (Arr, bool) {
	if m, ok := raw.([]any); ok {
		return Arr(m), true
	}

	return nil, false
}

func New(givenJSON any) (any, error) {
	var o any
	switch v := givenJSON.(type) {
	case string:
		data := givenJSON.(string)
		if err := json.Unmarshal([]byte(data), &o); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON string %s: %w", data, err)
		}
	case []byte:
		data := givenJSON.([]byte)
		if err := json.Unmarshal(data, &o); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON byte string %s: %w", string(data), err)
		}
	case Obj, Arr:
		// to normalize the values, marshal and unmarshal the object
		data, err := json.Marshal(givenJSON)
		if err != nil {
			return nil, fmt.Errorf("error building JSON object: marshaling %v: %w", givenJSON, err)
		}
		if err = json.Unmarshal(data, &o); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON object %v: %w", givenJSON, err)
		}
	default:
		return nil, fmt.Errorf("type %s not supported. use either string, []byte, jman.Obj, or jman.Arr", v)
	}

	return o, nil
}
