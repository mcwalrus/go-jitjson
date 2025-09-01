// go:build go1.25

package jitjson

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
)

type jsonType int

const (
	JsonV1 jsonType = 1
	JsonV2 jsonType = 2
)

// parser is used by jitjson.go.
var parser jsonParser = jsonV2{}

// SetParser sets the encoding/json version parser to use. By default, the json/v2 parser is used.
// If the version is not supported, the function panics.
func SetParser(version jsonType) {
	if version == JsonV1 {
		parser = jsonV1{}
	} else if version == JsonV2 {
		parser = jsonV2{}
	} else {
		panic("invalid json version")
	}
}

type jsonParser interface {
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

type jsonV2 struct{}

func (j jsonV2) marshal(v interface{}) ([]byte, error) {
	return jsonv2.Marshal(v)
}

func (j jsonV2) unmarshal(data []byte, v interface{}) error {
	return jsonv2.Unmarshal(data, v)
}
