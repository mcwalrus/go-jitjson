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

Benchmarks are run with the `-benchmem` flag to show memory allocations.

```bash
go test -bench=. -benchmem
```

To run the benchmarks with a specific percentage of the data parsed, set the `PARSE_PERCENTAGE` environment variable.

```bash
PARSE_PERCENTAGE=0.5 go test -bench='^BenchmarkParsePercentage$' -benchmem
```

Please note, jitjson benchmarks perform relative to the size and volume of data the library is applied to. Monolith applications will benifit the most which reductions in garbage collection cycles considering not all data needs to be parsed.

## Contributing

Please report any issues or feature requests to the [GitHub repository](https://github.com/mcwalrus/go-jitjson).

## About

This module is maintained by Max Collier under an MIT License Agreement.
