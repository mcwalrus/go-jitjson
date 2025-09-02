//go:build go1.25

// Build with: $ GOEXPERIMENT=jsonv2 go run .
package main

import (
	"fmt"

	jsonv2 "encoding/json/v2"

	"github.com/mcwalrus/go-jitjson"
)

// implement JSONParser for encoding/json/v2.
type jsonParserV2 struct{}

var _ jitjson.JSONParser = (*jsonParserV2)(nil)

func (j *jsonParserV2) Name() string {
	return "encoding/json/v2"
}

func (j *jsonParserV2) Marshal(v interface{}) ([]byte, error) {
	return jsonv2.Marshal(v)
}

func (j *jsonParserV2) Unmarshal(data []byte, v interface{}) error {
	return jsonv2.Unmarshal(data, v)
}

type Person struct {
	Name string
	Age  int
	City string
}

// Uses encoding/json/v2 for marshaling.
func main() {

	v2Parser := &jsonParserV2{}
	jitjson.MustRegisterParser(v2Parser)
	jitjson.MustSetDefaultParser("encoding/json/v2")

	// Marshal using encoding/json/v2.
	jit := jitjson.New(Person{
		Name: "John",
		Age:  30,
		City: "New York",
	})

	jsonEncoding, err := jit.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonEncoding))

	// Unmarshal using encoding/json/v2.
	jit = jitjson.NewFromBytes[Person]([]byte(`{
		"Name": "John",
		"Age":  30,
		"City": "New York"
	}`))

	person, err := jit.Unmarshal()
	if err != nil {
		panic(err)
	}
	fmt.Println(person)

	// Verify the default parser used.
	parser := jit.Parser()
	fmt.Println("parser:", parser)
	fmt.Println("default parser:", jitjson.DefaultParser())
}
