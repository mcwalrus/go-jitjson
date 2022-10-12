package jitjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JitJSON provides 'just-in-time' compilation to encode / decode json.
// JitJSON can be applied for any generic types expect interfaces and pointers of types.
type JitJSON[T any] struct {
	data []byte
	val  *T
}

// NewJitJSON creates new JitJSON from json encoding and a specified generic type. If the json encoding
// is not valid, or the generic type is either a interface or type pointer, an error will be returned.
//
// When data is nil, technically this is invalid json. To handle this, when JitJSON.Marshal is called,
// the default value of the generic type will be returned.
//
// Keep in mind, json validation only checks that the encoding is valid and not whether the encoding
// will unmarshal with the generic type since this would require parsing the json.
func NewJitJSON[T any](data []byte) (*JitJSON[T], error) {
	var val T
	kind := reflect.ValueOf(val).Kind()
	switch kind {
	case reflect.Ptr:
		return nil, fmt.Errorf("cannot parse json to pointer")
	case reflect.Interface:
		return nil, fmt.Errorf("cannot parse json to interface")
	case reflect.Invalid:
		return nil, fmt.Errorf("cannot parse json to invalid type")
	}

	if data == nil {
		return &JitJSON[T]{}, nil
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid json")
	}

	jit := JitJSON[T]{
		data: data,
	}

	return &jit, nil
}

// Set new value to JitJSON. Parsing of value will occur 'just-in-time' when 'Marshal' is called.
func (jit *JitJSON[T]) Set(val T) {
	jit.data = nil
	jit.val = &val
}

// Marshal performs 'just-in-time' compilation of json marshalling. If JitJSON contains the json
// encoding, this will returned to avoiding an unnecessary parse. In cases where the encoding and
// the value is not set, nil will be returned.
func (jit *JitJSON[T]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}

	if jit.val == nil {
		return nil, nil
	}

	var err error
	jit.data, err = json.Marshal(jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

// Unmarshal performs 'just-in-time' compilation of json unmarshalling. If JitJSON contains the
// decoded value, this will returned to avoiding an unnecessary parse. In cases where the value
// and the encoding is not set, the default value of the generic type will be returned.
func (jit *JitJSON[T]) Unmarshal() (T, error) {
	if jit.val != nil {
		return *jit.val, nil
	}

	var val T
	if jit.data == nil {
		return val, nil
	}

	jit.val = &val
	err := json.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}
