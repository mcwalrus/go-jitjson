//go:build go1.25 && goexperiment.jsonv2

package jitjson_test

import (
	"bytes"
	"strings"
	"testing"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/mcwalrus/go-jitjson"
)

func TestNewJitJSONV2(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}

	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	t.Run("Marshal Person", func(t *testing.T) {
		jit := jitjson.NewV2(person)

		data, err := jit.Marshal()
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(data, jsonData) {
			t.Error("data do not match")
		}
	})

	t.Run("Decode Person", func(t *testing.T) {
		jit := jitjson.NewFromBytesV2[Person](jsonData)

		p1, err := jit.Unmarshal()
		if err != nil {
			t.Error(err)
		}

		if p1.Name != person.Name || p1.Age != person.Age || p1.City != person.City {
			t.Error("values do not match")
		}
	})

	t.Run("Marshal result through jsonv2.Marshal", func(t *testing.T) {
		jit := jitjson.NewV2(person)
		data, err := jsonv2.Marshal(jit)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(data, jsonData) {
			t.Error("data do not match")
		}
	})

	t.Run("Unmarshal result through jsonv2.Unmarshal", func(t *testing.T) {
		var jit jitjson.JitJSON[Person]
		err := jsonv2.Unmarshal(jsonData, &jit)
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

}

// TestJitJSON_Set methods should provide consistency between the value and the encoding stored / returned.
func TestJitJSONV2_Set(t *testing.T) {
	person1 := Person{Name: "John", Age: 30, City: "New York"}
	person2 := Person{Name: "Jane", Age: 25, City: "Los Angeles"}

	person1Data := []byte(`{"Name":"John","Age":30,"City":"New York"}`)
	person2Data := []byte(`{"Name":"Jane","Age":25,"City":"Los Angeles"}`)

	t.Run("SetValue", func(t *testing.T) {
		jit := jitjson.NewFromBytesV2[Person](person1Data)
		jit.Set(person1)

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
		jit := jitjson.NewV2(person1)
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

func TestJitJSONV2_Nil(t *testing.T) {
	jit := jitjson.NewV2[*int](nil)

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

func TestJitJSONV2_Slice(t *testing.T) {
	jsonData := []byte(`[
		{"Name":"John","Age":30,"City":"New York"},
		{"Name":"Jane","Age":25,"City":"Los Angeles"}
	]`)

	var result []*jitjson.JitJSONV2[Person]
	err := jsonv2.Unmarshal(jsonData, &result)
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

func TestJitJSONV2_Map(t *testing.T) {
	jsonData := []byte(`{
		"person1": {"Name":"John","Age":30,"City":"New York"},
		"person2": {"Name":"Jane","Age":25,"City":"Los Angeles"}
	}`)

	var result map[string]*jitjson.JitJSONV2[Person]
	err := jsonv2.Unmarshal(jsonData, &result)
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

func TestJitJSONV2_MarshalJSONTo(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}
	expectedJSON := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	t.Run("MarshalJSONTo with value", func(t *testing.T) {
		jit := jitjson.NewV2(person)
		var buf strings.Builder
		enc := jsontext.NewEncoder(&buf)

		err := jit.MarshalJSONTo(enc)
		if err != nil {
			t.Errorf("MarshalJSONTo failed: %v", err)
		}

		result := buf.String()
		if !jsonEqual(result, string(expectedJSON)) {
			t.Errorf("Expected %s, got %s", string(expectedJSON), result)
		}
	})

	t.Run("MarshalJSONTo with cached data", func(t *testing.T) {
		jit := jitjson.NewFromBytesV2[Person](expectedJSON)
		var buf strings.Builder
		enc := jsontext.NewEncoder(&buf)

		err := jit.MarshalJSONTo(enc)
		if err != nil {
			t.Errorf("MarshalJSONTo failed: %v", err)
		}

		result := buf.String()
		if !jsonEqual(result, string(expectedJSON)) {
			t.Errorf("Expected %s, got %s", string(expectedJSON), result)
		}
	})

	t.Run("MarshalJSONTo with nil value", func(t *testing.T) {
		jit := jitjson.JitJSONV2[Person]{}
		var buf strings.Builder
		enc := jsontext.NewEncoder(&buf)

		err := jit.MarshalJSONTo(enc)
		if err != nil {
			t.Errorf("MarshalJSONTo failed: %v", err)
		}

		result := strings.TrimSpace(buf.String())
		expectedNull := "null"
		if result != expectedNull {
			t.Errorf("Expected %s, got %s", expectedNull, result)
		}
	})

	t.Run("MarshalJSONTo returns error on invalid json", func(t *testing.T) {
		jit := jitjson.JitJSONV2[Person]{}
		invalidJSON := []byte(`{"invalid": json}`)
		jit.SetBytes(invalidJSON)

		var buf strings.Builder
		enc := jsontext.NewEncoder(&buf)

		err := jit.MarshalJSONTo(enc)
		if err == nil {
			t.Error("MarshalJSONTo should fail with invalid json")
		}
	})
}

func TestJitJSONV2_UnmarshalJSONFrom(t *testing.T) {
	person := Person{
		Name: "John",
		Age:  30,
		City: "New York",
	}
	jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

	t.Run("UnmarshalJSONFrom basic", func(t *testing.T) {
		jit := jitjson.JitJSONV2[Person]{}
		dec := jsontext.NewDecoder(strings.NewReader(string(jsonData)))

		err := jit.UnmarshalJSONFrom(dec)
		if err != nil {
			t.Errorf("UnmarshalJSONFrom failed: %v", err)
		}

		result, err := jit.Unmarshal()
		if err != nil {
			t.Errorf("Unmarshal after UnmarshalJSONFrom failed: %v", err)
		}

		if result != person {
			t.Errorf("Expected %+v, got %+v", person, result)
		}
	})

	t.Run("UnmarshalJSONFrom clears existing value", func(t *testing.T) {
		existingPerson := Person{Name: "Jane", Age: 25, City: "LA"}
		jit := jitjson.NewV2(existingPerson)

		// Verify the first existing value
		result, err := jit.Unmarshal()
		if err != nil {
			t.Errorf("Initial unmarshal failed: %v", err)
		}
		if result != existingPerson {
			t.Error("Initial value not set correctly")
		}

		// Now unmarshal from decoder
		dec := jsontext.NewDecoder(strings.NewReader(string(jsonData)))
		err = jit.UnmarshalJSONFrom(dec)
		if err != nil {
			t.Errorf("UnmarshalJSONFrom failed: %v", err)
		}

		// Verify the new data overwrote the old value
		result, err = jit.Unmarshal()
		if err != nil {
			t.Errorf("Unmarshal after UnmarshalJSONFrom failed: %v", err)
		}

		if result != person {
			t.Errorf("Expected %+v, got %+v", person, result)
		}
	})

	t.Run("UnmarshalJSONFrom with invalid JSON", func(t *testing.T) {
		jit := jitjson.JitJSONV2[Person]{}
		invalidJSON := `{"invalid": json}`
		dec := jsontext.NewDecoder(strings.NewReader(invalidJSON))
		err := jit.UnmarshalJSONFrom(dec)
		if err == nil {
			t.Error("UnmarshalJSONFrom should fail with invalid json")
		}
	})
}

func jsonEqual(a, b string) bool {
	var objA, objB interface{}
	if err := jsonv2.Unmarshal([]byte(a), &objA); err != nil {
		return false
	}
	if err := jsonv2.Unmarshal([]byte(b), &objB); err != nil {
		return false
	}
	dataA, errA := jsonv2.Marshal(objA)
	dataB, errB := jsonv2.Marshal(objB)
	if errA != nil || errB != nil {
		return false
	}
	return string(dataA) == string(dataB)
}
