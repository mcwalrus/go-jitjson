package jitjson_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

func TestJitJSONInterface(t *testing.T) {
	var _ jitjson.JitJSONInterface = &jitjson.JitJSON[int]{}
}

func TestJitJSONInterfaceSwitch(t *testing.T) {
	var i jitjson.JitJSONInterface = &jitjson.JitJSON[int]{}
	switch i.(type) {
	case *jitjson.JitJSON[int]:
		break
	default:
		t.Error("expected type above")
	}
}

func TestJitJSONInterfaceMultipleTypes(t *testing.T) {
	var (
		err error
		iS  = make([]jitjson.JitJSONInterface, 3)
	)

	iS[0], err = jitjson.NewJitJSON[int](1)
	if err != nil {
		t.Error(err)
	}

	iS[1], err = jitjson.NewJitJSON[float64](2.0)
	if err != nil {
		t.Error(err)
	}

	iS[2], err = jitjson.NewJitJSON[string]("this works!")
	if err != nil {
		t.Error(err)
	}
}

func TestNilTypeUnmarshal(t *testing.T) {
	var expected = struct{}{}

	jit, err := jitjson.NewJitJSON[struct{}](nil)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestNilValueMarshal(t *testing.T) {
	var expected = []byte{}

	jit, err := jitjson.NewJitJSON[struct{}](nil)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestEmptyStructTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`null`)
		expected = struct{}{}
	)

	jit, err := jitjson.NewJitJSON[struct{}](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestEmptyStructValueMarshal(t *testing.T) {
	var (
		value    = struct{}{}
		expected = []byte(`{}`)
	)

	jit, err := jitjson.NewJitJSON[struct{}](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestIntTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`1`)
		expected = 1
	)

	jit, err := jitjson.NewJitJSON[int](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestIntValueMarshal(t *testing.T) {
	var (
		value    = 1
		expected = []byte(`1`)
	)

	jit, err := jitjson.NewJitJSON[int](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestFloatTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`1.0125`)
		expected = 1.0125
	)

	jit, err := jitjson.NewJitJSON[float64](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestFloatValueMarshal(t *testing.T) {
	var (
		value    = 1.0125
		expected = []byte(`1.0125`)
	)

	jit, err := jitjson.NewJitJSON[float64](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestStringTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`"json encoded"`)
		expected = "json encoded"
	)

	jit, err := jitjson.NewJitJSON[string](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestStringValueMarshal(t *testing.T) {
	var (
		value    = "json encoded"
		expected = []byte(`"json encoded"`)
	)

	jit, err := jitjson.NewJitJSON[string](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestPtrJitJSON(t *testing.T) {
	data := []byte(`1`)

	jit, err := jitjson.NewJitJSON[*int](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if *value != 1 {
		t.Error("value not equal to expected")
	}
}

func TestPointerTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`1`)
		expected = 1
	)

	jit, err := jitjson.NewJitJSON[*int](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if *value != expected {
		t.Error("value not equal to expected")
	}
}

func TestPointerValueMarshal(t *testing.T) {
	var (
		value    = 1
		expected = []byte(`1`)
	)

	jit, err := jitjson.NewJitJSON[*int](&value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestArrayTypeUnmarshal(t *testing.T) {
	var (
		data     = []byte(`[1,2,3,4,5]`)
		expected = []int{1, 2, 3, 4, 5}
	)

	jit, err := jitjson.NewJitJSON[[]int](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if !cmpArrExact(value, expected) {
		t.Error("value not equal to expected")
	}
}

func TestArrayValueMarshal(t *testing.T) {
	var (
		value    = []int{1, 2, 3, 4, 5}
		expected = []byte(`[1,2,3,4,5]`)
	)

	jit, err := jitjson.NewJitJSON[[]int](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestPtrArrayTypeUnmarshal(t *testing.T) {
	var (
		_1, _4, _5 = 1, 4, 5
		data       = []byte(`[1,null,null,4,5]`)
		expected   = []*int{&_1, nil, nil, &_4, &_5}
	)

	jit, err := jitjson.NewJitJSON[[]*int](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if !cmpArrPointers(value, expected) {
		t.Error("value not equal to expected")
	}
}

func TestPtrArrayValueMarshal(t *testing.T) {
	var (
		_1, _4, _5 = 1, 4, 5
		value      = []*int{&_1, nil, nil, &_4, &_5}
		expected   = []byte(`[1,null,null,4,5]`)
	)

	jit, err := jitjson.NewJitJSON[[]*int](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestMapTypeUnmarshal(t *testing.T) {
	data := []byte(`{
		"try": "me",
		"knock": "knock",
		"who there": null
	}`)

	jit, err := jitjson.NewJitJSON[map[string]interface{}](data)
	if err != nil {
		t.Error(err)
	}

	m, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if len(m) != 3 {
		t.Error("expected map to be set")
	}
}

func TestMapValueMarshal(t *testing.T) {
	value := map[string]interface{}{
		"try":       "me",
		"knock":     "knock",
		"who there": nil,
	}

	jit, err := jitjson.NewJitJSON[map[string]interface{}](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if len(data) == 0 {
		t.Error("expected data from map")
	}
}

type TestType struct {
	Field *string `json:"field"`
}

func TestStructFieldTypeNilValueUnmarshal(t *testing.T) {
	var (
		data     = []byte(`{"field":null}`)
		expected TestType
	)

	jit, err := jitjson.NewJitJSON[TestType](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if value != expected {
		t.Error("value not equal to expected")
	}
}

func TestStructFieldTypeNilValueMarshal(t *testing.T) {
	var (
		value    TestType
		expected = []byte(`{"field":null}`)
	)

	jit, err := jitjson.NewJitJSON[TestType](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestStructFieldTypeValueUnmarshal(t *testing.T) {
	var (
		str      = "goal"
		data     = []byte(`{"field":"goal"}`)
		expected = TestType{&str}
	)

	jit, err := jitjson.NewJitJSON[TestType](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if *value.Field != *expected.Field {
		t.Error("value not equal to expected")
	}
}

func TestStructFieldTypeValueMarshal(t *testing.T) {
	var (
		str      = "goal"
		value    = TestType{&str}
		expected = []byte(`{"field":"goal"}`)
	)

	jit, err := jitjson.NewJitJSON[TestType](value)
	if err != nil {
		t.Error(err)
	}

	data, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

// note: 'Unmarshal' return value can be updated outside of struct.
func TestUpdatedPointedValueUnmarshal(t *testing.T) {
	type myType struct {
		Value *int
	}

	data := []byte(`{
        "Value": 42
    }`)

	jit, err := jitjson.NewJitJSON[myType](data)
	if err != nil {
		t.Error(err)
	}

	value, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	*value.Value = 12
	value, err = jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if *value.Value != 12 {
		t.Error("expected to receive updated value")
	}
}

// note: 'Marshal' return value cannot be updated outside of struct.
func TestUpdatedPointedValueMarshal(t *testing.T) {
	type myType struct {
		Value *int
	}

	var (
		number = 42
		value  = myType{
			Value: &number,
		}
	)

	data := []byte(`{
        "Value": 42
    }`)

	jit, err := jitjson.NewJitJSON[myType](data)
	if err != nil {
		t.Error(err)
	}

	data1, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	*value.Value = 12
	data2, err := jit.Marshal()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data1, data2) {
		t.Error("expected to receive updated value")
	}
}

func TestNilInterfacesJitJSON(t *testing.T) {
	var err error

	type testInterface interface {
		method() int
	}

	_, err = jitjson.NewJitJSON[interface{}](nil)
	if err != nil {
		t.Error(err)
	}

	_, err = jitjson.NewJitJSON[*interface{}](nil)
	if err != nil {
		t.Error(err)
	}

	_, err = jitjson.NewJitJSON[**interface{}](nil)
	if err != nil {
		t.Error(err)
	}

	_, err = jitjson.NewJitJSON[testInterface](nil)
	if err != nil {
		t.Error(err)
	}

	_, err = jitjson.NewJitJSON[*testInterface](nil)
	if err != nil {
		t.Error(err)
	}

	_, err = jitjson.NewJitJSON[***interface{}](nil)
	if err != nil {
		t.Error(err)
	}
}

func TestInterfacesValueNilUnmarshal(t *testing.T) {
	var data = []byte(`null`)

	jit, err := jitjson.NewJitJSON[interface{}](data)
	if err != nil {
		t.Error(err)
	}

	x, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	if x != nil {
		t.Error("value not nil")
	}
}

func TestInterfacesValueNumberUnmarshal(t *testing.T) {
	var (
		data     = []byte(`1`)
		expected = float64(1)
	)

	jit, err := jitjson.NewJitJSON[interface{}](data)
	if err != nil {
		t.Error(err)
	}

	x, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	v, ok := x.(float64)
	if !ok {
		t.Error("expected float64")
	}

	if v != expected {
		t.Error("value not equal to expected")
	}
}

func TestInterfacesValueStringUnmarshal(t *testing.T) {
	var (
		data     = []byte(`"string"`)
		expected = "string"
	)

	jit, err := jitjson.NewJitJSON[interface{}](data)
	if err != nil {
		t.Error(err)
	}

	x, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	v, ok := x.(string)
	if !ok {
		t.Error("expected string")
	}

	if v != expected {
		t.Error("value not equal to expected")
	}
}

func TestInterfacesEmptyObjectUnmarshal(t *testing.T) {
	var data = []byte(`{}`)

	jit, err := jitjson.NewJitJSON[interface{}](data)
	if err != nil {
		t.Error(err)
	}

	x, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	_, ok := x.(map[string]interface{})
	if !ok {
		t.Error("expected map")
	}
}

func TestInterfacesObjectUnmarshal(t *testing.T) {
	var data = []byte(`{
		"Name": "No-name"
	}`)

	jit, err := jitjson.NewJitJSON[interface{}](data)
	if err != nil {
		t.Error(err)
	}

	x, err := jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}

	m, ok := x.(map[string]interface{})
	if !ok {
		t.Error("expected map")
	}

	v, ok := m["Name"]
	if !ok {
		t.Error("expected key 'Name'")
	}

	if v != "No-name" {
		t.Error("value not equal to expected")
	}
}

func TestArrayUnmarshalJSON(t *testing.T) {
	var (
		data     = []byte(`[1,2,3,4,5]`)
		expected = []int{1, 2, 3, 4, 5}
	)

	var jit []jitjson.JitJSON[int]
	err := json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}

	var arr = make([]int, len(jit))
	for i, js := range jit {
		arr[i], err = js.Unmarshal()
		if err != nil {
			t.Error(err)
		}
	}

	if !cmpArrExact[int](arr, expected) {
		t.Error("values not equal to expected")
	}
}

func TestArrayMarshalJSON(t *testing.T) {
	var (
		err      error
		expected = []byte(`[1,2,3,4,5]`)
	)

	var jit = make([]*jitjson.JitJSON[int], 5)
	for i := range jit {
		jit[i], err = jitjson.NewJitJSON[int](i + 1)
		if err != nil {
			t.Error(err)
		}
	}

	data, err := json.Marshal(&jit)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(data, expected) {
		t.Error("encoding not equal to expected")
	}
}

func TestMultipleUnmarshalJSON(t *testing.T) {
	var (
		v    int
		data []byte
		jit  jitjson.JitJSON[int]
		err  error
	)

	data = []byte(`1`)
	err = json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}
	v, err = jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if v != 1 {
		t.Error("not expected value")
	}

	data = []byte(`2`)
	err = json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}
	v, err = jit.Unmarshal()
	if err != nil {
		t.Error(err)
	}
	if v != 2 {
		t.Error("not expected value")
	}
}

func cmpArrExact[T comparable](as, bs []T) bool {
	if len(as) != len(bs) {
		return false
	}
	for i := range as {
		if as[i] != bs[i] {
			return false
		}
	}
	return true
}

func cmpArrPointers[T comparable](as, bs []*T) bool {
	if len(as) != len(bs) {
		return false
	}
	for i := range as {
		var aPtr, bPtr = as[i], bs[i]
		if (aPtr == nil) != (bPtr == nil) {
			return false
		} else if aPtr != nil && bPtr != nil {
			return *aPtr == *bPtr
		}
	}
	return true
}
