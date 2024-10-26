# Go-JitJSON

`go-jitjson` is a Go library that provides just-in-time (JIT) JSON to allow defered marshaling and unmarshaling until if, or when the data is actually needed.

## Highlights

- 🚀 Improved performance for JSON datasets by avoiding unnecessary parsing
- 💾 Reduced memory usage when working with multiple JSON objects
- 🔄 Seamless integration with existing Go JSON interfaces
- 🏃‍♂️ Improved handling of streaming JSON data

## Installation

```bash
go get github.com/mcwalrus/go-jitjson
```

## Usage

JitJSON provides a generic type `JitJSON[T any]` that can hold either JSON-encoded data or a value of type `T`. 
The type supports both marshaling (Go → JSON) and unmarshaling (JSON → Go) operations, performing the conversion 
only when needed. Both parse operations are expensive when unnecessarily performed and provides opportunity for
conditional parsing on select usage.

Caveat: If you intend to use all json data at once

## Examples

The standard approach is for json marshalling is:

```Go
// Parses everything immediately:
var p Person
err  = json.Unmarshal(data, &p)
if err != nil {
    return nil, err
}
```



```Go
// Delayed parsing:
var jitPerson jitjson.JitJSON[Person]
_ = json.Unmarshal(largeJSON, &jitPerson) 

// Do some work ...
runSomeCode(...)
countToTen(...)
tieSocks(...)

// Then parse!
parsedPerson, _ := jitPerson.Unmarshal()
fmt.Println(parsedPerson)
```

### Updating Existing Code to Use JitJSON

To update an existing codebase to use `jitjson`, you need to replace instances of standard JSON marshaling and unmarshaling with `jitjson` equivalents. Below is an example demonstrating this transition:

For decoding, this looks like:

```Go
// Parse immediately
var person Person
err := json.Unmarshal(jsonData, &person)
if err != nil {
    panic(err)
}
```

Becomes

```Go
// Deferred parsing:
var jitPerson jitjson.JitJSON[Person]
err = json.Unmarshal(jsonData, &jitPerson)
if err != nil {
    panic(err)
}

// Parse only when needed:
person, err := jitPerson.Unmarshal()
if err != nil {
    panic(err)
}
```

Likewise, for encoding:

```Go
// Marshal immediately
jsonData, err := json.Marshal(person)
if err != nil {
    panic(err)
}
```

Becomes

```Go
// Deferred marshaling:
jitPerson := jitjson.NewJitJSON[Person](person)

// Marshal only when needed:
jsonData, err := jitPerson.Marshal()
if err != nil {
    panic(err)
}
```

By following this pattern, you can incrementally update your codebase to leverage the performance and memory benefits of `jitjson`.

jitjson follows the `json.Marshaller` and `json.Unmarshaller` interfaces to allow easy replacement of the standard library where appropriate.


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

### AnyJitJSON 

Sometimes we have a range of types we are trying to parse or reference by jitjson. To do so, we use the `jitjson.AnyJitJSON` construct which can be dynamically type checked to decode later over a range of types.

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
    var jit jitjson.AnyJitJSON

    // ... T of int.
    jit = jitjson.NewJitJSON[int](1)

    // ... of float64.
    jit = jitjson.NewJitJSON[float64](2.0)

    // ... of string.
    jit = jitjson.NewJitJSON[string]("another type!")

    // resolve with json.Marshal
    data, err := json.Marshal(jit)
    if err != nil {
        panic(err)
    }

    // or with type interference
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

Note, this is typically only used for unmarshalling. In the case of marshalling, the type of `jitjson.JitJSON[T]` needs to be known upfront. Attempts to marshal jitjson with `json.Unmarshal(data, &jit)` will work, but will fail when jit.Unmarshal() is called. If `jit` is nil, a nil pointer panic will occur, so `jitjson.AnyJitJSON` so jit will need to be set to a defined `jitjson.JitJSON` type first. 

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
    dec := json.NewDecoder(jit)
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

### Unmarshal array iteratively

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

    // Create JitJSON for array of Person:
    var jitArray []jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonArray, &jitArray)
    if err != nil {
        panic(err)
    }

    // Iterate and print each person:
    for _, jitPerson := range jitArray {
        person, err := jitPerson.Unmarshal()
        if err != nil {
            panic(err)
        }

        fmt.Printf("Person: %+v\n", person)
    }
}
```

### Unmarshal map iteratively

In some cases, you might have a map of JSON objects which you want to iterate across:

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

    // Create JitJSON for map of UUID to Person:
    var jitMap map[uuid.UUID]jitjson.JitJSON[Person]
    err := json.Unmarshal(jsonMap, &jitMap)
    if err != nil {
        panic(err)
    }

    // Iterate and print each person:
    for id, jitPerson := range jitMap {
        person, err := jitPerson.Unmarshal()
        if err != nil {
            panic(err)
        }

        fmt.Printf("ID: %s, Person: %+v\n", id, person)
    }
}
```

This example demonstrates how to handle a map where the keys are UUIDs and the values are JSON objects. The `jitjson.JitJSON` type is used to defer the unmarshaling of the map until it's actually needed.







