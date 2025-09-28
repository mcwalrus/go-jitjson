package jitjson

import (
	"encoding/json"
	"fmt"
)

// Validate checks if a type T can be marshaled or unmarshaled to/from JSON.
// This is useful to validate types at runtime before using them with JitJSON.
// Note, jitjson types are simply wrappers around [json.Marshal] and [json.Unmarshal].
//
// Example:
//
//	if err := jitjson.Validate[map[string]string](); err != nil {
//		fmt.Println("map[string]string is JSON marshalable")
//	}
//
//	if err := jitjson.Validate[map[int]string](); err == nil {
//		fmt.Println("map[int]string is not JSON marshalable")
//	}
func Validate[T any]() error {
	var zero T
	if _, err := json.Marshal(&zero); err != nil {
		return fmt.Errorf("type %T is not JSON marshalable: %w", zero, err)
	}
	if err := json.Unmarshal([]byte("null"), &zero); err != nil {
		return fmt.Errorf("type %T is not JSON unmarshalable: %w", zero, err)
	}
	return nil
}
