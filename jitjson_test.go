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

	jit := jitjson.NewJitJSON[Person](person)
	if jit == nil {
		t.Error("unexpected nil value")
	}
}

func TestBytesToJitJSON(t *testing.T) {
	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	jit := jitjson.BytesToJitJSON[Person](jsonData)
	if jit == nil {
		t.Error("unexpected nil value")
	}
}

func TestJitJSON_Nil(t *testing.T) {

	t.Run("nil value pointer", func(t *testing.T) {
		jit := jitjson.NewJitJSON[*int](nil)

		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if data != nil {
			t.Error("unexpected data")
		}

		val, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if val != nil {
			t.Error("unexpected value")
		}
	})

	t.Run("nil value non pointer", func(t *testing.T) {
		jit := jitjson.NewJitJSON[int](0)

		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}
		if data != nil {
			t.Error("unexpected data")
		}

		val, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}
		if val != 0 {
			t.Error("unexpected value")
		}
	})
}

func TestJitJSON_Marshal(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}

	jit := jitjson.NewJitJSON[Person](person)

	p1, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	p2, err := json.Marshal(jit)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(p1, p2) {
		t.Error("values do not match")
	}
}

func TestJitJSON_Unmarshal(t *testing.T) {
	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	jit := jitjson.BytesToJitJSON[Person](jsonData)

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

// func TestAnyJitJSON(t *testing.T) {

// 	// basic types
// 	var arr = []jitjson.AnyJit{
// 		jitjson.NewJitJSON[int](1),
// 		jitjson.NewJitJSON[float64](2.0),
// 		jitjson.NewJitJSON[string]("it works!"),
// 	}

// 	for _, v := range arr {
// 		switch v := v.(type) {

// 		case *jitjson.JitJSON[int]:
// 			i, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if i != 1 {
// 				t.Error("unexpected value")
// 			}

// 		case *jitjson.JitJSON[float64]:
// 			f, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if f != 2.0 {
// 				t.Error("unexpected value")
// 			}

// 		case *jitjson.JitJSON[string]:
// 			s, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if s != "it works!" {
// 				t.Error("unexpected value")
// 			}

// 		default:
// 			t.Error("unexpected type")
// 		}
// 	}

// 	// more types
// 	arr = []jitjson.AnyJitJSON{
// 		jitjson.NewJitJSON[bool](true),
// 		jitjson.NewJitJSON[[]int]([]int{1, 2, 3}),
// 		jitjson.NewJitJSON[map[string]string](map[string]string{"key": "value"}),
// 	}

// 	for _, v := range arr {
// 		switch v := v.(type) {

// 		case *jitjson.JitJSON[bool]:
// 			b, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if b != true {
// 				t.Error("unexpected value")
// 			}

// 		case *jitjson.JitJSON[[]int]:
// 			arr, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if len(arr) != 3 || arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
// 				t.Error("unexpected value")
// 			}

// 		case *jitjson.JitJSON[map[string]string]:
// 			m, err := v.Unmarshal()
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			if len(m) != 1 || m["key"] != "value" {
// 				t.Error("unexpected value")
// 			}

// 		default:
// 			t.Error("unexpected type")
// 		}
// 	}
// }

func TestJitJSON_JsonDecoder(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)
		jit := jitjson.BytesToJitJSON[Person](jsonData)

		var p Person
		dec := json.NewDecoder(jit)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&p); err != nil {
			t.Error(err)
		}

		if p.Name != "John" || p.Age != 30 || p.City != "New York" {
			t.Error("values do not match")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonData := []byte(`{"Name":"John","Age":30,"City":"New York","Country":"USA"}`)
		jit := jitjson.BytesToJitJSON[Person](jsonData)

		var p Person
		dec := json.NewDecoder(jit)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&p); err == nil {
			t.Error("expected error")
		}
	})
}

func TestJitJSON_UnmarshalMap(t *testing.T) {
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

func TestJitJSON_UnmarshalSlice(t *testing.T) {
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
