//go:build go1.25 && goexperiment.jsonv2

package jitjson

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

var _ jsonv2.MarshalerTo = (*JitJSONV2[any])(nil)
var _ jsonv2.UnmarshalerFrom = (*JitJSONV2[any])(nil)

// JitJSONV2 provides just-in-time (JIT) JSON parsing in Go for a value of type T.
// Parsing to or from JSON is deferred until the Marshal and Unmarshal methods are called.
// Type implements parsing with the encoding/json/v2 library and supports new json/v2 interfaces.
type JitJSONV2[T any] struct {
	data []byte
	val  *T
}

// NewV2 creates JitJSON[T] from a value, with the default parser set.
func NewV2[T any](val T) *JitJSONV2[T] {
	return &JitJSONV2[T]{val: &val}
}

// NewFromBytesV2 creates a JitJSON[T] from an encoding, with the default parser set.
// If the encoding is invalid JSON, an error will be observed once Marshal is called.
func NewFromBytesV2[T any](data []byte) *JitJSONV2[T] {
	return &JitJSONV2[T]{data: data}
}

// Set sets a new value to JitJSON[T].
func (jit *JitJSONV2[T]) Set(val T) {
	jit.val = &val
	jit.data = nil
}

// SetBytes sets a new encoding to JitJSON[T].
func (jit *JitJSONV2[T]) SetBytes(data []byte) {
	jit.val = nil
	jit.data = data
}

// Marshal performs deferred json marshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Marshal' if the value has been marshaled previously. Once marshaled, the encoded value is stored with the
// jitjson for future retrieval. If there is no value to marshal, the method returns nil, nil.
func (jit *JitJSONV2[T]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}
	if jit.val == nil {
		return nil, nil
	}

	var err error
	jit.data, err = jsonv2.Marshal(jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

// Unmarshal performs deferred json unmarshaling for the value of JitJSON[T]. The method can return without evaluating
// 'json.Unmarshal' if the value has been unmarshaled previously. Once unmarshaled, the decoded value is stored with
// the jitjson for future retrieval. If there is no JSON data to unmarshal, the zero value of type T is returned.
// If the JSON data does not unmarshal into the type T, the method will return an error.
func (jit *JitJSONV2[T]) Unmarshal() (T, error) {
	if jit.val != nil {
		return *jit.val, nil
	}
	var val T
	if jit.data == nil {
		return val, nil
	}

	jit.val = &val
	err := jsonv2.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}

// MarshalJSON can be used to marshal JitJSON[T] to JSON.
// This method is compatable with the encoding/json/v1.Marshaler interface.
// When using encoding/json/v2, the json.MarshalJSONTo method is used to marshal instead.
func (jit *JitJSONV2[T]) MarshalJSON() ([]byte, error) {
	return jit.Marshal()
}

// UnmarshalJSON stores JSON data to be unmarshaled later.
// This method is compatable with the encoding/json/v1.Unmarshaler interface.
// When using encoding/json/v2, the json.UnmarshalJSONTo method is used to unmarshal instead.
func (jit *JitJSONV2[T]) UnmarshalJSON(data []byte) error {
	jit.val = nil
	jit.data = data
	return nil
}

// MarshalJSONTo implements the encoding/json/v2.MarshalerTo interface for efficient
// streaming JSON marshaling with json/v2. This method writes the JSON representation
// directly to the provided jsontext.Encoder, which is considered more efficient than
// the encoding/json/v1.MarshalJSON method.
func (jit *JitJSONV2[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	data, err := jit.Marshal()
	if err != nil {
		return err
	}
	if data == nil {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteValue(jsontext.Value(data))
}

// UnmarshalJSONFrom implements the encoding/json/v2.UnmarshalerFrom interface for
// more efficient decoding with JSON with json/v2 library. This method reads JSON
// data directly from the provided jsontext.Decoder, which is considered more efficient
// than the encoding/json/v1.UnmarshalJSON method.
func (jit *JitJSONV2[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	value, err := dec.ReadValue()
	if err != nil {
		return err
	}
	jit.val = nil
	jit.data = []byte(value)
	return nil
}
