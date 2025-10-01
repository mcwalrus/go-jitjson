package jitjson

import (
	"encoding/json"
	"fmt"
)

// Validate checks if a type T can be marshaled or unmarshaled to/from JSON.
// This is useful to validate types at runtime before using them with JitJSON.
// Note, jitjson are simply wrappers around [json.Marshal] and [json.Unmarshal].
//
// Example:
//
//	// Key is not comparable
//	if err := jitjson.Validate[map[int]string](); err != nil {
//		fmt.Println("map[int]string is not JSON marshalable")
//	}
func Validate[T any]() error {
	var zero T
	_, err1 := json.Marshal(&zero)
	err2 := json.Unmarshal([]byte("null"), &zero)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("type %T is not JSON parseable", zero)
	}
	return nil
}
