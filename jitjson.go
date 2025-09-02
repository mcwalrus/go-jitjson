// Package jitjson provides a Just-In-Time JSON parser for Go.

package jitjson

// JitJSON[T] provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to or from JSON is deferred until the Marshal and Unmarshal methods are called.
// You can think of JitJSON[T] as a lazy two way JSON parser, implemented for deferred value retrieval.
// Caching is always enabled and will store the parsed values for future retrieval.
type JitJSON[T any] struct {
	data   []byte
	val    *T
	parser jsonParser
}

// New creates JitJSON[T] from a value.
// The default parser is used, which can be changed with the SetParser method.
func New[T any](val T) *JitJSON[T] {
	return &JitJSON[T]{val: &val, parser: defaultParser()}
}

// NewFromBytes creates a JitJSON[T] from JSON byte data.
// The default parser is used, which can be changed with the SetParser method.
func NewFromBytes[T any](data []byte) *JitJSON[T] {
	return &JitJSON[T]{data: data, parser: defaultParser()}
}

// Set JitJSON[T] to a new value.
func (jit *JitJSON[T]) Set(val T) {
	jit.val = &val
	jit.data = nil
}

// SetParser sets the parser to use for the JitJSON[T].
// If the version is not supported, the function returns an error.
func (jit *JitJSON[T]) SetParser(version ParserVersion) error {
	p, err := parserFromVersion(version)
	if err != nil {
		return err
	}
	jit.parser = p
	return nil
}

// MustSetParser sets the parser and panics if the version is not supported.
func (jit *JitJSON[T]) MustSetParser(version ParserVersion) {
	if err := jit.SetParser(version); err != nil {
		panic(err)
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

	var err error
	jit.data, err = jit.parser.marshal(jit.val)
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

	jit.val = &val
	err := jit.parser.unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}

// Parser returns the name of the json library version.
func (jit *JitJSON[T]) Parser() string {
	return jit.parser.name()
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
