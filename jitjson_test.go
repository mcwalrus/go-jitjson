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

func TestJITInterface(t *testing.T) {
	var (
		err error
		arr = make([]jitjson.JITInterface, 3)
	)

	arr[0], err = jitjson.NewJitJSON[int](1)
	if err != nil {
		t.Error(err)
	}

	arr[1], err = jitjson.NewJitJSON[float64](2.0)
	if err != nil {
		t.Error(err)
	}

	arr[2], err = jitjson.NewJitJSON[string]("it works!")
	if err != nil {
		t.Error(err)
	}

	for _, v := range arr {
		switch v := v.(type) {

		case jitjson.JitJSON[int]:
			i, err := v.Unmarshal()
			if err != nil {
				t.Error(err)
			}
			if i != 1 {
				t.Error("unexpected value")
			}

		case jitjson.JitJSON[float64]:
			f, err := v.Unmarshal()
			if err != nil {
				t.Error(err)
			}
			if f != 2.0 {
				t.Error("unexpected value")
			}

		case jitjson.JitJSON[string]:
			s, err := v.Unmarshal()
			if err != nil {
				t.Error(err)
			}
			if s != "it works!" {
				t.Error("unexpected value")
			}

		default:
			t.Error("unexpected type")
		}
	}
}

func TestJSONDecoder(t *testing.T) {
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

	var p Person
	dec := json.NewDecoder(jit)
	dec.DisallowUnknownFields()
	err = dec.Decode(&p)
	if err != nil {
		t.Error(err)
	}

	if p.Name != "John" || p.Age != 30 || p.City != "New York" {
		t.Error("values do not match")
	}
}
