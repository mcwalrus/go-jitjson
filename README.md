# Go JitJSON

The go-jitjson library provides a just-in-time (JIT) approach to JSON encoding and decoding in Go. It's designed to be a lightweight wrapper of [encoding/json](https://pkg.go.dev/encoding/json) module. See [module documentation](pkg.go.dev/github.com/mcwalrus/go-jitjson) for details.

## Usage

### Encode with JitJSON:

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
    jit, err := jitjson.NewJitJSON[Person](value)
    if err != nil {
        panic(err)
    }

    // Just-in-time encoding:
    jsonEncoding, err := jit.Marshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonEncoding)) // Output: {"age":30,"city":"New York","name":"John"}
}
```

### Decode with JitJSON:

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
    jit, err := jitjson.NewJitJSON[Person](jsonEncoding)
    if err != nil {
        panic(err)
    }
    
    // Just-in-time decoding:
    value, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(value) // Output: {John 30 New York}
}
```
