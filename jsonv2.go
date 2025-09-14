//go:build go1.25 && goexperiment.jsonv2

package jitjson

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// jsonParserV2 is a JSONParser using the encoding/json/v2 package.
type jsonParserV2 struct{}

var _ JSONParser = (*jsonParserV2)(nil)

func (j *jsonParserV2) Name() string {
	return "encoding/json/v2"
}

func (j *jsonParserV2) Marshal(v interface{}) ([]byte, error) {
	return jsonv2.Marshal(v)
}

func (j *jsonParserV2) Unmarshal(data []byte, v interface{}) error {
	return jsonv2.Unmarshal(data, v)
}

func init() {
	initParserRegistry()
	MustRegisterParser(&jsonParserV2{})
}

var _ jsonv2.MarshalerTo = (*JitJSON[any])(nil)
var _ jsonv2.UnmarshalerFrom = (*JitJSON[any])(nil)

// MarshalJSONTo implements the encoding/json/v2.MarshalerTo interface for efficient
// streaming JSON marshaling with json/v2. This method writes the JSON representation
// directly to the provided jsontext.Encoder.
func (jit *JitJSON[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
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
// efficient streaming JSON unmarshaling with json/v2. This method reads JSON data
// directly from the provided jsontext.Decoder.
func (jit *JitJSON[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	value, err := dec.ReadValue()
	if err != nil {
		return err
	}
	jit.val = nil
	jit.data = []byte(value)
	return nil
}
