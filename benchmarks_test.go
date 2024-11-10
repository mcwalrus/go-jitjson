package jitjson_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

type Object struct {
	Nil    interface{}   `json:"nil"`
	Bool   bool          `json:"bool"`
	Number float64       `json:"number"`
	String string        `json:"string"`
	Slice  []interface{} `json:"slice"`
	Object *NestedObject `json:"object"`
}

type NestedObject struct {
	String string        `json:"nestedString"`
	Number int           `json:"nestedNumber"`
	Object *NestedObject `json:"nestedObject"` // recursive
}

// Generate a slice of JSON objects of a given count.
func generateArrayOfObjects(count int) []byte {
	objects := make([]string, count)
	for i := 0; i < count; i++ {
		objects[i] = fmt.Sprintf(`{
			"nil": null,
			"bool": true,
			"number": 123.45,
			"string": "Hello, World!",
			"slice": [1, "two", false, null, %d],
			"object": {
				"nestedString": "Nested",
				"nestedNumber": 456,
				"nestedObject": {
					"nestedString": "Deep Nested",
					"nestedNumber": 678,
					"nestedObject": null
				}
			}
		}`, i)
	}
	return []byte(fmt.Sprintf(`[%s]`, strings.Join(objects, ",")))
}

var (
	smallData  = generateArrayOfObjects(10)
	mediumData = generateArrayOfObjects(100)
	largeData  = generateArrayOfObjects(1000)
)

// Benchmark 1: Full parsing comparison
func BenchmarkFullParse(b *testing.B) {
	b.Run("JitJSON/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(smallData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
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
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(mediumData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
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
		for i := 0; i < b.N; i++ {
			var arr []*jitjson.JitJSON[Object]
			err := json.Unmarshal(largeData, &arr)
			if err != nil {
				b.Fatal(err)
			}

			// just in time unmarshal
			for _, obj := range arr {
				_, err := obj.Unmarshal()
				if err != nil {
					b.Fatal(err)
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

// shouldParseIterator returns a function that returns true if the percentage has been reached.
// It is used to generate a random percentage of the time that a parse should be done.
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

// Benchmark 1: Partial Comparison
func BenchmarkPartialParse(b *testing.B) {
	var parsePercent float64 = 0.4

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

const (
	objectData = `{
		"nil": null,
		"bool": true,
		"number": 123.45,
		"string": "Hello, World!",
		"slice": [1, "two", false, null],
		"object": {
			"nestedString": "Nested",
			"nestedNumber": 456,
			"nestedObject": {
				"nestedString": "Deep Nested",
				"nestedNumber": 678,
				"nestedObject": null
			}
		}
	}`
)

func buildNestedObject(depth int) *NestedObject {
	if depth <= 0 {
		return nil
	}
	return &NestedObject{
		String: "Nested",
		Number: depth,
		Object: buildNestedObject(depth - 1),
	}
}

func generateNestedObjects(depth int) []byte {
	if depth <= 0 {
		return nil
	}

	var obj *Object = &Object{}
	err := json.Unmarshal([]byte(objectData), obj)
	if err != nil {
		panic(err)
	}

	nestedObj := buildNestedObject(depth)
	obj.Object = nestedObj
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return data
}

var (
	smallNestedObjects  = generateNestedObjects(10)
	mediumNestedObjects = generateNestedObjects(100)
	largeNestedObjects  = generateNestedObjects(1000)
)

type JitObject struct {
	Nil    interface{}                        `json:"nil"`
	Bool   bool                               `json:"bool"`
	Number float64                            `json:"number"`
	String string                             `json:"string"`
	Slice  []interface{}                      `json:"slice"`
	Object *jitjson.JitJSON[*JitNestedObject] `json:"object"` // marshalled as a pointer
}

type JitNestedObject struct {
	String string                             `json:"nestedString"`
	Number int                                `json:"nestedNumber"`
	Object *jitjson.JitJSON[*JitNestedObject] `json:"nestedObject"` // recursive
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

// Implement with percentage iterator. If The percentage is 100% then we are doing the worst case. 0% is best case.
func BenchmarkNestedParseWorstCase(b *testing.B) {
	b.Run("JitJSON/Small", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var obj JitObject
			err := json.Unmarshal(smallNestedObjects, &obj)
			if err != nil {
				b.Fatal(err)
			}

			// Unmarshal all nested objects recursively
			nested, err := obj.Object.Unmarshal()
			if err != nil {
				b.Fatal(err)
			}
			for nested != nil && nested.Object != nil {
				nested, err = nested.Object.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
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

			// Unmarshal all nested objects recursively
			nested, err := obj.Object.Unmarshal()
			if err != nil {
				b.Fatal(err)
			}
			for nested != nil && nested.Object != nil {
				nested, err = nested.Object.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
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

			// Unmarshal all nested objects recursively
			nested, err := obj.Object.Unmarshal()
			if err != nil {
				b.Fatal(err)
			}
			for nested != nil && nested.Object != nil {
				nested, err = nested.Object.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
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

// type Animal string

// const (
// 	AnimalCat   Animal = "cat"
// 	AnimalDog   Animal = "dog"
// 	AnimalSnake Animal = "snake"
// 	AnimalMoa   Animal = "moa"
// )

// type AnimalObject interface {
// 	AnimalNoise() string
// }

// type Dog struct {
// 	Name string   `json:"name"`
// 	Age  int      `json:"age"`
// 	Diet []string `json:"diet"`
// }

// func (d Dog) AnimalNoise() string {
// 	return "woof"
// }

// type Cat struct {
// 	Name       string   `json:"name"`
// 	Age        int      `json:"age"`
// 	VetHistory []string `json:"vetHistory"`
// }

// func (c Cat) AnimalNoise() string {
// 	return "meow"
// }

// type Snake struct {
// 	Name           string   `json:"name"`
// 	Age            int      `json:"age"`
// 	Venomous       bool     `json:"venomous"`
// 	PlacesOfOrigin []string `json:"placesOfOrigin"`
// }

// func (s Snake) AnimalNoise() string {
// 	return "hiss"
// }

// type Moa struct {
// 	Name    string `json:"name"`
// 	Age     int    `json:"age"`
// 	Height  int    `json:"height"`
// 	Weight  int    `json:"weight"`
// 	Extinct bool   `json:"extinct"`
// }

// func (m Moa) AnimalNoise() string {
// 	return "honk"
// }

// const (
// 	animalDogData = `{
// 		"type": "dog",
// 		"name": "Buddy",
// 		"age": 5,
// 		"diet": ["meat", "bones"]
// 	}`
// 	animalCatData = `{
// 		"type": "cat",
// 		"name": "Whiskers",
// 		"age": 3,
// 		"vetHistory": ["vaccinations", "checkup"]
// 	}`
// 	animalSnakeData = `{
// 		"type": "snake",
// 		"name": "Slither",
// 		"age": 2,
// 		"venomous": true,
// 		"placesOfOrigin": ["Australia", "Africa", "Asia", "Europe", "North America", "South America"]
// 	}`
// 	animalMoaData = `{
// 		"type": "moa",
// 		"name": "Moa",
// 		"age": 10,
// 		"extinct": true
// 	}`
// )

// // Benchmark the case where we have a JSON string and we want to parse it into different types
// // based on the "type" field. Provide a benchmark for the case where we use AnyJitJSON and another
// // for the standard library.
// func BenchmarkAnyJitJSON_Bench2(b *testing.B) {

// }

// func BenchmarkStdlib_Bench2(b *testing.B) {

// }
