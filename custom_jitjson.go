// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

import (
	"fmt"
)

// JitJSON provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to or from JSON is deferred until the Marshal and Unmarshal methods are called.
// You can think of JitJSON[T] as a lazy two way JSON parser with results caching implemented.
type CustomJitJSON[T any] struct {
	data   []byte
	val    *T
	parser JSONParser
}

// New creates JitJSON[T] from a value, with the default parser set.
func NewCustom[T any](val T, parser JSONParser) *CustomJitJSON[T] {
	return &CustomJitJSON[T]{val: &val, parser: parser}
}

// NewFromBytes creates a JitJSON[T] from an encoding, with the default parser set.
// If the encoding is invalid JSON, an error will be observed once Marshal is called.
func NewCustomFromBytes[T any](data []byte, parser JSONParser) *CustomJitJSON[T] {
	return &CustomJitJSON[T]{data: data, parser: parser}
}

// Set sets a new value to JitJSON[T].
func (jit *CustomJitJSON[T]) Set(val T) {
	jit.val = &val
	jit.data = nil
}

// SetBytes sets a new encoding to JitJSON[T].
func (jit *CustomJitJSON[T]) SetBytes(data []byte) {
	jit.val = nil
	jit.data = data
}

// SetParser sets a new parser to JitJSON[T].
func (jit *CustomJitJSON[T]) SetParser(parser JSONParser) {
	jit.parser = parser
}

// Marshal performs deferred json marshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Marshal' if the value has been marshaled previously. Once marshaled, the encoded value is stored with the
// jitjson for future retrieval. If there is no value to marshal, the method returns nil, nil.
func (jit *CustomJitJSON[T]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}
	if jit.val == nil {
		return nil, nil
	}
	if jit.parser == nil {
		return nil, fmt.Errorf("parser is nil")
	}

	var err error
	jit.data, err = jit.parser.Marshal(jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

// Unmarshal performs deferred json unmarshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Unmarshal' if the value has been unmarshaled previously. Once unmarshaled, the decoded value is stored with
// the jitjson for future retrieval. If there is no JSON data to unmarshal, the zero value of type T is returned.
// If the JSON data does not unmarshal into the type T, the method will return an error.
func (jit *CustomJitJSON[T]) Unmarshal() (T, error) {
	if jit.val != nil {
		return *jit.val, nil
	}
	var val T
	if jit.data == nil {
		return val, nil
	}
	if jit.parser == nil {
		return val, fmt.Errorf("parser is nil")
	}

	jit.val = &val
	err := jit.parser.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}

// MarshalJSON can be used to marshal JitJSON[T] to JSON.
func (jit *JitJSON[T]) MarshalJSON() ([]byte, error) {
	return jit.Marshal()
}

// UnmarshalJSON stores JSON data to be unmarshaled later.
func (jit *JitJSON[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data
	return nil
}
