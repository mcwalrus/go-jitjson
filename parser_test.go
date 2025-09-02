package jitjson_test

import (
	"sync"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

func TestSetDefaultParser(t *testing.T) {
	// Save original parser
	originalParser := jitjson.GetDefaultParser()
	defer func() {
		// Restore original parser
		if originalParser == "encoding/json" {
			jitjson.MustSetDefaultParser(jitjson.JsonV1)
		} else {
			jitjson.MustSetDefaultParser(jitjson.JsonV2)
		}
	}()

	t.Run("valid parsers", func(t *testing.T) {
		err := jitjson.SetDefaultParser(jitjson.JsonV1)
		if err != nil {
			t.Errorf("SetDefaultParser(JsonV1) failed: %v", err)
		}
		if got := jitjson.GetDefaultParser(); got != "encoding/json" {
			t.Errorf("GetDefaultParser() = %v, want %v", got, "encoding/json")
		}

		err = jitjson.SetDefaultParser(jitjson.JsonV2)
		if err != nil {
			t.Errorf("SetDefaultParser(JsonV2) failed: %v", err)
		}
		if got := jitjson.GetDefaultParser(); got != "encoding/json/v2" {
			t.Errorf("GetDefaultParser() = %v, want %v", got, "encoding/json/v2")
		}
	})

	t.Run("invalid parser", func(t *testing.T) {
		err := jitjson.SetDefaultParser(99) // Invalid jsonType
		if err == nil {
			t.Error("SetDefaultParser with invalid version should return error")
		}
	})
}

func TestMustSetDefaultParser(t *testing.T) {
	// Save original parser
	originalParser := jitjson.GetDefaultParser()
	defer func() {
		// Restore original parser
		if originalParser == "encoding/json" {
			jitjson.MustSetDefaultParser(jitjson.JsonV1)
		} else {
			jitjson.MustSetDefaultParser(jitjson.JsonV2)
		}
	}()

	t.Run("valid parser", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustSetDefaultParser should not panic with valid parser: %v", r)
			}
		}()
		jitjson.MustSetDefaultParser(jitjson.JsonV1)
	})

	t.Run("invalid parser panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSetDefaultParser should panic with invalid parser")
			}
		}()
		jitjson.MustSetDefaultParser(99) // Invalid jsonType
	})
}

func TestJitJSONSetParser(t *testing.T) {
	type TestStruct struct {
		Value string `json:"value"`
	}

	jit := jitjson.New(TestStruct{Value: "test"})

	t.Run("valid parsers", func(t *testing.T) {
		err := jit.SetParser(jitjson.JsonV1)
		if err != nil {
			t.Errorf("SetParser(JsonV1) failed: %v", err)
		}
		if got := jit.Parser(); got != "encoding/json" {
			t.Errorf("Parser() = %v, want %v", got, "encoding/json")
		}

		err = jit.SetParser(jitjson.JsonV2)
		if err != nil {
			t.Errorf("SetParser(JsonV2) failed: %v", err)
		}
		if got := jit.Parser(); got != "encoding/json/v2" {
			t.Errorf("Parser() = %v, want %v", got, "encoding/json/v2")
		}
	})

	t.Run("invalid parser", func(t *testing.T) {
		err := jit.SetParser(99) // Invalid jsonType
		if err == nil {
			t.Error("SetParser with invalid version should return error")
		}
	})
}

func TestJitJSONMustSetParser(t *testing.T) {
	type TestStruct struct {
		Value string `json:"value"`
	}

	jit := jitjson.New(TestStruct{Value: "test"})

	t.Run("valid parser", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustSetParser should not panic with valid parser: %v", r)
			}
		}()
		jit.MustSetParser(jitjson.JsonV1)
	})

	t.Run("invalid parser panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustSetParser should panic with invalid parser")
			}
		}()
		jit.MustSetParser(99) // Invalid jsonType
	})
}

func TestThreadSafety(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numIterations)

	// Test concurrent access to default parser
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Alternate between parsers
				if (id+j)%2 == 0 {
					if err := jitjson.SetDefaultParser(jitjson.JsonV1); err != nil {
						errors <- err
						return
					}
				} else {
					if err := jitjson.SetDefaultParser(jitjson.JsonV2); err != nil {
						errors <- err
						return
					}
				}

				// Read parser name
				_ = jitjson.GetDefaultParser()

				// Create new JitJSON instances
				jit := jitjson.New("test")
				_ = jit.Parser()
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent operation failed: %v", err)
	}
}

func TestParserConsistency(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	testData := TestStruct{Name: "test", Value: 42}

	// Test with JsonV1
	t.Run("JsonV1 consistency", func(t *testing.T) {
		jit := jitjson.New(testData)
		err := jit.SetParser(jitjson.JsonV1)
		if err != nil {
			t.Fatal(err)
		}

		// Marshal and unmarshal
		data, err := jit.Marshal()
		if err != nil {
			t.Fatal(err)
		}

		jit2 := jitjson.NewFromBytes[TestStruct](data)
		err = jit2.SetParser(jitjson.JsonV1)
		if err != nil {
			t.Fatal(err)
		}

		result, err := jit2.Unmarshal()
		if err != nil {
			t.Fatal(err)
		}

		if result.Name != testData.Name || result.Value != testData.Value {
			t.Errorf("Data mismatch: got %+v, want %+v", result, testData)
		}
	})

	// Test with JsonV2
	t.Run("JsonV2 consistency", func(t *testing.T) {
		jit := jitjson.New(testData)
		err := jit.SetParser(jitjson.JsonV2)
		if err != nil {
			t.Fatal(err)
		}

		// Marshal and unmarshal
		data, err := jit.Marshal()
		if err != nil {
			t.Fatal(err)
		}

		jit2 := jitjson.NewFromBytes[TestStruct](data)
		err = jit2.SetParser(jitjson.JsonV2)
		if err != nil {
			t.Fatal(err)
		}

		result, err := jit2.Unmarshal()
		if err != nil {
			t.Fatal(err)
		}

		if result.Name != testData.Name || result.Value != testData.Value {
			t.Errorf("Data mismatch: got %+v, want %+v", result, testData)
		}
	})
}
