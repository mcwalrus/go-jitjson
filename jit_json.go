package jitjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JitJSON provides 'just-in-time' compilation to parse json encoding to value type
// and vice versa. JitJSON can be applied for any generic types expect interfaces and
// type pointers. See NewJitJSON for more details.
type JitJSON[T any] struct {
	data []byte
	val  *T
}

// NewJitJSON creates new JitJSON from json encoding and a specified generic type.
// If the json encoding is not valid, or the generic type is either a interface or
// type pointer, an error will be returned.
//
// There will be a slight overhead using a JitJSON to perform the additional json
// validation and reflection. Keep in mind that json validation does not ensure the
// encoded json will marshal into the given type, only that it is of valid encoding.
func NewJitJSON[T any](data []byte) (*JitJSON[T], error) {
	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid json")
	}

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

	jit := JitJSON[T]{
		data: data,
	}

	return &jit, nil
}

// Set new value to JitJSON. Value will not be parsed to json encoding until 'Marshal'
// is called.
func (jit *JitJSON[T]) Set(val T) {
	jit.data = nil
	jit.val = &val
}

// Marshal provides the byte representation with 'just-in-time' compilation. If the
// value has not yet been marshalled, the initial byte representation will be returned.

// Marshal performs 'just-in-time' compilation of parsing type to json encoding. If the
// encoding of the , the initial byte representation will be returned.
func (jit *JitJSON[T]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}

	var err error
	jit.data, err = json.Marshal(jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

// Unmarshal provides the value with 'just-in-time' compilation. After the first
// unmarshal, the value can be returned without further repeated unmarshalling.
func (jit *JitJSON[T]) Unmarshal() (T, error) {
	if jit.val != nil {
		return *jit.val, nil
	}

	var val T
	jit.val = &val
	err := json.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}
