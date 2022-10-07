package jitjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// JSON value requires json marshal and unmarshal methods or json tags.
type JSON interface {
	json.Marshaler
	json.Unmarshaler
}

// JitJSON provides 'just-in-time' compilation for json marshal and unmarshal methods
// for some JSON value with type T.
type JitJSON[T JSON] struct {
	data []byte
	val  *T
}

// NewJitJSON creates new JitJSON from json byte representation.
func NewJitJSON[T JSON](data []byte) (*JitJSON[T], error) {
	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid json")
	}

	// TODO: check if cases need to be considered.
	var val T
	kind := reflect.ValueOf(val).Kind()
	switch kind {
	case reflect.Ptr:
		return nil, fmt.Errorf("cannot parse json to pointer")
	case reflect.Interface:
		return nil, fmt.Errorf("cannot parse json to interface")
	}

	jit := JitJSON[T]{
		data: data,
	}

	return &jit, nil
}

// Set new value to jitJSON.
// JitJSON byte representation is wiped to avoid data mismatch.
// Set stores a new value within JitJSON. The byte representation is cleared to avoid
// data mismatch.
func (jit *JitJSON[JSON]) Set(val JSON) {
	jit.data = nil
	jit.val = &val
}

// Marshal provides the byte representation with 'just-in-time' compilation.
// If the value is never marshalled, the initial representation will be returned.
// Since the underlying JSON type could be a pointer to some type, the representation
// is ignored when the value is set as it could become out of sync with the value.
func (jit *JitJSON[JSON]) Marshal() ([]byte, error) {
	if jit.val == nil {
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
func (jit *JitJSON[JSON]) Unmarshal() (JSON, error) {
	if jit.val != nil {
		return *jit.val, nil
	}

	var val JSON
	err := json.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}
