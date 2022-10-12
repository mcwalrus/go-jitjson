package jitjson_test

import (
	"bytes"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

func TestNilJitJSON(t *testing.T) {
	data := []byte(`null`)

	jit, err := jitjson.NewJitJSON[struct{}](data)
	if err != nil {
		t.Error("failed to create struct{} JitJSON")
		t.FailNow()
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != struct{}{} {
		t.Error("unexpected value from Unmarshal")
	}
}

func TestIntJitJSON(t *testing.T) {
	data := []byte(`1`)

	jit, err := jitjson.NewJitJSON[int](data)
	if err != nil {
		t.Error("failed to create int JitJSON")
		t.FailNow()
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != 1 {
		t.Error("unexpected value from Unmarshal")
	}
}

func TestFloatJitJSON(t *testing.T) {
	data := []byte(`1.00000001`)

	jit, err := jitjson.NewJitJSON[float64](data)
	if err != nil {
		t.Error("failed to create float64 JitJSON")
		t.FailNow()
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != 1.00000001 {
		t.Error("unexpected value from Unmarshal")
	}
}

func TestStringJitJSON(t *testing.T) {
	data := []byte(`"json encoded"`)

	jit, err := jitjson.NewJitJSON[string](data)
	if err != nil {
		t.Error("failed to create string JitJSON")
		t.FailNow()
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
		t.FailNow()
	}
	if value != "json encoded" {
		t.Error("unexpected value from Unmarshal")
		t.FailNow()
	}
}

func TestArrayJitJSON(t *testing.T) {
	data := []byte(`[1, 2, 3, 4, 5]`)

	jit, err := jitjson.NewJitJSON[[]int](data)
	if err != nil {
		t.Error("failed to create []int JitJSON")
		t.FailNow()
	}

	values, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
		t.FailNow()
	}

	expected := []int{1, 2, 3, 4, 5}
	if !compareValues(t, expected, values) {
		t.Error("unexpected values")
	}
}

func TestMapJitJSON(t *testing.T) {
	data := []byte(`{
		"1": "two",
		"three": 4.0,
		"5": [6, 7, "eight", 9, "ten"]
	}`)

	jit, err := jitjson.NewJitJSON[map[string]interface{}](data)
	if err != nil {
		t.Error("failed to create map JitJSON")
		t.FailNow()
	}

	m, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
		t.FailNow()
	} else if m == nil {
		t.Error("expected map to be filled")
		t.FailNow()
	}

	value, ok := m["1"]
	if !ok {
		t.Error("expected map to contain key \"1\"")
	} else if value != "two" {
		t.Error("expected value to equal \"two\"")
	}

	value, ok = m["three"]
	if !ok {
		t.Error("expected map to contain key \"three\"")
	} else if value != 4.0 {
		t.Error("expected value to equal 4.0")
	}

	value, ok = m["5"]
	if !ok {
		t.Error("expected map to contain key \"5\"")
		t.FailNow()
	}

	v, ok := (value).([]interface{})
	if !ok {
		t.Error("expected value to have type []interface{}")
		t.FailNow()
	}

	expected := []interface{}{6.0, 7.0, "eight", 9.0, "ten"}
	if !compareInterfaceValues(t, expected, v) {
		t.Error("unexpected values")
	}
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestExample1(t *testing.T) {
	data := []byte(`
        {
            "name": "Willy Wonka",
            "age":  42
        }
    `)

	expected := Person{
		Name: "Willy Wonka",
		Age:  42,
	}

	jit, err := jitjson.NewJitJSON[Person](data)
	if err != nil {
		t.Error("failed to create Person JitJSON")
		t.FailNow()
	}

	person, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if person != expected {
		t.Error("unexpected person struct")
	}
}

func TestExample2(t *testing.T) {
	person := Person{
		Name: "Charlie Bucket",
		Age:  12,
	}

	expected := []byte(`{"name":"Charlie Bucket","age":12}`)

	jit, err := jitjson.NewJitJSON[Person](nil)
	if err != nil {
		t.Error("failed to create nil Person JitJSON")
		t.FailNow()
	}

	jit.Set(person)
	data, err := jit.Marshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if bytes.Compare(expected, data) != 0 {
		t.Error("unexpected json encoding")
	}
}

type TestType struct {
	Field *string `json:"field"`
}

func TestNilFieldJitJSON(t *testing.T) {
	data := []byte(`{
		"field": null
	}`)

	jit, err := jitjson.NewJitJSON[TestType](data)
	if err != nil {
		t.Error("failed to create TestType JitJSON")
		t.FailNow()
	}

	var expected TestType
	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != expected {
		t.Error("unexpected TestType")
	}
}

func TestUpdateNilFieldJitJSON(t *testing.T) {
	data := []byte(`{
		"field": null
	}`)

	jit, err := jitjson.NewJitJSON[TestType](data)
	if err != nil {
		t.Error("failed to create TestType JitJSON")
		t.FailNow()
	}

	var expected TestType
	value, err := jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != expected {
		t.Error("unexpected TestType")
	}

	field := "some value"
	expected = TestType{
		Field: &field,
	}

	jit.Set(expected)
	value, err = jit.Unmarshal()
	if err != nil {
		t.Error("failed Unmarshal JitJSON")
	} else if value != expected {
		t.Error("unexpected TestType after Set")
	}
}

type testInterface interface {
	method() int
}

func TestErrorsJitJSON(t *testing.T) {

	// triggers pointer type error.
	t.Run("pointer", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[*int](nil)
		if err == nil {
			t.Error("expected error")
		}
	})

	// triggers invalid type error.
	t.Run("interface", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[interface{}](nil)
		if err == nil {
			t.Error("expected error")
		}
	})

	// triggers invalid type error.
	t.Run("invalid type", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[testInterface](nil)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func compareValues[T comparable](t *testing.T, expected, values []T) bool {
	if len(expected) != len(values) {
		return false
	}

	for i, expect := range expected {
		if expect != values[i] {
			return false
		}
	}

	return true
}

func compareInterfaceValues(t *testing.T, expected, values []interface{}) bool {
	if len(expected) != len(values) {
		return false
	}

	for i, expect := range expected {
		switch e := expect.(type) {
		case int:
			v, ok := (values[i]).(int)
			if !ok {
				return false
			}
			if e != v {
				return false
			}
		case float64:
			v, ok := (values[i]).(float64)
			if !ok {
				return false
			}
			if e != v {
				return false
			}
		case string:
			v, ok := (values[i]).(string)
			if !ok {
				return false
			}
			if e != v {
				return false
			}
		default:
			return false
		}
	}

	return true
}
