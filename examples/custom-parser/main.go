package main

import (
	"fmt"

	"github.com/wI2L/jettison"

	jitjson "github.com/mcwalrus/go-jitjson"
)

type jettisonParser struct{}

var _ jitjson.JSONParser = (*jettisonParser)(nil)

func (c *jettisonParser) Name() string {
	return "jettison"
}

func (c *jettisonParser) Marshal(v interface{}) ([]byte, error) {
	return jettison.Marshal(v)
}

func (c *jettisonParser) Unmarshal(data []byte, v interface{}) error {
	return fmt.Errorf("not implemented by jettison parser")
}

type Person struct {
	Name string
	Age  int
	City string
}

// Uses jettison for marshaling only.
func main() {

	jetParser := &jettisonParser{}
	jitjson.MustRegisterParser(jetParser)
	jitjson.MustSetDefaultParser("jettison")

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

	jit = jitjson.NewFromBytes[Person](jsonEncoding)
	_, err = jit.Unmarshal()
	if err == nil {
		panic("expected unmarshal error")
	}

	fmt.Println("error unmarshalling:", err)

	// Verify the default parser used.
	parser := jit.Parser()
	fmt.Println("parser:", parser)
	fmt.Println("default parser:", jitjson.DefaultParser())
}
