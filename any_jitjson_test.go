package jitjson

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON_TypeMatching(t *testing.T) {
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

func TestAnyJitJSON_Type(t *testing.T) {

	t.Run("nil value", func(t *testing.T) {
		var a *AnyJitJSON
		if a.Type() != TypeNull {
			t.Errorf("expected TypeNull, got %v", a.Type())
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
