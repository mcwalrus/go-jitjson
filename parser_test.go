package jitjson

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// resetParserRegistry resets the parser registry to the default state.
func resetParserRegistry(t *testing.T) {
	t.Helper()
	parsers = make(map[string]JSONParser)
	setupParserRegistry()
}

// mockParser is a test implementation of JSONParser that adds prefixes to test behavior
type mockParser struct {
	name          string
	marshalPrefix string
	shouldFail    bool
}

var _ JSONParser = (*mockParser)(nil)

func (m *mockParser) Name() string {
	return m.name
}

func (m *mockParser) Marshal(v interface{}) ([]byte, error) {
	if m.shouldFail {
		return nil, fmt.Errorf("mock marshal error")
	}

	// Use standard library to marshal
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Add prefix to the JSON string if prefix is set
	if m.marshalPrefix != "" {
		result := fmt.Sprintf(`%s%s`, m.marshalPrefix, string(data))
		return []byte(result), nil
	}

	return data, nil
}

func (m *mockParser) Unmarshal(data []byte, v interface{}) error {
	if m.shouldFail {
		return fmt.Errorf("mock unmarshal error")
	}

	dataStr := string(data)

	// Remove test prefix if present
	if m.marshalPrefix != "" && strings.HasPrefix(dataStr, m.marshalPrefix) {
		dataStr = strings.TrimPrefix(dataStr, m.marshalPrefix)
	}

	return json.Unmarshal([]byte(dataStr), v)
}

// uppercaseParser converts all string values to uppercase
type uppercaseParser struct{}

var _ JSONParser = (*uppercaseParser)(nil)

func (u *uppercaseParser) Name() string {
	return "uppercase-json"
}

func (u *uppercaseParser) Marshal(v interface{}) ([]byte, error) {
	// First marshal normally
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// Convert to uppercase (this will break JSON format, but good for testing)
	return []byte(strings.ToUpper(string(data))), nil
}

func (u *uppercaseParser) Unmarshal(data []byte, v interface{}) error {
	// Convert back to lowercase before unmarshaling
	normalData := strings.ToLower(string(data))
	return json.Unmarshal([]byte(normalData), v)
}

func TestRegisterParser(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	t.Run("successful registration", func(t *testing.T) {
		parser := &mockParser{name: "test-parser"}
		err := RegisterParser(parser)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("nil parser", func(t *testing.T) {
		err := RegisterParser(nil)
		if err == nil {
			t.Error("expected error for nil parser")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		parser := &mockParser{name: ""}
		err := RegisterParser(parser)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		parser1 := &mockParser{name: "duplicate-parser"}
		parser2 := &mockParser{name: "duplicate-parser"}

		err := RegisterParser(parser1)
		if err != nil {
			t.Errorf("expected no error for first registration, got %v", err)
		}
		err = RegisterParser(parser2)
		if err == nil {
			t.Error("expected error for duplicate registration")
		}
	})
}

func TestMustRegisterParser(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	t.Run("successful registration", func(t *testing.T) {
		parser := &mockParser{name: "must-test-parser"}
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		MustRegisterParser(parser)
	})

	t.Run("panic on nil parser", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for nil parser")
			}
		}()
		MustRegisterParser(nil)
	})

	t.Run("panic on duplicate", func(t *testing.T) {
		parser1 := &mockParser{name: "panic-duplicate-parser"}
		parser2 := &mockParser{name: "panic-duplicate-parser"}

		MustRegisterParser(parser1)
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for duplicate parser")
			}
		}()
		MustRegisterParser(parser2)
	})
}

func TestSetDefaultParser(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	t.Run("successful change", func(t *testing.T) {
		parser := &mockParser{name: "new-parser"}
		MustRegisterParser(parser)

		err := SetDefaultParser("new-parser")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if DefaultParser() != "new-parser" {
			t.Errorf("expected default to be %q, got %q", "new-parser", DefaultParser())
		}
	})

	t.Run("unregistered parser", func(t *testing.T) {
		err := SetDefaultParser("non-existent-parser")
		if err == nil {
			t.Error("expected error for unregistered parser")
		}
	})
}

func TestMustSetDefaultParser(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	t.Run("successful change", func(t *testing.T) {
		parser := &mockParser{name: "must-parser"}
		MustRegisterParser(parser)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("unexpected panic: %v", r)
			}
		}()
		MustSetDefaultParser("must-parser")
		if DefaultParser() == "must-parser" {
			t.Errorf("expected default to be %q, got %q", "must-parser", DefaultParser())
		}
	})

	t.Run("panic on unregistered parser", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for unregistered parser")
			}
		}()
		MustSetDefaultParser("non-existent-parser")
	})
}

func TestJitJSONWithCustomParsers(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestData{Name: "test", Value: 42}

	t.Run("using uppercase parser", func(t *testing.T) {
		parser := &uppercaseParser{}
		MustRegisterParser(parser)
		MustSetDefaultParser("uppercase-json")

		jit := New(testData)

		// Test marshal
		data, err := jit.Marshal()
		if err != nil {
			t.Errorf("marshal error: %v", err)
		}

		// Should be uppercase
		expected := `{"NAME":"TEST","VALUE":42}`
		if string(data) != expected {
			t.Errorf("expected %q, got %q", expected, string(data))
		}

		// Test unmarshal
		jit2 := NewFromBytes[TestData](data)
		result, err := jit2.Unmarshal()
		if err != nil {
			t.Errorf("unmarshal error: %v", err)
		}

		if result.Name != testData.Name || result.Value != testData.Value {
			t.Errorf("expected %+v, got %+v", testData, result)
		}
	})

	t.Run("marshal error handling", func(t *testing.T) {
		parser := &mockParser{
			name:       "failing-marshal-parser",
			shouldFail: true,
		}
		MustRegisterParser(parser)
		MustSetDefaultParser("failing-marshal-parser")

		jit := New(testData)
		_, err := jit.Marshal()
		if err == nil {
			t.Error("expected marshal error")
		}
	})

	t.Run("unmarshal error handling", func(t *testing.T) {
		parser := &mockParser{
			name:       "failing-unmarshal-parser",
			shouldFail: true,
		}
		MustRegisterParser(parser)
		MustSetDefaultParser("failing-unmarshal-parser")

		data := []byte(`{"name":"test","value":42}`)
		jit := NewFromBytes[TestData](data)
		_, err := jit.Unmarshal()
		if err == nil {
			t.Error("expected unmarshal error")
		}
	})
}

func TestJitJSONSetParser(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	type TestData struct {
		Message string `json:"message"`
	}

	testData := TestData{Message: "hello"}
	originalDefault := DefaultParser()

	t.Run("successful parser change", func(t *testing.T) {
		parser := &mockParser{
			name:          "instance-parser",
			marshalPrefix: "INSTANCE:",
		}
		MustRegisterParser(parser)

		jit := New(testData)

		// Ceheck default parser
		if jit.Parser() != originalDefault {
			t.Errorf("expected parser %q, got %q", originalDefault, jit.Parser())
		}

		// Change parser to jit
		err := jit.SetParser("instance-parser")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify parser changed
		if jit.Parser() != "instance-parser" {
			t.Errorf("expected parser %q, got %q", "instance-parser", jit.Parser())
		}

		// Test that the new parser
		data, err := jit.Marshal()
		if err != nil {
			t.Errorf("marshal error: %v", err)
		}

		expected := `INSTANCE:{"message":"hello"}`
		if string(data) != expected {
			t.Errorf("expected %q, got %q", expected, string(data))
		}
	})

	t.Run("unregistered parser", func(t *testing.T) {
		jit := New(testData)
		err := jit.SetParser("non-existent-parser")
		if err == nil {
			t.Error("expected error for unregistered parser")
		}
	})
}

func TestParserNilCases(t *testing.T) {
	t.Cleanup(func() {
		resetParserRegistry(t)
	})

	type TestData struct {
		Value string `json:"value"`
	}

	t.Run("nil value marshal", func(t *testing.T) {
		parser := &mockParser{name: "nil-test-parser", marshalPrefix: "NIL:"}
		MustRegisterParser(parser)
		MustSetDefaultParser("nil-test-parser")

		var nilPtr *TestData
		jit := New(nilPtr)

		data, err := jit.Marshal()
		if err != nil {
			t.Errorf("expected no error for nil marshal, got %v", err)
		}

		expected := "NIL:null"
		if string(data) != expected {
			t.Errorf("expected %q, got %q", expected, string(data))
		}
	})

	t.Run("empty data unmarshal", func(t *testing.T) {
		parser := &mockParser{name: "empty-test-parser"}
		MustRegisterParser(parser)
		MustSetDefaultParser("empty-test-parser")

		jit := NewFromBytes[TestData](nil)

		result, err := jit.Unmarshal()
		if err != nil {
			t.Errorf("expected no error for empty data unmarshal, got %v", err)
		}

		expected := TestData{}
		if result != expected {
			t.Errorf("expected %+v, got %+v", expected, result)
		}
	})
}
