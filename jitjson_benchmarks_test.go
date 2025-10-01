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

// BenchmarkMarshalWorstCase benchmarks marshaling performance for both jitjson and encoding/json.
// This performs worst-case analysis by marshaling all objects in the dataset.
func BenchmarkMarshalWorstCase(b *testing.B) {
	b.Run("jitjson/Small", func(b *testing.B) {
		objects := make([]Object, 10)
		for i := 0; i < 10; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSON[Object], 10)
			for j, obj := range objects {
				jitObjects[j] = jitjson.New(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Small", func(b *testing.B) {
		objects := make([]Object, 10)
		for i := 0; i < 10; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := json.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("jitjson/Medium", func(b *testing.B) {
		objects := make([]Object, 1000)
		for i := 0; i < 1000; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSON[Object], 1000)
			for j, obj := range objects {
				jitObjects[j] = jitjson.New(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Medium", func(b *testing.B) {
		objects := make([]Object, 1000)
		for i := 0; i < 1000; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := json.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("jitjson/Large", func(b *testing.B) {
		objects := make([]Object, 100000)
		for i := 0; i < 100000; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSON[Object], 100000)
			for j, obj := range objects {
				jitObjects[j] = jitjson.New(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Large", func(b *testing.B) {
		objects := make([]Object, 100000)
		for i := 0; i < 100000; i++ {
			err := json.Unmarshal([]byte(fmt.Sprintf(objectTemplate, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := json.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// BenchmarkUnmarshalWorstCase benchmarks unmarshaling performance for both jitjson and encoding/json.
// This performs worst-case analysis by unmarshaling all objects in the dataset.
func BenchmarkUnmarshalWorstCase(b *testing.B) {
	b.Run("jitjson/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkParsePercentage benchmarks the parsing of JSON data with a given percentage of objects
// that are parsed. It compares the performance of JitJSON and the standard library.
func BenchmarkParsePercentage(b *testing.B) {
	parsePercent, err := strconv.ParseFloat(os.Getenv("PARSE_PERCENTAGE"), 64)
	if err != nil {
		b.Log("PARSE_PERCENTAGE not set, defaulting to 50%")
		parsePercent = 0.5
	} else {
		b.Logf("PARSE_PERCENTAGE is set to %f", parsePercent)
	}

	if parsePercent < 0 || parsePercent > 1 {
		b.Fatal("PARSE_PERCENTAGE must be between 0 and 1")
	}

	b.Run("jitjson/Small", func(b *testing.B) {
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

	b.Run("encoding-json/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Medium", func(b *testing.B) {
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

	b.Run("encoding-json/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Large", func(b *testing.B) {
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

	b.Run("encoding-json/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []Object
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func buildNestedObjects(num int) []*NestedObject {
	if num <= 0 {
		return nil
	}
	objs := make([]*NestedObject, num)
	for i := 0; i < num; i++ {
		objs[i] = &NestedObject{}
		data := fmt.Sprintf(objectTemplate, i)
		err := json.Unmarshal([]byte(data), objs[i])
		if err != nil {
			panic(err)
		}
	}
	return objs
}

func buildNestedObjectsData(num int) []byte {
	objs := buildNestedObjects(num)
	data, err := json.Marshal(objs)
	if err != nil {
		panic(err)
	}
	return data
}

// var (
// 	smallNestedObjects  = buildNestedObjects(10)
// 	mediumNestedObjects = buildNestedObjects(100)
// 	largeNestedObjects  = buildNestedObjects(1000)
// )

var (
	smallNestedObjectsData  = buildNestedObjectsData(10)
	mediumNestedObjectsData = buildNestedObjectsData(100)
	largeNestedObjectsData  = buildNestedObjectsData(1000)
)

// JitObject is a struct that is used to test the nested object parsing.
// The Object field is jitjson.JitJSON[*JitObject] which allows us to catch
// the recursive case by performing just in time unmarshaling.
type NestedObject struct {
	Nil    interface{}               `json:"nil"`
	Bool   bool                      `json:"bool"`
	Number float64                   `json:"number"`
	String string                    `json:"string"`
	Slice  []interface{}             `json:"slice"`
	Object *jitjson.JitJSON[*Object] `json:"object"` // marshalled as a pointer
}

// Benchmark 2: Nested object parsing comparison
func BenchmarkNestedParse(b *testing.B) {
	b.Run("jitjson/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObject
			err := json.Unmarshal(smallNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*Object
			err := json.Unmarshal(smallNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObject
			err := json.Unmarshal(mediumNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*Object
			err := json.Unmarshal(mediumNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObject
			err := json.Unmarshal(largeNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*Object
			err := json.Unmarshal(largeNestedObjectsData, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
