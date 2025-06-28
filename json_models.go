package jman

import (
	"encoding/json"
	"fmt"
)

const (
	ObjectType  = "object"
	ArrayType   = "array"
	UnknownType = "unknown"
)

type JSONModel struct {
	data any
}

func (j *JSONModel) Obj() (Obj, error) {
	dataObj, ok := j.data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("data is not a JSON object")
	}
	return Obj(dataObj), nil
}

func (j *JSONModel) IsObj() bool {
	_, ok := j.data.(map[string]any)
	return ok
}

func (j *JSONModel) Arr() (Arr, error) {
	dataObj, ok := j.data.([]any)
	if !ok {
		return nil, fmt.Errorf("data is not a JSON array")
	}
	return Arr(dataObj), nil
}

func (j *JSONModel) IsArr() bool {
	_, ok := j.data.([]any)
	return ok
}

func (j *JSONModel) UnderlyingType() string {
	switch j.data.(type) {
	case map[string]any:
		return ObjectType
	case []any:
		return ArrayType
	default:
		return UnknownType
	}
}

type Obj map[string]any

func NewObjFromAny(raw any) (Obj, bool) {
	if m, ok := raw.(map[string]any); ok {
		return Obj(m), true
	}

	return nil, false
}

func (o Obj) MustString() string {
	return string(o.MustMarshal())
}

func (o Obj) MustMarshal() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(fmt.Sprintf("error marshalling object %v: %v", o, err))
	}
	return data
}

type Arr []any

func NewArrFromAny(raw any) (Arr, bool) {
	if m, ok := raw.([]any); ok {
		return Arr(m), true
	}

	return nil, false
}

func (a Arr) MustString() string {
	return string(a.MustMarshal())
}

func (a Arr) MustMarshal() []byte {
	data, err := json.Marshal(a)
	if err != nil {
		panic(fmt.Sprintf("error marshalling array %v: %v", a, err))
	}
	return data
}

func New(givenJSON any) (JSONModel, error) {
	var result any
	switch v := givenJSON.(type) {
	case string:
		data := givenJSON.(string)
		if err := json.Unmarshal([]byte(data), &result); err != nil {
			return JSONModel{}, fmt.Errorf("error unmarshalling JSON string %s: %w", data, err)
		}
	case []byte:
		data := givenJSON.([]byte)
		if err := json.Unmarshal(data, &result); err != nil {
			return JSONModel{}, fmt.Errorf("error unmarshalling JSON byte string %s: %w", string(data), err)
		}
	case Obj, Arr:
		// to normalize the values, marshal and unmarshal the object
		data, err := json.Marshal(givenJSON)
		if err != nil {
			return JSONModel{}, fmt.Errorf("error building JSON object: marshaling %v: %w", givenJSON, err)
		}
		if err = json.Unmarshal(data, &result); err != nil {
			return JSONModel{}, fmt.Errorf("error unmarshalling JSON object %v: %w", givenJSON, err)
		}
	default:
		return JSONModel{}, fmt.Errorf("type %s not supported. use either string, []byte, jman.Obj, or jman.Arr", v)
	}

	return JSONModel{result}, nil
}
