package jitjson_test

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/mcwalrus/go-jitjson"
)

func TestNewJitJSON(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}

	_, err := jitjson.NewJitJSON[Person](person)
	if err != nil {
		t.Error(err)
	}
}

func TestMarshal(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}

	jit, err := jitjson.NewJitJSON[Person](person)
	if err != nil {
		t.Error(err)
	}

	p1, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	p2, err := json.Marshal(jit)
	if err != nil {
		t.Error(err)
	}

	if bytes.Equal(p1, p2) {
		t.Error("values do not match")
	}
}

func TestUnmarshal(t *testing.T) {
	type Person struct {
		Name string
		Age  int
		City string
	}

	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	jit, err := jitjson.NewJitJSON[Person](jsonData)
	if err != nil {
		t.Error(err)
	}

	p1, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	var p2 Person
	err = json.Unmarshal(jsonData, &p2)
	if err != nil {
		t.Error(err)
	}

	if p1 != p2 {
		t.Error("values do not match")
	}
}
