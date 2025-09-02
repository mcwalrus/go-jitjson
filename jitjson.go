// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

import "fmt"

// JitJSON[T] provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to or from JSON is deferred until the Marshal and Unmarshal methods are called.
// You can think of JitJSON[T] as a lazy two way JSON parser, implemented for deferred value retrieval.
// Caching is always enabled and will store the parsed values for future retrieval.
type JitJSON[T any] struct {
	data   []byte
	val    *T
	parser JSONParser
}

// New creates JitJSON[T] from a value, with the default parser set.
func New[T any](val T) *JitJSON[T] {
	return &JitJSON[T]{val: &val, parser: getDefaultParser()}
}

// NewFromBytes creates a JitJSON[T] from an encoding, with the default parser set.
// If the encoding is invalid JSON, an error will be observed once Marshal is called.
func NewFromBytes[T any](data []byte) *JitJSON[T] {
	return &JitJSON[T]{data: data, parser: getDefaultParser()}
}

// Set a new value to JitJSON[T].
func (jit *JitJSON[T]) Set(val T) {
	jit.val = &val
	jit.data = nil
}

// SetParser sets the parser to use for the JitJSON[T].
// Returns an error if the parser is not pre-registered by using [RegisterParser].
func (jit *JitJSON[T]) SetParser(name string) error {
	parser, exists := parsers[name]
	if !exists {
		return fmt.Errorf("parser %s not registered", name)
	}
	jit.parser = parser
	return nil
}

// Parser returns the name of the parser used by JitJSON[T].
// A parser might be nil when the jitjson was initialised without a parser.
// In this case, the JitJSON will return "<nil>" which is later set to the default parser.
func (jit *JitJSON[T]) Parser() string {
	if jit.parser == nil {
		return "<nil>"
	} else {
		return jit.parser.Name()
	}
}

// Marshal performs deferred json marshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Marshal' if the value has been marshaled previously. Once marshaled, the encoded value is stored with the
// jitjson for future retrieval. If there is no value to marshal, the method returns nil, nil.
func (jit *JitJSON[T]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}
	if jit.val == nil {
		return nil, nil
	}
	if jit.parser == nil {
		jit.parser = getDefaultParser()
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
func (jit *JitJSON[T]) Unmarshal() (T, error) {
	if jit.val != nil {
		return *jit.val, nil
	}
	var val T
	if jit.data == nil {
		return val, nil
	}
	if jit.parser == nil {
		jit.parser = getDefaultParser()
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
