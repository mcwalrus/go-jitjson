package jitjson

import (
	"encoding/json"
	"testing"
)

func TestValueType_String(t *testing.T) {
	tests := []struct {
		vt   ValueType
		want string
	}{
		{TypeNull, "TypeNull"},
		{TypeBool, "TypeBool"},
		{TypeNumber, "TypeNumber"},
		{TypeString, "TypeString"},
		{TypeArray, "TypeArray"},
		{TypeObject, "TypeObject"},
		{TypeInvalid, "TypeInvalid"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.vt.String(); got != tt.want {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyJitJSON_NewAny(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		data := []byte(`{"key": "value"}`)
		a, err := NewAny(data)
		if err != nil {
			t.Errorf("error = %v", err)
		}
		if a == nil {
			t.Error("returned nil")
		}
		if a.Type() != TypeObject {
			t.Errorf("type = %v, want %v", a.Type(), TypeObject)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(`{invalid}`)
		a, err := NewAny(data)
		if err == nil {
			t.Errorf("expected error")
		}
		if a != nil {
			t.Error("should return nil on error")
		}
	})

	t.Run("null JSON", func(t *testing.T) {
		data := []byte(`null`)
		a, err := NewAny(data)
		if err != nil {
			t.Errorf("error = %v", err)
		}
		if !a.IsNull() {
			t.Error("should create null value")
		}
	})
}

func TestAnyJitJSON_MarshalJSON(t *testing.T) {
	t.Run("marshal original data", func(t *testing.T) {
		originalData := []byte(`{"key": "value", "num": 42}`)
		var a AnyJitJSON
		err := a.UnmarshalJSON(originalData)
		if err != nil {
			t.Fatal(err)
		}

		marshaled, err := a.MarshalJSON()
		if err != nil {
			t.Error(err)
		}

		if string(marshaled) != string(originalData) {
			t.Errorf("marshal result do not match")
		}
	})

	t.Run("marshal null", func(t *testing.T) {
		originalData := []byte(`null`)
		var a AnyJitJSON
		err := a.UnmarshalJSON(originalData)
		if err != nil {
			t.Fatal(err)
		}

		marshaled, err := a.MarshalJSON()
		if err != nil {
			t.Error(err)
		}

		if string(marshaled) != "null" {
			t.Errorf("marshal result do not match")
		}
	})
}

func TestAnyJitJSON_UnmarshalJSON(t *testing.T) {
	// type matching tests
	tests := []struct {
		name     string
		input    string
		wantType ValueType
	}{
		// Simple types
		{
			name:     "null",
			input:    `null`,
			wantType: TypeNull,
		},
		{
			name:     "bool",
			input:    `true`,
			wantType: TypeBool,
		},
		{
			name:     "number",
			input:    `123`,
			wantType: TypeNumber,
		},
		{
			name:     "string",
			input:    `"hello"`,
			wantType: TypeString,
		},
		// Complex types that could be mistaken for simple types
		{
			name:     "array of bools",
			input:    `[true, false]`,
			wantType: TypeArray,
		},
		{
			name:     "array of numbers",
			input:    `[123, 456]`,
			wantType: TypeArray,
		},
		{
			name:     "array of strings",
			input:    `["hello", "world"]`,
			wantType: TypeArray,
		},
		// Objects that could be mistaken for strings
		{
			name:     "object with string keys",
			input:    `{"true": 1}`,
			wantType: TypeObject,
		},
		{
			name:     "object with number-like keys",
			input:    `{"123": 1}`,
			wantType: TypeObject,
		},
		// Nested structures
		{
			name:     "nested array",
			input:    `[[1,2], [3,4]]`,
			wantType: TypeArray,
		},
		{
			name:     "nested object",
			input:    `{"a": {"b": 1}}`,
			wantType: TypeObject,
		},

		// Mixed content
		{
			name:     "mixed array",
			input:    `[1, "text", true, null]`,
			wantType: TypeArray,
		},
		{
			name:     "mixed object",
			input:    `{"num": 1, "str": "text", "bool": true, "null": null}`,
			wantType: TypeObject,
		},
		// Whitespace handling
		{
			name:     "padded array",
			input:    `  [1,2,3]  `,
			wantType: TypeArray,
		},
		{
			name:     "padded object",
			input:    `  {"a":1}  `,
			wantType: TypeObject,
		},

		// Multi-line content
		{
			name:     "multiline array",
			input:    "[\n  1,\n  2\n]",
			wantType: TypeArray,
		},
		{
			name:     "multiline object",
			input:    "{\n  \"a\": 1\n}",
			wantType: TypeObject,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AnyJitJSON
			err := a.UnmarshalJSON([]byte(tt.input))
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}
			if got := a.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
		})
	}

	// failure cases
	t.Run("invalid JSON", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`{invalid`))
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("empty input", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(``))
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`   `))
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("malformed string", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"unclosed string`))
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("malformed number", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`123.456.789`))
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestAnyJitJSON_Type(t *testing.T) {
	t.Run("nil AnyJitJSON", func(t *testing.T) {
		var a *AnyJitJSON
		if a.Type() != TypeNull {
			t.Errorf("Type() on nil = %v, want %v", a.Type(), TypeNull)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		a := &AnyJitJSON{val: "invalid type"}
		if a.Type() != TypeInvalid {
			t.Errorf("Type() on invalid = %v, want %v", a.Type(), TypeInvalid)
		}
	})

	t.Run("all valid types", func(t *testing.T) {
		testCases := []struct {
			name     string
			json     string
			expType  ValueType
			checkVal func(t *testing.T, a *AnyJitJSON)
		}{
			{
				name:    "null",
				json:    "null",
				expType: TypeNull,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					if !a.IsNull() {
						t.Error("expected null value")
					}
				},
			},
			{
				name:    "boolean",
				json:    "true",
				expType: TypeBool,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					boolVal, ok := a.AsBool()
					if !ok {
						t.Fatal("expected boolean type")
					}
					if !boolVal {
						t.Error("expected true")
					}
				},
			},
			{
				name:    "number",
				json:    "42",
				expType: TypeNumber,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					numVal, ok := a.AsNumber()
					if !ok {
						t.Fatal("expected number type")
					}
					if numVal != json.Number("42") {
						t.Errorf("expected 42, got %v", numVal)
					}
				},
			},
			{
				name:    "string",
				json:    `"hello"`,
				expType: TypeString,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					strVal, ok := a.AsString()
					if !ok {
						t.Fatal("expected string type")
					}
					if strVal != "hello" {
						t.Errorf("expected 'hello', got %v", strVal)
					}
				},
			},
			{
				name:    "array",
				json:    "[1,2,3]",
				expType: TypeArray,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					arr, ok := a.AsArray()
					if !ok {
						t.Fatal("expected array type")
					}
					if len(arr) != 3 {
						t.Errorf("expected length 3, got %d", len(arr))
					}
				},
			},
			{
				name:    "object",
				json:    `{"key":"value"}`,
				expType: TypeObject,
				checkVal: func(t *testing.T, a *AnyJitJSON) {
					t.Helper()
					obj, ok := a.AsObject()
					if !ok {
						t.Fatal("expected object type")
					}
					if len(obj) != 1 {
						t.Errorf("expected length 1, got %d", len(obj))
					}
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var a AnyJitJSON
				if err := json.Unmarshal([]byte(tc.json), &a); err != nil {
					t.Fatal(err)
				}
				if a.Type() != tc.expType {
					t.Errorf("expected %v, got %v", tc.expType, a.Type())
				}
				tc.checkVal(t, &a)
			})
		}
	})
}

func TestAnyJitJSON_String(t *testing.T) {
	t.Run("valid JSON pretty print", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`{"key":"value","number":42}`))
		if err != nil {
			t.Fatal(err)
		}

		result := a.String()
		expected := "{\n  \"key\": \"value\",\n  \"number\": 42\n}"
		if result != expected {
			t.Errorf("String() = %q, want %q", result, expected)
		}
	})

	t.Run("invalid JSON returns raw", func(t *testing.T) {
		a := &AnyJitJSON{data: []byte(`{invalid`)}
		result := a.String()
		expected := "{invalid"
		if result != expected {
			t.Errorf("String() = %q, want %q", result, expected)
		}
	})

	t.Run("simple value", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"hello"`))
		if err != nil {
			t.Fatal(err)
		}

		result := a.String()
		expected := "\"hello\""
		if result != expected {
			t.Errorf("String() = %q, want %q", result, expected)
		}
	})
}

func TestAnyJitJSON_AsArray(t *testing.T) {
	t.Run("array of numbers", func(t *testing.T) {
		data := `{"numbers": [1, 2, 3, 4, 5]}`
		var a AnyJitJSON
		if err := json.Unmarshal([]byte(data), &a); err != nil {
			t.Fatal(err)
		}

		obj, ok := a.AsObject()
		if !ok {
			t.Fatal("expected object type")
		}
		if a.Type() != TypeObject {
			t.Errorf("expected TypeObject, got %v", a.Type())
		}

		arr, ok := obj["numbers"].AsArray()
		if !ok {
			t.Fatal("expected array type")
		}
		if obj["numbers"].Type() != TypeArray {
			t.Errorf("expected TypeArray, got %v", obj["numbers"].Type())
		}

		for _, num := range arr {
			if num.Type() != TypeNumber {
				t.Errorf("expected TypeNumber, got %v", num.Type())
			}
			if _, ok := num.AsNumber(); !ok {
				t.Error("expected number type")
			}
		}
	})

	t.Run("nested objects", func(t *testing.T) {
		data := `{"outer": {"inner": {"str": "hello", "num": 42}}}`
		var a AnyJitJSON
		if err := json.Unmarshal([]byte(data), &a); err != nil {
			t.Fatal(err)
		}

		if a.Type() != TypeObject {
			t.Errorf("expected TypeObject, got %v", a.Type())
		}

		obj, ok := a.AsObject()
		if !ok {
			t.Fatal("expected object type")
		}

		if obj["outer"].Type() != TypeObject {
			t.Errorf("expected TypeObject, got %v", obj["outer"].Type())
		}
	})

	t.Run("mixed type array", func(t *testing.T) {
		data := `[null, 42, "text", [1,2], {"key": "value"}]`
		var a AnyJitJSON
		if err := json.Unmarshal([]byte(data), &a); err != nil {
			t.Fatal(err)
		}

		if a.Type() != TypeArray {
			t.Errorf("expected TypeArray, got %v", a.Type())
		}

		arr, ok := a.AsArray()
		if !ok {
			t.Fatal("expected array type")
		}

		expectedTypes := []ValueType{TypeNull, TypeNumber, TypeString, TypeArray, TypeObject}
		for i, arrItem := range arr {
			if arrItem.Type() != expectedTypes[i] {
				t.Errorf("index %d: expected %v, got %v", i, expectedTypes[i], arrItem.Type())
			}
		}
	})
}

func TestAnyJitJSON_As_ErrorCases(t *testing.T) {
	t.Run("AsBool on non-bool", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"not a bool"`))
		if err != nil {
			t.Fatal(err)
		}

		_, ok := a.AsBool()
		if ok {
			t.Error("AsBool() should return false for non-bool value")
		}
	})

	t.Run("AsNumber on non-number", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"not a number"`))
		if err != nil {
			t.Fatal(err)
		}

		_, ok := a.AsNumber()
		if ok {
			t.Error("AsNumber() should return false for non-number value")
		}
	})

	t.Run("AsString on non-string", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`42`))
		if err != nil {
			t.Fatal(err)
		}

		_, ok := a.AsString()
		if ok {
			t.Error("AsString() should return false for non-string value")
		}
	})

	t.Run("AsArray on non-array", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"not an array"`))
		if err != nil {
			t.Fatal(err)
		}

		_, ok := a.AsArray()
		if ok {
			t.Error("AsArray() should return false for non-array value")
		}
	})

	t.Run("AsObject on non-object", func(t *testing.T) {
		var a AnyJitJSON
		err := a.UnmarshalJSON([]byte(`"not an object"`))
		if err != nil {
			t.Fatal(err)
		}

		_, ok := a.AsObject()
		if ok {
			t.Error("AsObject() should return false for non-object value")
		}
	})
}
