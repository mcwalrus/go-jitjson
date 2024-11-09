# Go-JitJSON

`go-jitjson` is a Go library that provides just-in-time (JIT) JSON to allow defered marshaling and unmarshaling until if, or when the data is actually needed.

## Key Features

- 🚀 Improve performance for JSON datasets by avoiding unnecessary parsing
- 💾 Reduce memory usage when working with multiple JSON objects
- 🔄 Seamless integration with existing Go JSON interfaces
- 🏃‍♂️ Improve handling of streaming JSON data
- 🧩 Dynamic type parsing of JSON

## Installation

```bash
go get github.com/mcwalrus/go-jitjson
```

## Usage

The module provides a generic `JitJSON[T any]` as a lazy JSON parser which can take .
The type supports both marshaling (Go → JSON) and unmarshaling (JSON → Go) operations, to perform conversions
only when needed. Both parsing operations are expensive when performed unnecessarily, which can be avoided by
conditional parsing.

`AnyJitJSON` is a primary means of unmarshaling JSON of a dynamic structure. This leverages JitJSON to store JSON
json, but will provide a significantly higher overhead compared to the standard form of parsing? This should be
considered for cases where parsing is a dynamic process. AnyJitJSON provides certainty over types.

## Examples

### Encoding with JitJSON:

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

    var value = Person{
        Name: "John",
        Age:  30,
        City: "New York",
    }

    // Create JitJSON:
    jit := jitjson.NewJitJSON[Person](value)

    // Just-in-time encoding:
    jsonEncoding, err := jit.Marshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonEncoding)) // Output: {"age":30,"city":"New York","name":"John"}
}
```

### Decoding with JitJSON:

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
    jsonEncoding := []byte(`{"Name":"John","Age":30,"City":"New York"}`)
    
    // Create JitJSON:
    jit := jitjson.NewJitJSON[Person](jsonEncoding)
    
    // Just-in-time decoding:
    value, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

Benefit: `jitjson.JitJSON[T]` provides `json.Marshaller` and `json.Unmarshaller` interface methods to allow easy replacement of the standard library where appropriate.

### Custom json.Decoder:

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
    jsonData := []byte(`{"Name":"John","Age":30,"City":"New York"}`)

    jit := jitjson.BytesToJitJSON[Person](jsonData)

    // Create a json.Decoder:
    dec := json.NewDecoder(jit.JSON JSON JSON)
    dec.DisallowUnknownFields()

    // Decode Person:
    var p Person
    err := dec.Decode(&p)
    if err != nil {
        panic(err)
    }

    fmt.Println(p) // Output: {John 30 New York}
}
```

### Unmarshalling a slice:

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
        {"Name":"Jane","Age":25,"City":"Los Angeles"}
    ]`)

    // A JitJSON slice of People:
    var jitSlice []jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonArray, &jitSlice)
    if err != nil {
        panic(err)
    }

    // Unmarshal only the first index:
    jit, err = jitSlice[0].Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

### Unmarshalling a map:

```Go
package main

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    jsonMap := []byte(`{
        "550e8400-e29b-41d4-a716-446655440000": {"Name":"John","Age":30,"City":"New York"},
        "550e8400-e29b-41d4-a716-446655440001": {"Name":"Jane","Age":25,"City":"Los Angeles"}
    }`)

    // A JitJSON map of UUIDs to People:
    var jitMap map[uuid.UUID]jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonMap, &jitMap)
    if err != nil {
        panic(err)
    }


    jit, ok := jitMap["550e8400-e29b-41d4-a716-446655440000"]
    if !ok {
        panic("missing person")
    }
    person, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }
    
    fmt.Println(person) // Output: {John 30 New York}
}
```

### Unmarshalling nested structures

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
    Address jitjson.JitJSON[Address]    
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

    jit := jitjson.NewJitJSON[Person](jsonData)

    // Decode person
    person, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    // Decode the address
    address, err := person.Address.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(address) // Output: {123 Main St New York 10001}
}
```

### Dynamic Type Parsing

To handle dynamic parsing of JSON, we can use `AnyJitJSON` to optionally set `NewJitJSON[T]` types.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
    var jit jitjson.AnyJitJSON

    // Support for multiple types
    jit = jitjson.NewJitJSON[int](1)
    jit = jitjson.NewJitJSON[float64](2.0)
    jit = jitjson.NewJitJSON[string]("another type!")

    // Resolve by json.Marshal
    data, err := json.Marshal(jit)
    if err != nil {
        panic(err)
    }

    // Unmarshal by type inference
    v := (jit).(jitjson.JitJSON[string])
    s, err := v.Unmarshal()
    if err != nil {
        panic(err)
    }

    // Output: another type!
    fmt.Println(string(data))
    fmt.Println(s)
}
```

## About

This module is maintained by Max Collier under an MIT License Agreement.
