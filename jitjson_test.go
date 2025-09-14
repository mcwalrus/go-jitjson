package jitjson_test

import (
	"bytes"
	"encoding/json"
	"testing"

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

	t.Run("Marshal result through json.Marshal", func(t *testing.T) {
		jit := jitjson.New(person)
		data, err := json.Marshal(jit)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, jsonData) {
			t.Error("data do not match")
		}
	})

	t.Run("Unmarshal result through json.Unmarshal", func(t *testing.T) {
		var jit jitjson.JitJSON[Person]
		err := json.Unmarshal(jsonData, &jit)
		if err != nil {
			t.Error(err)
		}
		value, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if value != person {
			t.Error("value do not match")
		}
	})

	t.Run("Marshal result nil without value set", func(t *testing.T) {
		jit := jitjson.JitJSON[Person]{}
		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if data != nil {
			t.Error("data should be nil")
		}
	})

	t.Run("Unmarshal result zero value without data set", func(t *testing.T) {
		jit := jitjson.JitJSON[Person]{}
		p, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if p != (Person{}) {
			t.Error("value should be zero value")
		}
	})

	t.Run("Verify default parser is set on marshal", func(t *testing.T) {
		jit := jitjson.JitJSON[Person]{}
		jit.SetValue(person)
		_, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if jit.Parser() != jitjson.DefaultParser() {
			t.Error("parser should be default parser")
		}
	})

	t.Run("Verify default parser is set on unmarshal", func(t *testing.T) {
		jit := jitjson.JitJSON[Person]{}
		jit.SetBytes(jsonData)
		_, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if jit.Parser() != jitjson.DefaultParser() {
			t.Error("parser should be default parser")
		}
	})
}

// TestJitJSON_Set methods should provide consistency between the value and the encoding stored / returned.
func TestJitJSON_Set(t *testing.T) {
	person1 := Person{Name: "John", Age: 30, City: "New York"}
	person2 := Person{Name: "Jane", Age: 25, City: "Los Angeles"}

	person1Data := []byte(`{"Name":"John","Age":30,"City":"New York"}`)
	person2Data := []byte(`{"Name":"Jane","Age":25,"City":"Los Angeles"}`)

	t.Run("SetValue", func(t *testing.T) {
		jit := jitjson.NewFromBytes[Person](person1Data)
		jit.SetValue(person1)

		//
		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, person1Data) {
			t.Error("data do not match")
		}

		value, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if value != person1 {
			t.Error("value do not match")
		}
	})

	t.Run("SetBytes", func(t *testing.T) {
		jit := jitjson.New(person1)
		jit.SetBytes(person2Data)
		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, person2Data) {
			t.Error("data do not match")
		}

		value, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if value != person2 {
			t.Error("value do not match")
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
