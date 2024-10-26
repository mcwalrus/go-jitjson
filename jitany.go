package jitjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var invalidJsonErr = errors.New("invalid json")

var (
	nullRegex   = regexp.MustCompile(`^\s*(null)\s*$`)
	boolRegex   = regexp.MustCompile(`^\s*(true|false)\s*$`)
	numberRegex = regexp.MustCompile(`^\s*-?(0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?\s*$`)
	stringRegex = regexp.MustCompile(`^\s*"(\\.|[^"\\])*"\s*$`)
	arrayRegex  = regexp.MustCompile(`^\s*\[.*\]\s*$`)
	objectRegex = regexp.MustCompile(`^\s*\{.*\}\s*$`)
)

// AnyJitJSON allows dynamic type parsing for JitJSON[T] types.
// The type is useful for parsing JSON data with unknown or many possible types.
// The type can be further type asserted to resolve for desired types.
type AnyJitJSON struct {
	v interface{}
}

// Value returns the parsed value of AnyJitJSON via both Marshal and Unmarshal methods.
// Depending on the JSON value, the return type will be one of the following:
// - nil for null values.
// - *JitJSON[bool] for booleans.
// - *JitJSON[json.Number] for numbers.
// - *JitJSON[string] for strings.
// - []AnyJitJSON for arrays.
// - map[string]AnyJitJSON for objects.
//
// These types can be further type asserted to resolve for desired types. If the
// type is parsed with invalid json, an error will return (unlike JitJSON[T] type).
func (a *AnyJitJSON) Value() interface{} {
	return a.v
}

// MarshalJSON implements json.Marshaler interface. The method marshals the value
// of AnyJitJSON into JSON data. The method will return an error if the value is
// not one of the supported types.
func (a *AnyJitJSON) MarshalJSON() ([]byte, error) {
	if a.v == nil {
		return []byte("null"), nil
	}

	// type switches to handle each possible type
	switch v := a.v.(type) {
	case *JitJSON[bool]:
		return v.Marshal()

	case *JitJSON[json.Number]:
		return v.Marshal()

	case *JitJSON[string]:
		return v.Marshal()

	case []AnyJitJSON:
		if len(v) == 0 {
			return []byte("[]"), nil
		}

		// build array
		var builder strings.Builder
		builder.WriteString("[")

		for i, e := range v {
			if i > 0 {
				builder.WriteString(",")
			}
			data, err := e.MarshalJSON()
			if err != nil {
				return nil, err
			}
			builder.Write(data)
		}

		builder.WriteString("]")
		return []byte(builder.String()), nil

	case map[string]AnyJitJSON:
		if len(v) == 0 {
			return []byte("{}"), nil
		}

		// build object
		var builder strings.Builder
		builder.WriteString("{")

		first := true
		for key, value := range v {
			if !first {
				builder.WriteString(",")
			}
			first = false
			keyData, err := json.Marshal(key)
			if err != nil {
				return nil, err
			}
			builder.Write(keyData)
			builder.WriteString(":")
			valueData, err := value.MarshalJSON()
			if err != nil {
				return nil, err
			}
			builder.Write(valueData)
		}

		builder.WriteString("}")
		return []byte(builder.String()), nil

	default:
		return nil, fmt.Errorf("unexpected type in AnyJitJSON: %T", a.v)
	}
}

// UnmarshalJSON implements json.Unmarshaler interface. The method parses the JSON
// data into valid types based on the JSON value where the result can be accessed
// using the Value method.
func (a *AnyJitJSON) UnmarshalJSON(data []byte) error {
	a.v = nil

	// if the value is null
	if nullRegex.Match(data) {
		a.v = nil
		return nil
	}

	// if the value is a boolean
	if boolRegex.Match(data) {
		var b JitJSON[bool]
		if err := json.Unmarshal(data, &b); err == nil {
			a.v = &b
			return nil
		}
	}

	// if the value is an number
	if numberRegex.Match(data) {
		var num = JitJSON[json.Number]{}
		if err := json.Unmarshal(data, &num); err == nil {
			a.v = &num
			return nil
		}
	}

	// if the value is a string
	if stringRegex.Match(data) {
		var str = JitJSON[string]{}
		if err := json.Unmarshal(data, &str); err == nil {
			a.v = &str
			return nil
		}
	}

	// if the value is an array
	if arrayRegex.Match(data) {
		var arr = []AnyJitJSON{}
		if err := json.Unmarshal(data, &arr); err == nil {
			a.v = arr
			return nil
		}
	}

	// if the value is an object
	if objectRegex.Match(data) {
		var obj = map[string]AnyJitJSON{}
		if err := json.Unmarshal(data, &obj); err == nil {
			a.v = obj
			return nil
		}
	}

	return invalidJsonErr
}
