//go:build go1.25 && goexperiment.jsonv2

package jitjson_test

import (
	"strings"
	"testing"

	"encoding/json"
	"encoding/json/jsontext"

	"github.com/mcwalrus/go-jitjson"
)

func TestJitJSONV2(t *testing.T) {
	// Test basic functionality without parser management
}

func TestJitJSON_MarshalJSONTo(t *testing.T) {
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

func TestJitJSON_UnmarshalJSONFrom(t *testing.T) {
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

// Helper function to compare JSON strings semantically
func jsonEqual(a, b string) bool {
	var objA, objB interface{}
	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return false
	}
	dataA, errA := json.Marshal(objA)
	dataB, errB := json.Marshal(objB)
	if errA != nil || errB != nil {
		return false
	}
	return string(dataA) == string(dataB)
}
