// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

import (
	"encoding/json"
	"fmt"
)

// JitJSON[T any] provides 'just-in-time' (JIT) JSON parsing capabilities in Go.
// It can hold a JSON encoding or a value of any type (T). The type T can be parsed to and from JSON/Go types only when needed.
// This is achieved through the 'Marshal' and 'Unmarshal' methods of JitJSON.
type JitJSON[T any] struct {
	data []byte
	val  *T
}

// AnyJitJSON is implemented by the JitJSON[T any] type, where T can be any type.
// This means you can use AnyJitJSON with any underlying type that can be marshaled or unmarshaled to / from JSON.
// See test file for examples.
type AnyJitJSON interface {
	private()
	json.Marshaler
	json.Unmarshaler
}

// NewJitJSON[T any] constructs a new JitJSON[T] type.
// The constructor accepts only values of JSON encoding or value of type T, otherwise an error is returned.
// JSON encoding can be either of type []byte or json.RawMessage. Nil is also a valid value for JitJSON[T] of no associated data.
// Empty values of JitJSON[T] is also valid constructors. If the value provided of T is nil, the method 'Unmarshal' will return the zero value of type T.
func NewJitJSON[T any](val interface{}) (*JitJSON[T], error) {
	var jit JitJSON[T]
	if val == nil {
		return &jit, nil
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
		return nil, fmt.Errorf("unexpected type: %T", val)
	}

	return &jit, nil
}

// private implements AnyJitJSON.
func (jit *JitJSON[T]) private() {}

// MarshalJSON implements the json.Marshaler interface.
// It returns the JSON encoding of the JitJSON[T] value by calling the Marshal method.
// If an error occurs during the marshaling process, it returns the error.
func (jit *JitJSON[T]) MarshalJSON() ([]byte, error) {
	return jit.Marshal()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It sets the internal data of the JitJSON[T] to the provided byte slice and resets the value to nil.
// This method does not perform the actual unmarshaling; the unmarshaling is deferred until the Unmarshal method is called.
// This method always returns nil as error since it only assigns the input byte slice to the internal data.
func (jit *JitJSON[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data
	return nil
}

// Marshal performs the equivalent of 'json.Marshal' for the value of JitJSON[T].
// If the value has already been resolved, the method returns the existing encoding without re-evaluating 'json.Marshal'.
// If the value of JitJSON[T] is nil, the method returns nil.
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

// Unmarshal performs the equivalent of 'json.Unmarshal' for the encoded value of JitJSON[T].
// If the value has already been resolved, the method returns the existing value without re-evaluating 'json.Unmarshal'.
// If the encoded value of JitJSON[T] is nil, the method returns the zero value of type T.
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

// Read implements the io.Reader interface.
// It marshals the value of JitJSON[T] and copies the result to the provided byte slice.
// If the value of JitJSON[T] is nil, the method returns 0, nil.
func (jit JitJSON[T]) Read(p []byte) (n int, err error) {
	data, err := jit.Marshal()
	if err != nil {
		return 0, nil
	}
	return copy(p, data), nil
}
