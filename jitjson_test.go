package jitjson_test

import (
	"bytes"
	"testing"

	"encoding/json"

	"github.com/mcwalrus/go-jitjson"
)

type Person struct {
	Name string
	Age  int
	City string
}

func TestNewJitJSON(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}

	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	t.Run("Marshal Person", func(t *testing.T) {
		jit := jitjson.New(person)

		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(data, jsonData) {
			t.Error("data do not match")
		}
	})

	t.Run("Decode Person", func(t *testing.T) {
		jit := jitjson.NewFromBytes[Person](jsonData)

		p1, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}

		if p1.Name != person.Name || p1.Age != person.Age || p1.City != person.City {
			t.Error("values do not match")
		}
	})
}

func TestJitJSON_Nil(t *testing.T) {
	jit := jitjson.New[*int](nil)

	t.Run("Marshal nil", func(t *testing.T) {
		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if string(data) != "null" {
			t.Error("expected null")
		}
	})

	t.Run("Decode nil", func(t *testing.T) {
		val, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if val != nil {
			t.Error("expected nil")
		}
	})
}

func TestJitJSON_Slice(t *testing.T) {
	jsonData := []byte(`[
		{"Name":"John","Age":30,"City":"New York"},
		{"Name":"Jane","Age":25,"City":"Los Angeles"}
	]`)

	var result []jitjson.JitJSON[Person]
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 elements, got %d", len(result))
	}

	person1, err := result[0].Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if person1.Name != "John" || person1.Age != 30 || person1.City != "New York" {
		t.Error("values do not match for person1")
	}

	person2, err := result[1].Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if person2.Name != "Jane" || person2.Age != 25 || person2.City != "Los Angeles" {
		t.Error("values do not match for person2")
	}
}

func TestJitJSON_Map(t *testing.T) {
	jsonData := []byte(`{
		"person1": {"Name":"John","Age":30,"City":"New York"},
		"person2": {"Name":"Jane","Age":25,"City":"Los Angeles"}
	}`)

	var result map[string]*jitjson.JitJSON[Person]
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Error(err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 elements, got %d", len(result))
	}

	person1, err := result["person1"].Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if person1.Name != "John" || person1.Age != 30 || person1.City != "New York" {
		t.Error("values do not match for person1")
	}

	person2, err := result["person2"].Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if person2.Name != "Jane" || person2.Age != 25 || person2.City != "Los Angeles" {
		t.Error("values do not match for person2")
	}
}
