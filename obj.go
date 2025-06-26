package jman

import (
	"encoding/json"
	"fmt"
)

type Obj map[string]any

func NewObject(givenJSON any) (Obj, error) {
	var o Obj
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
	case Obj:
		// to normalize the values, marshal and unmarshal the object
		data, err := json.Marshal(givenJSON)
		if err != nil {
			return nil, fmt.Errorf("error building JSON object: marshaling %v: %w", givenJSON, err)
		}
		if err = json.Unmarshal(data, &o); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSON object %v: %w", givenJSON, err)
		}
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
