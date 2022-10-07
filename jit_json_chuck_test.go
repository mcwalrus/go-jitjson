package jitjson_test

import (
	"encoding/json"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

// see tests for jit-json parsing for arrays / objects.
var (
	arrayChuckData = []byte(`{
		"array": [
			{
				"Name": "Harry Potter",
				"Age": 12
			},
			{
				"Name": "Hermione Granger",
				"Age": 13
			},
			{
				"Name": "Ron Weasley",
				"Age": 11
			}
		]
	}`)
	objectChuckData = []byte(`{
		"object": {
			"1": {
				"Name": "Harry Potter",
				"Age": 12
			},
			"2": {
				"Name": "Hermione Granger",
				"Age": 13
			},
			"3": {
				"Name": "Ron Weasley",
				"Age": 11
			}
		}
	}`)
)

func TestArrayChuckingJitJSON(t *testing.T) {
	var data = arrayChuckData

	type myObject struct {
		MyArray []*jitjson.JitJSON[struct {
			Age  int
			Name string
		}] `json:"array"`
	}

	var jit myObject
	err := json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}

	// decode values of array just-in-time.
	if len(jit.MyArray) != 3 {
		t.Error("array not populated")
	}

	for i := range jit.MyArray {
		val := jit.MyArray[i]

		_, err = val.Marshal()
		if err != nil {
			t.Error(err)
		}

		_, err = val.Unmarshal()
		if err != nil {
			t.Error(err)
		}
	}

	err = json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}
}

func TestObjectChuckingJitJSON(t *testing.T) {
	var data = objectChuckData

	type myObject struct {
		MyObject map[string]*jitjson.JitJSON[struct {
			Age  int
			Name string
		}] `json:"object"`
	}

	var jit myObject
	err := json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}

	// decode values of map just-in-time.
	if len(jit.MyObject) != 3 {
		t.Error("object not populated")
	}

	for key := range jit.MyObject {
		val := jit.MyObject[key]

		_, err = val.Marshal()
		if err != nil {
			t.Error(err)
		}

		_, err = val.Unmarshal()
		if err != nil {
			t.Error(err)
		}
	}

	err = json.Unmarshal(data, &jit)
	if err != nil {
		t.Error(err)
	}
}
