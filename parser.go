// go:build go1.25

package jitjson

import (
	"fmt"

	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
)

// ParserVersion is the version of the JSON library to use.
type ParserVersion int

const (
	JsonV1 ParserVersion = 1
	JsonV2 ParserVersion = 2
)

// parser is used by jitjson.go.
// No synchronization is applied for performance reasons.
var (
	parser jsonParser = jsonV1{}
)

// defaultParser returns the v2 json parser.
func defaultParser() jsonParser {
	return parser
}

func parserFromVersion(version ParserVersion) (jsonParser, error) {
	if version == JsonV1 {
		return jsonV1{}, nil
	} else if version == JsonV2 {
		return jsonV2{}, nil
	} else {
		return nil, fmt.Errorf("unsupported json version: %d", version)
	}
}

// SetDefaultParser sets the encoding/json version parser to use.
// By default, the library uses the encoding/json/v2 parser.
// If the version is not supported, the function returns an error.
func SetDefaultParser(version ParserVersion) error {
	p, err := parserFromVersion(version)
	if err != nil {
		return err
	}
	parser = p
	return nil
}

// MustSetDefaultParser sets the default parser and panics if the version is not supported.
// This is a convenience function for backward compatibility.
func MustSetDefaultParser(version ParserVersion) {
	if err := SetDefaultParser(version); err != nil {
		panic(err)
	}
}

type jsonParser interface {
	name() string
	marshal(v interface{}) ([]byte, error)
	unmarshal(data []byte, v interface{}) error
}

type jsonV1 struct{}

func (j jsonV1) marshal(v interface{}) ([]byte, error) {
	return jsonv1.Marshal(v)
}

func (j jsonV1) unmarshal(data []byte, v interface{}) error {
	return jsonv1.Unmarshal(data, v)
}

func (j jsonV1) name() string {
	return "encoding/json"
}

type jsonV2 struct{}

func (j jsonV2) marshal(v interface{}) ([]byte, error) {
	return jsonv2.Marshal(v)
}

func (j jsonV2) unmarshal(data []byte, v interface{}) error {
	return jsonv2.Unmarshal(data, v)
}

func (j jsonV2) name() string {
	return "encoding/json/v2"
}

// GetDefaultParser returns the name of the currently configured default parser.
func GetDefaultParser() string {
	return parser.name()
}
