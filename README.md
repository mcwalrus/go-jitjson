# Go-JitJSON

`go-jitjson` is a Go library to provide defered just-in-time (JIT) JSON parsing for if and when the data is actually needed.

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

## When to Use?

Use jitjson in cases where you can conditionally avoid parsing json data. With avoidance, bypasses to `encoding/json` can improve CPU op/s, memory allocation and reduce the number of garbage collection (GC) cycles performed. There are other cases to consider for partially parsing json data. You can parse json from or into structs with nested jitjson fields, or between slices and maps with jitjson elements. For dynamic unmarshalling of json data that you do know the specifcation of, use `AnyJitJSON`. If you intend to parse all your data, jitjson will not provide any benefit. This library is intended for use in real-time, monolithic, or large data driven applications.


## Quick Start

### Marshaling

Use the `New` method to create a `JitJSON` from a value of any type.

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
    // Create new JitJSON:
    jit := jitjson.New(Person{
        Name: "John",
        Age:  30,
        City: "New York",
    })

    // Marshal value just-in-time:
    jsonEncoding, err := jit.Marshal()
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
    // Create new JitJSON:
    jit := jitjson.NewFromBytes[Person]([]byte(`{
        "Name": "John",
        "Age": 30,
        "City": "New York"
    }`))

    // Unmarshal value just-in-time:
    value, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

### Updating Values

Use the `Set` method to update the value of a `JitJSON`.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
    // Create new JitJSON:
    jit := jitjson.New(Person{
        Name: "John",
        Age:  30,
        City: "New York",
    })

    // Marshal the initial value:
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

    // Marshal the updated value:
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

    _, err := jit.Marshal() // Initial Marshal
    if err != nil {
        panic(err)
    }
    for i := 0; i < 10; i++ {
        _, err = jit.Marshal() // No new allocations
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

    _, err = jit.Unmarshal() // Initial Unmarshal
    if err != nil {
        panic(err)
    }
    for i := 0; i < 10; i++ {
        _, err = jit.Unmarshal() // No new allocations
        if err != nil {
            panic(err)
        }
    }
}
```

### Advanced Usage

#### Using Slices

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

    // Unmarshal the first person just-in-time
    value, err := jit[0].Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```

#### Using Maps

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

    // Unmarshal only person one just-in-time
    person, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(person) // Output: {John 30 New York}
}
```

#### Nested Fields

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

    // Unmarshal person's address just-in-time
    address, err := person.Address.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(address) // Output: {123 Main St New York 10001}
}
```

### AnyJitJSON

#### Basic Usage

Dynamic type inference of `AnyJitJSON`.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

type Person struct {
    Name    string
    Age     int
    Friends []Person
}

func main() {
    jsonData := []byte(`{
        "Name": "John",
        "Age": 30,
        "Friends": [
            {"Name": "Jane", "Age": 25},
            {"Name": "Jim", "Age": 35},
            {"Name": "Jill", "Age": 45}
        ]
    }`)

    // Support for multiple types
    var jit jitjson.AnyJitJSON
    err := json.Unmarshal(jsonData, &jit)
    if err != nil {
        panic(err)
    }

    // Get the object
    obj, ok := jit.AsObject()
    if !ok {
        panic("not object")
    }

    // Get the name
    if name, ok := obj["Name"].AsString(); ok {
        fmt.Println(name)
    }

    // Get the friends
    if friends, ok := obj["Friends"].AsArray(); ok {
        for _, friend := range friends {
            fmt.Println(friend)
        }

        data, err := friends[0].Marshal()
        if err != nil {
            panic(err)
        }

        // Unmarshal the friend just-in-time
        var person Person
        err = json.Unmarshal(data, &person)
        if err != nil {
            panic(err)
        }
    }
}
```

#### With Arrays

Dynamic type inference of `AnyJitJSON` with arrays.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
    jsonData := []byte(`[
        1.23,
        "Hello, world!",
        {"Name": "John", "Age": 30},
        true
    ]`)

    // Support for multiple types
    var jit []*jitjson.AnyJitJSON
    err := json.Unmarshal(jsonData, &jit)
    if err != nil {
        panic(err)
    }

    num, ok := jit[0].AsNumber()
    if !ok {
        panic("not a number")
    }
    fmt.Println(num) // Output: 1.23

    str, ok := jit[1].AsString()
    if !ok {
        panic("not a string")
    }
    fmt.Println(str) // Output: Hello, world!

    if jit[2].Type() != jitjson.Object {
        panic("not an object")
    }
    fmt.Println(jit[2]) // Output: {"Name": "John", "Age": 30}

    if jit[3].Type() != jitjson.Boolean {
        panic("not a boolean")
    }
    fmt.Println(jit[3]) // Output: true
}
```

#### Conditional Types

Dynamic type inference of `AnyJitJSON` across multiple possible types.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func whichType(data []byte) {
    var jit *jitjson.AnyJitJSON
    err := json.Unmarshal(data, &jit)
    if err != nil {
        panic(err)
    }
    switch typ := jit.Type(); typ {
    case jitjson.Null:
        fmt.Println("null")
    case jitjson.Object:
        fmt.Println("Hmmm, an object?")
    case jitjson.Array:
        fmt.Println("An array? Interesting...")
    default:
        fmt.Println("Huh, I have no idea what this is...")
    }
}

func main() {
    whichType([]byte(`null`)) // Output: null
    whichType([]byte(`{"Name": "John", "Age": 30}`)) // Output: Hmmm, an object?
    whichType([]byte(`[1, 2, 3]`)) // Output: An array? Interesting...
    whichType([]byte(`true`)) // Output: Huh, I have no idea what this is...
}
```

## Benchmarks

Benchmarks are run with the `-benchmem` flag to show memory allocations.

```bash
go test -bench=. -benchmem
```

To run the benchmarks with a specific percentage of the data parsed, set the `PARSE_PERCENTAGE` environment variable.

```bash
PARSE_PERCENTAGE=0.3 go test -bench='^BenchmarkParsePercentage$' -benchmem
```

Please note, jitjson benchmarks perform relative to the size and volume of data the library is applied to. Monolith applications will benifit the most which reductions in garbage collection cycles considering not all data needs to be parsed.

## Contributing

Please report any issues or feature requests to the [GitHub repository](https://github.com/mcwalrus/go-jitjson).

## About

This module is maintained by Max Collier under an MIT License Agreement.
