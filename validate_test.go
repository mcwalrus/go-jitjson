package jitjson_test

import (
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

func TestValidate_ValidTypes(t *testing.T) {
	type ValidStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name string
		test func() error
	}{
		{
			name: "Valid Struct",
			test: func() error { return jitjson.Validate[ValidStruct]() },
		},
		{
			name: "Valid Map",
			test: func() error { return jitjson.Validate[map[string]interface{}]() },
		},
		{
			name: "Valid Map with int keys",
			test: func() error { return jitjson.Validate[map[int]string]() },
		},
		{
			name: "Valid Slice",
			test: func() error { return jitjson.Validate[[]string]() },
		},
		{
			name: "Valid Primitive",
			test: func() error { return jitjson.Validate[string]() },
		},
		{
			name: "Builtin String",
			test: func() error { return jitjson.Validate[string]() },
		},
		{
			name: "Builtin Int",
			test: func() error { return jitjson.Validate[int]() },
		},
		{
			name: "Builtin Float64",
			test: func() error { return jitjson.Validate[float64]() },
		},
		{
			name: "Builtin Bool",
			test: func() error { return jitjson.Validate[bool]() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.test(); err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
		})
	}

	// Edge cases

	t.Run("EmptyStruct", func(t *testing.T) {
		type EmptyStruct struct{}
		if err := jitjson.Validate[EmptyStruct](); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("PointerType", func(t *testing.T) {
		if err := jitjson.Validate[*ValidStruct](); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("InterfaceType", func(t *testing.T) {
		if err := jitjson.Validate[interface{}](); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})
}

func TestValidate_InvalidTypes(t *testing.T) {
	// channels cannot be marshaled to JSON
	type InvalidStruct struct {
		Channel chan struct{} `json:"channel"`
	}

	tests := []struct {
		name string
		test func() error
	}{
		{
			name: "InvalidStruct",
			test: func() error { return jitjson.Validate[InvalidStruct]() },
		},
		{
			name: "Invalid Map",
			test: func() error { return jitjson.Validate[map[struct{}]string]() },
		},
		{
			name: "ChannelType",
			test: func() error { return jitjson.Validate[chan int]() },
		},
		{
			name: "FunctionType",
			test: func() error { return jitjson.Validate[func()]() },
		},
		{
			name: "ComplexType",
			test: func() error { return jitjson.Validate[complex64]() },
		},
		// Dumb cases which should be addressed by the standard library...
		// {
		// 	name: "Invalid Slice",
		// 	test: func() error { return jitjson.Validate[[]chan struct{}]() },
		// },
		// {
		// 	name: "Invalid Slice of Functions",
		// 	test: func() error { return jitjson.Validate[[]func()]() },
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing %s", tt.name)
			if err := tt.test(); err == nil {
				t.Errorf("Validate() error = nil, want non-nil error")
			}
		})
	}
}

func TestValidate_ErrorMessages(t *testing.T) {
	tests := []struct {
		name           string
		test           func() error
		expectedPrefix string
	}{
		{
			name:           "ChannelError",
			test:           func() error { return jitjson.Validate[chan int]() },
			expectedPrefix: "type chan int is not JSON marshalable",
		},
		{
			name:           "FunctionError",
			test:           func() error { return jitjson.Validate[func()]() },
			expectedPrefix: "type func() is not JSON marshalable",
		},
		{
			name:           "ComplexError",
			test:           func() error { return jitjson.Validate[complex64]() },
			expectedPrefix: "type complex64 is not JSON marshalable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.test()
			if err == nil {
				t.Errorf("Validate() error = nil, want non-nil error")
				return
			}

			if err.Error()[:len(tt.expectedPrefix)] != tt.expectedPrefix {
				t.Errorf("Validate() error = %v, want error with prefix %q", err, tt.expectedPrefix)
			}
		})
	}
}
