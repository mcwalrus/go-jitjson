package jitjson_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

type Object struct {
	Nil    interface{}   `json:"nil"`
	Bool   bool          `json:"bool"`
	Number float64       `json:"number"`
	String string        `json:"string"`
	Slice  []interface{} `json:"slice"`
	Object *Object       `json:"object"` // recursive
}

const objectTemplate = `{
	"nil": null,
	"bool": true,
	"number": 123.45,
	"string": "Hello, World!",
	"slice": [1, "two", false, null, %d],
	"object": {
		"nil": null,
		"bool": false,
		"number": 456.78,
		"string": "Nested",
		"slice": [4, "five", true, null],
		"object": {
			"nil": null,
			"bool": true,
			"number": 123.45,
			"string": "Hello, World!",
			"slice": [1, "two", false, null, 1],
			"object": {
				"nil": null,
				"bool": false,
				"number": 456.78,
				"string": "Nested",
				"slice": [4, "five", true, null],
				"object": null
			}
		}
	}
}`

func generateObjects(count int) []byte {
	objects := make([]Object, count)
	for i := 0; i < count; i++ {
		err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
		if err != nil {
			panic(err)
		}
	}
	data, err := json.Marshal(objects)
	if err != nil {
		panic(err)
	}
	return data
}

var (
	smallData  = generateObjects(10)
	mediumData = generateObjects(1000)
	largeData  = generateObjects(100000)
)

func shouldParseIterator(parsePercent float64) func() bool {
	var record float64 = 0
	return func() bool {
		record += parsePercent
		if record >= 1 {
			record = record - 1
			return true
		} else {
			return false
		}
	}
}

// BenchmarkParsePercent benchmarks the parsing of JSON data with a given percentage of objects
// that are parsed. It compares the performance of JitJSON and the standard library.
func BenchmarkParsePercent(b *testing.B) {
	parsePercent, err := strconv.ParseFloat(os.Getenv("JITJSON_PARSE_PERCENT"), 64)
	if err != nil {
		b.Fatal(err)
	}

	if parsePercent < 0 || parsePercent > 1 {
		b.Fatal("JITJSON_PARSE_PERCENT must be between 0 and 1")
	}

	b.Run("JitJSON/Small", func(b *testing.B) {
		shouldParse := shouldParseIterator(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				if shouldParse() {
					_, err := obj.Unmarshal()
					if err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	})

	b.Run("Stdlib/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("JitJSON/Medium", func(b *testing.B) {
		shouldParse := shouldParseIterator(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				if shouldParse() {
					_, err := obj.Unmarshal()
					if err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	})

	b.Run("Stdlib/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("JitJSON/Large", func(b *testing.B) {
		shouldParse := shouldParseIterator(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				if shouldParse() {
					_, err := obj.Unmarshal()
					if err != nil {
						b.Fatal(err)
					}
				}
			}
		}
	})

	b.Run("Stdlib/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

const nestedObjectTemplate = `{
	"nil": null,
	"bool": true,
	"number": 123.45,
	"string": "Hello, World!",
	"slice": [1, "two", false, null],
	"object": null
}`

func generateNestedObjects(depth int) []byte {
	if depth <= 0 {
		return nil
	}
	var obj *Object = &Object{}
	err := json.Unmarshal([]byte(nestedObjectTemplate), obj)
	if err != nil {
		panic(err)
	}

	obj.Object = buildNestedObject(depth - 1)
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return data
}

func buildNestedObject(depth int) *Object {
	if depth <= 0 {
		return nil
	}

	var obj = &Object{}
	err := json.Unmarshal([]byte(nestedObjectTemplate), obj)
	if err != nil {
		panic(err)
	}

	obj.Object = buildNestedObject(depth - 1)
	return obj
}

var (
	smallNestedObjects  = generateNestedObjects(10)
	mediumNestedObjects = generateNestedObjects(1000)
	largeNestedObjects  = generateNestedObjects(100000)
)

// JitObject is a struct that is used to test the nested object parsing.
// The Object field is jitjson.JitJSON[*JitObject] which allows us to catch
// the recursive case by performing just in time unmarshaling.
type JitObject struct {
	Nil    interface{}                  `json:"nil"`
	Bool   bool                         `json:"bool"`
	Number float64                      `json:"number"`
	String string                       `json:"string"`
	Slice  []interface{}                `json:"slice"`
	Object *jitjson.JitJSON[*JitObject] `json:"object"` // marshalled as a pointer
}

// Benchmark 2: Nested object parsing comparison
func BenchmarkNestedParse(b *testing.B) {
	b.Run("JitJSON/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj JitObject
			err := json.Unmarshal(smallNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Stdlib/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj Object
			err := json.Unmarshal(smallNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("JitJSON/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj JitObject
			err := json.Unmarshal(mediumNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Stdlib/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj Object
			err := json.Unmarshal(mediumNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("JitJSON/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj JitObject
			err := json.Unmarshal(largeNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Stdlib/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj Object
			err := json.Unmarshal(largeNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
