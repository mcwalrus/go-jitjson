package jitjson

import (
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
			var any AnyJitJSON
			err := any.UnmarshalJSON([]byte(tt.input))
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}
			if got := any.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
		})
	}
}
