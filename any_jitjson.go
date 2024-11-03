package jitjson

import (
	"encoding/json"
	"fmt"
	"regexp"
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
//	err := json.Unmarshal([]byte(`{"key": [1, "text", true]}`), &any)
//	if err != nil {
//		panic(err)
//	}
//
//	// Get parsed value
//	value := any.Value()
//	fmt.Println(value) // Output: map[key:[1 text true]]
//
//	// Access the object value
//	m, ok := value.(map[string]*jitjson.AnyJitJSON)
//	if !ok {
//		panic("not a map")
//	}
//
//	// Access the array value
//	sl, ok := m["key"].Value().([]*jitjson.AnyJitJSON)
//	if !ok {
//		panic("not an array")
//	}
//
//	// Access number value
//	num := sl[0].Value().(*jitjson.JitJSON[json.Number])
//	n, err := num.Unmarshal()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(n.Int64()) // Output: 1 <nil>
//
//	// Access string value
//	str := sl[1].Value().(*jitjson.JitJSON[string])
//	s, err := str.Unmarshal()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(s) // Output: text
//
//	// Access boolean value
//	boo := sl[2].Value().(*jitjson.JitJSON[bool])
//	b, err := boo.Unmarshal()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(b) // Output: true
type AnyJitJSON struct {
	v interface{}
}

// Value returns the underlying value of AnyJitJSON which is parsed via the UnmarshalJSON
// method. The value can be one of the following:
//
//	var any AnyJitJSON
//	value := any.Value()
//
//	switch any.Type() {
//	case TypeNull:
//	    // handle nil value
//	case TypeBool:
//		var b = value.(*JitJSON[bool])
//	case TypeNumber:
//		var num = value.(*JitJSON[json.Number])
//	case TypeString:
//	    var str = value.(*JitJSON[string])
//	case TypeArray:
//	    var arr = value.([]*AnyJitJSON)
//	case TypeObject:
//	    var obj = value.(map[string]*AnyJitJSON)
//	}
//
// These cover all possible types that can be returned from the Value method. Alternatively,
// the Type method can be used to determine the type of the underlying value for type assertion.
func (a *AnyJitJSON) Value() interface{} {
	return a.v
}

// Type returns the ValueType of the current AnyJitJSON value. This method can be used
// to determine the type of the underlying value alternatively to type assertion.
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
