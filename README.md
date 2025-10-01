# Go-JitJSON

[![Go Version](https://img.shields.io/github/go-mod/go-version/mcwalrus/go-jitjson)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/mcwalrus/go-jitjson)](https://goreportcard.com/report/github.com/mcwalrus/go-jitjson)
[![codecov](https://codecov.io/gh/mcwalrus/go-jitjson/branch/main/graph/badge.svg)](https://codecov.io/gh/mcwalrus/go-jitjson) 
[![GoDoc](https://godoc.org/github.com/mcwalrus/go-jitjson?status.svg)](https://godoc.org/github.com/mcwalrus/go-jitjson)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go library to provide lazy just-in-time (JIT) parsing of JSON encodings and values.

## Motivation

Traditional parsing with [json.Marshal](https://pkg.go.dev/encoding/json#Marshal) or [json.Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal) processes all data immediately, even if it may never be used. Unnecessary parsing leads to wasted CPU cycles on unused data, unnecessary memory allocations, and increased pressure on garbage collection operations. If you intend to parse all your data, jitjson will not provide any benefit. Think of the library as a lazy two-way parser with caching implemented.

## Key Features

- 🪶 Zero dependencies
- 🛠️ Support for encoding/json/v2
- 🔗 Integrates with encoding/json interfaces
- 💾 Reduce memory when working with large JSON datasets 
- 🚀 Improves performance by avoiding unnecessary parsing of JSON

## Installation

This library requires Go version >=1.18:

```bash
$ go get github.com/mcwalrus/go-jitjson
```

### Using json/v2

For Go 1.25+, jitjson provides support for [encoding/json/v2](https://pkg.go.dev/encoding/json/v2). 

```bash
$ export GOEXPERIMENT=jsonv2
$ go doc github.com/mcwalrus/go-jitjson.JitJSONV2
```

For more information on json/v2, see the [offical blog post](https://go.dev/blog/jsonv2-exp).

## Quick Start

### Marshaling

Use the `New` method to create a `JitJSON` from a Go value of any type.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jit := jitjson.New(Person{
        Name: "John",
        Age:  30,
        City: "New York",
    })

    jsonEncoding, err := jit.Marshal() // just-in-time ...
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonEncoding)) // Output: {"age":30,"city":"New York","name":"John"}
}
```

### Unmarshaling

Use the `NewFromBytes` method to create a `JitJSON` from a JSON encoded string.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jit := jitjson.NewFromBytes[Person]([]byte(`{
        "Name": "John",
        "Age": 30,
        "City": "New York"
    }`))

    value, err := jit.Unmarshal() // just-in-time ...
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

### Using Slices

Consider when your target value is a slice:

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jsonArray := []byte(`[
        {"Name":"John","Age":30,"City":"New York"},
        {"Name":"Jane","Age":25,"City":"Los Angeles"},
        {"Name":"Jim","Age":35,"City":"Chicago"}
    ]`)

    // Unmarshal slice
    var jit []*jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonArray, &jit)
    if err != nil {
        panic(err)
    }

    // Unmarshal only the first person ...
    value, err := jit[0].Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

### Using Maps

Consider when your target value is a map:

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jsonMap := []byte(`{
        1: {"Name":"John","Age":30,"City":"New York"},
        2: {"Name":"Jane","Age":25,"City":"Los Angeles"},
        3: {"Name":"Jim","Age":35,"City":"Chicago"}
    }`)

    // Unmarshal map
    var jitMap map[int]*jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonMap, &jitMap)
    if err != nil {
        panic(err)
    }

    // Select a person
    jit, ok := jitMap[1]
    if !ok {
        panic("missing person")
    }

    // Unmarshal only person (1) ...
    person, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(person) // Output: {John 30 New York}
}
```

### Nested Fields

Consider when you have nested fields you would want to avoid parsing:

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Address struct {
    Street string
    City   string
    Zip    string
}

type Person struct {
    Name    string
    Age     int
    Address *jitjson.JitJSON[Address]    
}

func main() {
    jsonData := []byte(`{
        "Name": "John",
        "Age": 30,
        "Address": {
            "Street": "123 Main St",
            "City": "New York",
            "Zip": "10001"
        }
    }`)

    // Unmarshal person
    var person Person
    err := json.Unmarshal(jsonData, &person)
    if err != nil {
        panic(err)
    }

    // Unmarshal person's address just-in-time ...
    address, err := person.Address.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(address) // Output: {123 Main St New York 10001}
}
```

### Updating Values

Use the `Set` \ `SetBytes` methods to update values of `JitJSON`. Using these methods will force re-parsing of results:

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
    // New JitJSON:
    jit := jitjson.New(Person{
        Name: "John",
        Age:  30,
        City: "New York",
    })

    // Initial encoding ...
    jsonEncoding, err := jit.Marshal()
    if err != nil {
        panic(err)
    }

    // Update the value:
    jit.Set(Person{
        Name: "Jane",
        Age:  25,
        City: "Los Angeles",
    })

    // Updated encoding ...
    jsonEncoding, err = jit.Marshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonEncoding)) // Output: {"age":25,"city":"Los Angeles","name":"Jane"}
}
```

### Parsing Multiple Times

Values can be parsed multiple times, which doesn't come with any performance penalty.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jit := jitjson.New(Person{
        Name: "John",
        Age:  30,   
        City: "New York",
    })

     // Initial Marshal ...
    _, err := jit.Marshal()
    if err != nil {
        panic(err)
    }

    // No new allocations
    for i := 0; i < 10; i++ {
        _, err = jit.Marshal()
        if err != nil {
            panic(err)
        }
    }
◊
    jit = jitjson.NewFromBytes([]byte(`{
        "name": "John",
        "age": 30,
        "city": "New York"
    }`))

    // Initial Unmarshal ...
    _, err = jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    // No new allocations
    for i := 0; i < 10; i++ {
        _, err = jit.Unmarshal()
        if err != nil {
            panic(err)
        }
    }
}
```

## Benchmarks

The library provides comprehensive benchmarks for both `encoding/json` (v1) and `encoding/json/v2` implementations. Note jitjson benchmarks perform relative to the size and volume of data the library is applied to. Benefits across applications will be varied with most significant effects in reduced pressure in garbage collection cycles considering different ways data can avoid being parsed.

### Run Benchmarks

An assumed data structure is used for the benefits of jitjson benchmarks:

```bash
export GOEXPERIMENT=jsonv2 
# Run all benchmarks (encoding/json v1)
go test -bench=. -benchmem
# Run all benchmarks (encoding/json v2) - requires Go 1.25+
go test -bench=. -benchmem
```

### Marshaling Benchmarks Worst-Case

These benchmarks perform worst-case analysis by marshaling all objects in the dataset:

```bash
go test -bench='^BenchmarkMarshalWorstCase$' -benchmem
go test -bench='^BenchmarkMarshalWorstCaseV2$' -benchmem
```

### Unmarshaling Benchmarks Worst-Case

These benchmarks perform worst-case analysis by unmarshaling all objects in the dataset:

```bash
go test -bench='^BenchmarkUnmarshalWorstCase$' -benchmem
go test -bench='^BenchmarkUnmarshalWorstCaseV2$' -benchmem
```

### Partial Dataset Parsing Benchmarks

These benchmarks test parsing at a configurable percentage of the total objects available:

```bash
PARSE_PERCENTAGE=0.5 go test -bench='^BenchmarkParsePercentage$' -benchmem
PARSE_PERCENTAGE=0.5 go test -bench='^BenchmarkParsePercentageV2$' -benchmem
```

### Nested JitJSON Structs Benchmarks

Test performance with nested jitjson structures:

```bash
go test -bench='^BenchmarkNestedParse$' -benchmem
go test -bench='^BenchmarkNestedParseV2$' -benchmem
```

### Custom Benchmarking

For your own accurate performance evaluation, consider creating benchmarks with your own data:

```go
// use json or jsonv2
import "encoding/json"

type MyStruct struct { /* ... */ } 

func BenchmarkMyData(b *testing.B) {

    // Read json from file
    data, err := os.ReadFile(os.Getenv("JSON_FILE"))
    if err != nil {
        b.Fatalf("failed to retrieve json-data: %v", err)
    }
    
    // Parse only what you need
    b.Run("jitjson", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var jit []*jitjson.JitJSON[MyStruct]
            json.Unmarshal(data, &jit)
            for _, item := range jit[:len(jit)/2] {
                item.Unmarshal()
            }
        }
    })
    
    // Compare against the standard library
    b.Run("encoding-json", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var result []MyStruct
            json.Unmarshal(data, &result)
        }
    })
}
```

## Contributing

Please report any issues or feature requests to the [GitHub repository](https://github.com/mcwalrus/go-jitjson).

## About

This module is maintained by Max Collier under an MIT License Agreement.
