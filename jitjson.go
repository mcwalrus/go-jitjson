// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

import (
	"encoding/json"
)

// JitJSON[T] provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to/from JSON is deferred until needed via Marshal and Unmarshal methods.
// You can think of JitJSON[T] as a lazy two way JSON parser, with cache.
type JitJSON[T any] struct {
	data []byte
	val  *T
}

// NewJitJSON constructs a new JitJSON[T] from a value of type T.
func NewJitJSON[T any](val T) *JitJSON[T] {
	return &JitJSON[T]{val: &val}
}

// NewJitJSONFromBytes constructs a new JitJSON[T] from a JSON encoding for a value of type T.
func NewJitJSONFromBytes[T any](data []byte) *JitJSON[T] {
	return &JitJSON[T]{data: data}
}

// MarshalJSON returns the JSON encoding of JitJSON[T].
func (jit *JitJSON[T]) MarshalJSON() ([]byte, error) {
	return jit.Encode()
}

// Encode performs deferred json marshaling of the value of JitJSON[T]. The method can return without evaluating
// 'json.Marshal' if the value has been constructed from bytes or has already been marshaled. Once marshaled, the
// encoded value is cached for future use. If the value of T is nil, the method returns nil, nil.
func (jit *JitJSON[T]) Encode() ([]byte, error) {
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

// UnmarshalJSON stores the encoding to JitJSON[T].
func (jit *JitJSON[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data
	return nil
}

// Decode performs deferred json unmarshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Unmarshal' if the value is stored by the JitJSON or has already been unmarshaled. Once unmarshaled, the
// decoded value is cached for future use. If the JitJSON[T] is empty, the zero value of type T is returned.
func (jit *JitJSON[T]) Decode() (T, error) {
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

// Set sets the value of JitJSON[T] to the provided value of T.
func (jit *JitJSON[T]) Set(val T) {
	jit.val = &val
	jit.data = nil
}
