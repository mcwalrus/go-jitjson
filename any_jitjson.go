package jitjson

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	wrongTypeErr = fmt.Errorf("wrong type in AnyJitJSON")
)

var (
	nullRegex   = regexp.MustCompile(`^\s*(null)\s*$`)
	boolRegex   = regexp.MustCompile(`^\s*(true|false)\s*$`)
	numberRegex = regexp.MustCompile(`^\s*-?(0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?\s*$`)
	stringRegex = regexp.MustCompile(`^\s*"(\\.|[^"\\])*"\s*$`)
	arrayRegex  = regexp.MustCompile(`^\s*\[.*\]\s*$`)
	objectRegex = regexp.MustCompile(`^\s*\{.*\}\s*$`)
)

// ValueType represents the type of JSON value stored in AnyJitJSON.
type ValueType int

const (
	TypeNull ValueType = iota
	TypeBool
	TypeNumber
	TypeString
	TypeArray
	TypeObject
)

// AnyJitJSON provides a type for handling arbitrary JSON values with just-in-time parsing.
// It can unmarshal and store any valid JSON value type (null, boolean, number, string, array,
// or object) and defers parsing of JSON values until needed. The result of unmarshaling can
// be accessed via the Value method.
//
// Example:
//
//	var any jitjson.AnyJitJSON
//	err := json.Unmarshal([]byte(`{"key": [1, "text", true, null]}`), &any)
//	if err != nil {
//		panic(err)
//	}
//
//	// Access object
//	m, ok := any.AsObject()
//	if !ok {
//		panic("not a map")
//	}
//
//	// Access array value
//	sl, ok := m["key"].AsArray()
//	if !ok {
//		panic("not an array")
//	}
//
//	// Access number value
//	n, ok := sl[0].AsNumber()
//	if !ok {
//		panic("not a number")
//	}
//	fmt.Println(n.Int64()) // Output: 1 <nil>
//
//	// Access string value
//	str, ok := sl[1].AsString()
//	if !ok {
//		panic("not a string")
//	}
//	fmt.Println(s) // Output: text
//
//	// Access boolean value
//	b, ok := sl[2].AsBool()
//	if !ok {
//		panic("not a boolean")
//	}
//	fmt.Println(b) // Output: true
//
//	// Access null value
//	fmt.Println(sl[3].IsNull()) // Output: true
type AnyJitJSON struct {
	v interface{}
}

// NewAnyJitJSON creates a new AnyJitJSON value from the given JSON data. The method
// unmarshals the JSON data and stores the value in AnyJitJSON. If the JSON data is
// invalid, the method returns an error.
func NewAnyJitJSON(data []byte) *AnyJitJSON {
	var a = &AnyJitJSON{}
	_ = a.UnmarshalJSON(data)
	return a
}

// Type returns the ValueType of the current AnyJitJSON value. This method can be used
// to determine the type of the underlying value alternatively to type assertion.
//
//	var any AnyJitJSON
//	switch any.Type() {
//	case TypeNull:
//	    // handle nil value
//	case TypeBool:
//		b, _ := any.AsBool()
//		// handle boolean value
//	case TypeNumber:
//		num, _ := any.AsNumber()
//		// handle number value
//	case TypeString:
//	    str, _ := any.AsString()
//		// handle string value
//	case TypeArray:
//	    arr, _ := any.AsArray()
//		// handle array value
//	case TypeObject:
//	    obj, _ := any.AsObject()
//		// handle object value
//	}
//
// These cover all possible types that can be returned from the Value method. Alternatively,
// the Type method can be used to determine the type of the underlying value for type assertion.
func (a *AnyJitJSON) Type() ValueType {
	switch a.v.(type) {
	case nil:
		return TypeNull
	case *JitJSON[bool]:
		return TypeBool
	case *JitJSON[json.Number]:
		return TypeNumber
	case *JitJSON[string]:
		return TypeString
	case []AnyJitJSON:
		return TypeArray
	case map[string]AnyJitJSON:
		return TypeObject
	default:
		return TypeNull
	}
}

// IsNull returns true if the AnyJitJSON is a null value.
func (a *AnyJitJSON) IsNull() bool {
	return a.Type() == TypeNull
}

// AsBool returns value of AnyJitJSON as a bool if possible.
func (a *AnyJitJSON) AsBool() (bool, bool) {
	jit, ok := (a.v).(*JitJSON[bool])
	if !ok {
		return false, false
	}
	val, _ := jit.Decode()
	return val, true
}

// AsNumber returns value of AnyJitJSON as a json.Number if possible.
func (a *AnyJitJSON) AsNumber() (json.Number, bool) {
	jit, ok := (a.v).(*JitJSON[json.Number])
	if !ok {
		return "", false
	}
	val, _ := jit.Decode()
	return val, true
}

// AsString returns value of AnyJitJSON as a string if possible.
func (a *AnyJitJSON) AsString() (string, bool) {
	jit, ok := (a.v).(*JitJSON[string])
	if !ok {
		return "", false
	}
	val, _ := jit.Decode()
	return val, true
}

// AsArray returns value of AnyJitJSON as []*AnyJitJSON if possible.
func (a *AnyJitJSON) AsArray() ([]*AnyJitJSON, bool) {
	arr, ok := a.v.([]*AnyJitJSON)
	if !ok {
		return nil, false
	}
	return arr, true
}

// AsObject returns value of AnyJitJSON as map[string]*AnyJitJSON if possible.
func (a *AnyJitJSON) AsObject() (map[string]*AnyJitJSON, bool) {
	obj, ok := a.v.(map[string]*AnyJitJSON)
	if !ok {
		return nil, false
	}
	return obj, true
}

// MarshalJSON parses the value to return the JSON encoding of the value. The method
// will either achieves this through JitJSON cache or by marshaling the value.
func (a *AnyJitJSON) MarshalJSON() ([]byte, error) {
	if a == nil || a.v == nil {
		return []byte("null"), nil
	}

	// type switches to handle each possible type
	switch v := a.v.(type) {
	case *JitJSON[bool]:
		return v.Encode()

	case *JitJSON[json.Number]:
		return v.Encode()

	case *JitJSON[string]:
		return v.Encode()

	case []*AnyJitJSON:
		if len(v) == 0 {
			return []byte("[]"), nil
		}

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

	case map[string]*AnyJitJSON:
		if len(v) == 0 {
			return []byte("{}"), nil
		}

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

// UnmarshalJSON parses the JSON data and stores the value in AnyJitJSON. The method
// supports all valid JSON value types (null, boolean, number, string, array, object).
// If the json is invalid, the method returns an error.
func (a *AnyJitJSON) UnmarshalJSON(data []byte) error {
	a.v = nil
	var err error

	// if the value is null
	if nullRegex.Match(data) {
		a.v = nil
		return nil
	}

	// if the value is a boolean
	if boolRegex.Match(data) {
		var b JitJSON[bool]
		if err = json.Unmarshal(data, &b); err == nil {
			a.v = &b
			return nil
		}
	}

	// if the value is an number
	if numberRegex.Match(data) {
		var num = JitJSON[json.Number]{}
		if err = json.Unmarshal(data, &num); err == nil {
			a.v = &num
			return nil
		}
	}

	// if the value is a string
	if stringRegex.Match(data) {
		var str = JitJSON[string]{}
		if err = json.Unmarshal(data, &str); err == nil {
			a.v = &str
			return nil
		}
	}

	// if the value is an array
	if arrayRegex.Match(data) {
		var arr = []*AnyJitJSON{}
		if err = json.Unmarshal(data, &arr); err == nil {
			a.v = arr
			return nil
		}
	}

	// if the value is an object
	if objectRegex.Match(data) {
		var obj = map[string]*AnyJitJSON{}
		if err = json.Unmarshal(data, &obj); err == nil {
			a.v = obj
			return nil
		}
	}

	return fmt.Errorf("invalid json: %w", err)
}
