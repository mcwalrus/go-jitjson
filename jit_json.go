package jitjson

import (
	"encoding/json"
	"fmt"
)

// JitJSON provides 'just-in-time' compilation to encode or decode json data or value of T.
// Use type to parse values to and from json / Go types only when needed.
// See 'Marshal' and 'Unmarshal' methods for usage.
type JitJSON[T any] struct {
	data []byte
	val  *T
}

// NewJitJSON creates new jit-json based on either json encoding or a value of T.
//
// The json encoding can be either of type []byte or json.RawMessage. nil is also a valid
// value for an empty jit-json of no associated data. An empty jit-json can be useful for
// later use of unmarshalling. If the argument is neither of a json encoding or of type T,
// an error will be returned.
func NewJitJSON[T any](val interface{}) (JitJSON[T], error) {
	var jit JitJSON[T]
	if val == nil {
		return jit, nil
	}

	switch v := val.(type) {
	case []byte:
		jit = JitJSON[T]{
			data: v,
		}
	case json.RawMessage:
		jit = JitJSON[T]{
			data: []byte(v),
		}
	case T:
		jit = JitJSON[T]{
			val: &v,
		}
	default:
		return jit, fmt.Errorf("unexpected type: %T", val)
	}

	return jit, nil
}

// MarshalJSON implements json.Marshaler.
// For more documentation, see: https://pkg.go.dev/encoding/json#Marshaler.
func (jit *JitJSON[T]) MarshalJSON() ([]byte, error) {
	return jit.Marshal()
}

// Marshal performs 'json.Marshal' for value of JitJSON. If the value is already resolved
// for, the encoding will be returned without re-evaluation of 'json.Marshal'.
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

// UnmarshalJSON implements json.Unmarshaler.
// For more documentation, see: https://pkg.go.dev/encoding/json#Unmarshaler.
func (jit *JitJSON[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data

	return nil
}

// Unmarshal performs 'json.Unmarshal' for encoding of JitJSON. If the value is already
// resolved for, the value is returned without re-evaluation of 'json.Unmarshal'.
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
