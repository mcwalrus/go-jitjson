package jitjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	nullRegex   = regexp.MustCompile(`^\s*null\s*$`)
	arrayRegex  = regexp.MustCompile(`^\s*\[\s*(.|\n)*\]\s*$`)
	objectRegex = regexp.MustCompile(`^\s*\{\s*(.|\n)*\}\s*$`)
	boolRegex   = regexp.MustCompile(`^\s*(true|false)\s*$`)
	numberRegex = regexp.MustCompile(`^\s*-?(0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?\s*$`)
	stringRegex = regexp.MustCompile(`^\s*"(\\.|[^"\\])*"\s*$`)
)

// ValueType represents the JSON type of the value stored in AnyJitJSON.
// This type is used to determine the type of the value without performing the
// unmarshalling operation.
type ValueType int

const (
	TypeNull ValueType = iota
	TypeBool
	TypeNumber
	TypeString
	TypeArray
	TypeObject
	TypeInvalid
)

func (v ValueType) String() string {
	return []string{
		"TypeNull",
		"TypeBool",
		"TypeNumber",
		"TypeString",
		"TypeArray",
		"TypeObject",
		"TypeInvalid",
	}[v]
}

// AnyJitJSON can unmarshal arbitrary JSON encodings with just-in-time parsing.
// The value of the encoding can be accessed programmatically by using the methods
// AsBool, AsNumber, AsString, AsArray, and AsObject. // It stores the raw data and
// only unmarshals them when specific values are requested, making it memory efficient
// for large JSON structures where only parts need to be accessed.
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
	val  interface{}
	data []byte
}

// NewAny creates a new AnyJitJSON from JSON data.
func NewAny(data []byte) (*AnyJitJSON, error) {
	var a = &AnyJitJSON{}
	err := a.UnmarshalJSON(data)
	return a, err
}

func (a *AnyJitJSON) String() string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, a.data, "", "  ")
	if err != nil {
		return string(a.data)
	}
	return prettyJSON.String()
}

// Type returns the ValueType of the current AnyJitJSON value. This method can be used
// to determine the type of the underlying value alternatively to type assertion. These
// cover all possible types that can be returned from the Value method:
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
//	case TypeInvalid:
//		panic("invalid type")
//	}
func (a *AnyJitJSON) Type() ValueType {
	if a == nil || a.val == nil {
		return TypeNull
	}
	switch a.val.(type) {
	case *JitJSON[bool]:
		return TypeBool
	case *JitJSON[json.Number]:
		return TypeNumber
	case *JitJSON[string]:
		return TypeString
	case []*AnyJitJSON:
		return TypeArray
	case map[string]*AnyJitJSON:
		return TypeObject
	default:
		return TypeInvalid
	}
}

// IsNull returns true if the AnyJitJSON value is nil.
func (a *AnyJitJSON) IsNull() bool {
	return a.Type() == TypeNull
}

// AsBool returns a bool from AnyJitJSON if possible.
// This method will return false if the value is not a boolean.
func (a *AnyJitJSON) AsBool() (bool, bool) {
	jit, ok := (a.val).(*JitJSON[bool])
	if !ok {
		return false, false
	}
	val, _ := jit.Unmarshal()
	return val, true
}

// AsNumber returns a json.Number from AnyJitJSON if possible.
// This method will return false if the value is not a number.
func (a *AnyJitJSON) AsNumber() (json.Number, bool) {
	jit, ok := (a.val).(*JitJSON[json.Number])
	if !ok {
		return "", false
	}
	val, _ := jit.Unmarshal()
	return val, true
}

// AsString returns a string from AnyJitJSON if possible.
// This method will return false if the value is not a string.
func (a *AnyJitJSON) AsString() (string, bool) {
	jit, ok := (a.val).(*JitJSON[string])
	if !ok {
		return "", false
	}
	val, _ := jit.Unmarshal()
	return val, true
}

// AsArray returns a []*AnyJitJSON from AnyJitJSON if possible.
// This method will return false if the value is not an array.
func (a *AnyJitJSON) AsArray() ([]*AnyJitJSON, bool) {
	if a.data == nil {
		return nil, false
	}
	if _, ok := a.val.([]*AnyJitJSON); !ok {
		return nil, false
	}

	var arr []*AnyJitJSON
	if err := json.Unmarshal(a.data, &arr); err != nil {
		return nil, false
	}

	a.data = nil
	return arr, true
}

// AsObject returns a map[string]*AnyJitJSON from AnyJitJSON if possible.
// This method will return false if the value is not an object.
func (a *AnyJitJSON) AsObject() (map[string]*AnyJitJSON, bool) {
	if a.data == nil {
		return nil, false
	}
	if _, ok := a.val.(map[string]*AnyJitJSON); !ok {
		return nil, false
	}

	var obj map[string]*AnyJitJSON
	if err := json.Unmarshal(a.data, &obj); err != nil {
		return nil, false
	}

	a.data = nil
	return obj, true
}

// MarshalJSON returns the JSON encoding of the value.
func (a *AnyJitJSON) MarshalJSON() ([]byte, error) {
	return a.data, nil
}

// UnmarshalJSON parses the JSON data and stores the value in AnyJitJSON. The method
// supports all valid JSON value types (null, boolean, number, string, array, object).
func (a *AnyJitJSON) UnmarshalJSON(data []byte) error {
	a.val = nil
	a.data = data
	var err error

	// if the value is null
	if nullRegex.Match(data) {
		a.val = nil
		return nil
	}

	// if the value is a boolean
	if boolRegex.Match(data) {
		var b JitJSON[bool]
		if err = json.Unmarshal(data, &b); err == nil {
			a.val = &b
			return nil
		}
	}

	// if the value is an number
	if numberRegex.Match(data) {
		var num JitJSON[json.Number]
		if err = json.Unmarshal(data, &num); err == nil {
			a.val = &num
			return nil
		}
	}

	// if the value is a string
	if stringRegex.Match(data) {
		var str JitJSON[string]
		if err = json.Unmarshal(data, &str); err == nil {
			a.val = &str
			return nil
		}
	}

	// if the value is an array
	if arrayRegex.Match(data) {
		a.val = []*AnyJitJSON{}
		a.data = make([]byte, len(data))
		copy(a.data, data)
		return nil
	}

	// if the value is an object
	if objectRegex.Match(data) {
		a.val = map[string]*AnyJitJSON{}
		a.data = make([]byte, len(data))
		copy(a.data, data)
		return nil
	}

	return fmt.Errorf("invalid json: %w", err)
}
