//go:build go1.25

package jitjson

import (
	jsonv2 "encoding/json/v2"
)

func init() {
	var v2Parser JSONParser = &jsonParserV2{}
	parsers[v2Parser.Name()] = v2Parser
}

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
