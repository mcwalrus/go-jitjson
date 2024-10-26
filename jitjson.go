// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

import (
	"encoding/json"
)

// JitJSON[T] provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to/from JSON is deferred until needed via Marshal and Unmarshal methods.
// You can think of JitJSON[T] as a lazy JSON parser.
type JitJSON[T any] struct {
	data json.RawMessage
	val  *T
}

// NewJitJSON constructs a new JitJSON[T] from a value of type T.
func NewJitJSON[T any](val T) *JitJSON[T] {
	return &JitJSON[T]{val: &val}
}

// BytesToJitJSON constructs a new JitJSON[T] from a JSON encoding.
func BytesToJitJSON[T any](data []byte) *JitJSON[T] {
	return &JitJSON[T]{data: data}
}

// MarshalJSON implements json.Marshaler, simply calling the Marshal method.
func (jit *JitJSON[T]) MarshalJSON() ([]byte, error) {
	return jit.Marshal()
}

// Marshal performs deferred marshaling of the value of JitJSON[T]. If the encoding has already been resolved,
// the method returns the existing encoding without re-evaluating 'json.Marshal'. If the JitJSON[T] is empty,
// the method returns nil, nil. The encoded value is cached for future use.
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

// UnmarshalJSON implements json.Unmarshaler by storing the JSON encoding for deferred unmarshaling.
// If jit has been used for previous unmarshaling, the method resets the jit for the new encoding.
func (jit *JitJSON[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data
	return nil
}

// Unmarshal performs deferred unmarshaling of the value of JitJSON[T]. If the value has already been resolved,
// the method returns the existing value without re-evaluating 'json.Unmarshal'. If the JitJSON[T] is empty, the
// zero value of type T is returned. The decoded value is cached for future use.
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

// Read allows JitJSON[T] to be used with json.Decoder or other readers which use io.Reader.
func (jit JitJSON[T]) Read(p []byte) (n int, err error) {
	data, err := jit.Marshal()
	if err != nil {
		return 0, nil
	}
	return copy(p, data), nil
}
