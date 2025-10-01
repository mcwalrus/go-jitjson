//go:build go1.25 && goexperiment.jsonv2

package jitjson_test

import (
	jsonv2 "encoding/json/v2"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

// ObjectV2 is the same as Object but for jsonv2 benchmarks
type ObjectV2 struct {
	Nil    interface{}   `json:"nil"`
	Bool   bool          `json:"bool"`
	Number float64       `json:"number"`
	String string        `json:"string"`
	Slice  []interface{} `json:"slice"`
	Object *ObjectV2     `json:"object"` // recursive
}

const objectV2Template = `{
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

func generateObjectsV2(count int) []byte {
	objects := make([]ObjectV2, count)
	for i := 0; i < count; i++ {
		err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
		if err != nil {
			panic(err)
		}
	}
	data, err := jsonv2.Marshal(objects)
	if err != nil {
		panic(err)
	}
	return data
}

var (
	smallDataV2  = generateObjectsV2(10)
	mediumDataV2 = generateObjectsV2(1000)
	largeDataV2  = generateObjectsV2(100000)
)

func shouldParseIteratorV2(parsePercent float64) func() bool {
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

// BenchmarkMarshalV2 benchmarks marshaling performance for both jitjson v2 and encoding/json v2
func BenchmarkMarshalV2(b *testing.B) {
	b.Run("jitjson-v2/Small", func(b *testing.B) {
		objects := make([]ObjectV2, 10)
		for i := 0; i < 10; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSONV2[ObjectV2], 10)
			for j, obj := range objects {
				jitObjects[j] = jitjson.NewV2(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json-v2/Small", func(b *testing.B) {
		objects := make([]ObjectV2, 10)
		for i := 0; i < 10; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := jsonv2.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("jitjson-v2/Medium", func(b *testing.B) {
		objects := make([]ObjectV2, 1000)
		for i := 0; i < 1000; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSONV2[ObjectV2], 1000)
			for j, obj := range objects {
				jitObjects[j] = jitjson.NewV2(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json-v2/Medium", func(b *testing.B) {
		objects := make([]ObjectV2, 1000)
		for i := 0; i < 1000; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := jsonv2.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("jitjson-v2/Large", func(b *testing.B) {
		objects := make([]ObjectV2, 100000)
		for i := 0; i < 100000; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jitObjects := make([]*jitjson.JitJSONV2[ObjectV2], 100000)
			for j, obj := range objects {
				jitObjects[j] = jitjson.NewV2(obj)
			}
			for _, jit := range jitObjects {
				_, err := jit.Marshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("encoding-json-v2/Large", func(b *testing.B) {
		objects := make([]ObjectV2, 100000)
		for i := 0; i < 100000; i++ {
			err := jsonv2.Unmarshal([]byte(fmt.Sprintf(objectV2Template, i)), &objects[i])
			if err != nil {
				b.Fatal(err)
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, obj := range objects {
				_, err := jsonv2.Marshal(obj)
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}

// BenchmarkUnmarshalV2 benchmarks unmarshaling performance for both jitjson v2 and encoding/json v2
func BenchmarkUnmarshalV2(b *testing.B) {
	b.Run("jitjson-v2/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(smallDataV2, &arr)
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

	b.Run("encoding-json-v2/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(smallDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(mediumDataV2, &arr)
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

	b.Run("encoding-json-v2/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(mediumDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(largeDataV2, &arr)
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

	b.Run("encoding-json-v2/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(largeDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkParsePercentageV2 benchmarks the parsing of JSON data with a given percentage of objects
// that are parsed. It compares the performance of JitJSONV2 and the standard library v2.
func BenchmarkParsePercentageV2(b *testing.B) {
	parsePercent, err := strconv.ParseFloat(os.Getenv("PARSE_PERCENTAGE"), 64)
	if err != nil {
		b.Log("PARSE_PERCENTAGE not set, defaulting to 0.3")
		parsePercent = 0.3
	} else {
		b.Logf("PARSE_PERCENTAGE is set to %f", parsePercent)
	}

	if parsePercent < 0 || parsePercent > 1 {
		b.Fatal("PARSE_PERCENTAGE must be between 0 and 1")
	}

	b.Run("jitjson-v2/Small", func(b *testing.B) {
		shouldParse := shouldParseIteratorV2(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(smallDataV2, &arr)
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

	b.Run("encoding-json-v2/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(smallDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Medium", func(b *testing.B) {
		shouldParse := shouldParseIteratorV2(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(mediumDataV2, &arr)
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

	b.Run("encoding-json-v2/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(mediumDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Large", func(b *testing.B) {
		shouldParse := shouldParseIteratorV2(parsePercent)

		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSONV2[ObjectV2]
			err := jsonv2.Unmarshal(largeDataV2, &arr)
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

	b.Run("encoding-json-v2/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []ObjectV2
			err := jsonv2.Unmarshal(largeDataV2, &arr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func buildNestedObjectsV2(num int) []*NestedObjectV2 {
	if num <= 0 {
		return nil
	}
	objs := make([]*NestedObjectV2, num)
	for i := 0; i < num; i++ {
		objs[i] = &NestedObjectV2{}
		data := fmt.Sprintf(objectV2Template, i)
		err := jsonv2.Unmarshal([]byte(data), objs[i])
		if err != nil {
			panic(err)
		}
	}
	return objs
}

func buildNestedObjectsDataV2(num int) []byte {
	objs := buildNestedObjectsV2(num)
	data, err := jsonv2.Marshal(objs)
	if err != nil {
		panic(err)
	}
	return data
}

var (
	smallNestedObjectsDataV2  = buildNestedObjectsDataV2(10)
	mediumNestedObjectsDataV2 = buildNestedObjectsDataV2(100)
	largeNestedObjectsDataV2  = buildNestedObjectsDataV2(1000)
)

// NestedObjectV2 is a struct that is used to test the nested object parsing with jsonv2.
// The Object field is jitjson.JitJSONV2[*ObjectV2] which allows us to catch
// the recursive case by performing just in time unmarshaling.
type NestedObjectV2 struct {
	Nil    interface{}                   `json:"nil"`
	Bool   bool                          `json:"bool"`
	Number float64                       `json:"number"`
	String string                        `json:"string"`
	Slice  []interface{}                 `json:"slice"`
	Object *jitjson.JitJSONV2[*ObjectV2] `json:"object"` // marshalled as a pointer
}

// BenchmarkNestedParseV2 benchmarks nested object parsing comparison with jsonv2
func BenchmarkNestedParseV2(b *testing.B) {
	b.Run("jitjson-v2/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObjectV2
			err := jsonv2.Unmarshal(smallNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json-v2/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*ObjectV2
			err := jsonv2.Unmarshal(smallNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObjectV2
			err := jsonv2.Unmarshal(mediumNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json-v2/Medium", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*ObjectV2
			err := jsonv2.Unmarshal(mediumNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jitjson-v2/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*NestedObjectV2
			err := jsonv2.Unmarshal(largeNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("encoding-json-v2/Large", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var objs []*ObjectV2
			err := jsonv2.Unmarshal(largeNestedObjectsDataV2, &objs)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
