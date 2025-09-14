//go:build go1.25 && goexperiment.jsonv2

// Build with: $ GOEXPERIMENT=jsonv2 go run .
package main

import (
	"fmt"

	"encoding/json"
	jsonv2 "encoding/json/v2"

	"github.com/mcwalrus/go-jitjson"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

// Example usage of encoding/json/v2 parser.
func main() {
	jitjson.MustSetDefaultParser("encoding/json/v2")
	verifyDefaultParser()

	simpleExample()
	jsonv2Example()
	invalidJsonv2Example()
}

// Verify default parser used.
func verifyDefaultParser() {
	fmt.Println("--------------------------------")
	fmt.Println("Default parser")
	fmt.Println("--------------------------------")

	// Verify the default parser used
	var jit = jitjson.New(Person{})
	parser := jit.Parser()
	fmt.Println("parser:", parser)
	fmt.Println("default parser:", jitjson.DefaultParser())
}

// Simple example of using encoding/json/v2 parser.
func simpleExample() {
	fmt.Println("--------------------------------")
	fmt.Println("Simple example")
	fmt.Println("--------------------------------")

	// Marshal
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

	// Unmarshal
	jit = jitjson.NewFromBytes[Person]([]byte(`{
			"name": "John",
			"age":  30,
			"city": "New York"
		}`))

	person, err := jit.Unmarshal()
	if err != nil {
		panic(err)
	}
	fmt.Println(person)
}

// Explains details how to use encoding/json/v2 parser.
func jsonv2Example() {

	fmt.Println("--------------------------------")
	fmt.Println("Marshaling")
	fmt.Println("--------------------------------")

	person := Person{
		Name: "John Doe",
		Age:  30,
		City: "New York",
	}

	// json/v1.Marshal calls jit.MarshalJSON()
	stdJSON, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json/v1.Marshal: %s\n", stdJSON)

	jit := jitjson.New(person)

	// json/v2.Marshal calls jit.MarshalJSONTo()
	jitJSON, err := jsonv2.Marshal(jit)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json/v2.Marshal: %s\n", jitJSON)

	fmt.Println("--------------------------------")
	fmt.Println("Unmarshaling")
	fmt.Println("--------------------------------")

	jsonData := `{
		"name":"Jane Smith",
		"age":28,
		"city":"San Francisco"
	}`

	// json/v1.Unmarshal calls jit2.UnmarshalJSON()
	var stdPerson Person
	err = json.Unmarshal([]byte(jsonData), &stdPerson)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json/v1.Unmarshal: %+v\n", stdPerson)

	// json/v2.Unmarshal calls jit2.UnmarshalJSONFrom()
	var jit2 jitjson.JitJSON[Person]
	err = jsonv2.Unmarshal([]byte(jsonData), &jit2)
	if err != nil {
		panic(err)
	}
	result, err := jit2.Unmarshal()
	if err != nil {
		panic(err)
	}
	fmt.Printf("json/v2.Unmarshal: %+v\n", result)
}

func invalidJsonv2Example() {
	fmt.Println("--------------------------------")
	fmt.Println("Invalid JSON parsing")
	fmt.Println("--------------------------------")

	// Data has incorrect case-sensitivity for field names
	// This would work with json/v1 but has been addressed in json/v2
	// For more information, see the offical blog post: https://go.dev/blog/jsonv2-exp
	jsonData := `{
		"Name": "John",
		"Age":  30,
		"City": "New York"
	}`

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}

	// unmarshal with json/v1
	var jit = jitjson.JitJSON[Person]{}
	jit.SetParser("encoding/json")

	err := json.Unmarshal([]byte(jsonData), &jit)
	if err != nil {
		panic(err)
	}
	person, err := jit.Unmarshal()
	if err != nil {
		panic(err)
	}
	fmt.Println("json/v1 failed to unmarshal:", Person{} == person, "parser:", jit.Parser())

	// unmarshal with json/v2
	var jit2 = jitjson.JitJSON[Person]{}
	jit2.SetParser("encoding/json/v2")

	err = jsonv2.Unmarshal([]byte(jsonData), &jit2)
	if err != nil {
		panic(err)
	}

	person, err = jit2.Unmarshal()
	if err != nil {
		panic(err)
	}
	fmt.Println("json/v2 failed to unmarshal:", Person{} == person, "parser:", jit2.Parser())

}
