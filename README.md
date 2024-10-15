# Go-JitJSON

`go-jitjson` is a Go library that provides just-in-time (JIT) JSON parsing capability to defer marshaling and unmarshaling processes until they are actually needed. The library supports type safety through use of generics and uses standard Go interfaces (`json.Marshaler`, `json.Unmarshaler`, and `io.Reader`) to make it easy to integrate with existing projects.

## Installation

```bash
go get github.com/mcwalrus/gp-jitjson
```

## Usage

The library provides a generic type `JitJSON[T any]`, capable of holding either JSON-encoded data or a value of type `T`. The `AnyJitJSON` interface allows for flexible dynamic type handling with `JitJSON[T any]`, supporting any Go type. Additionally, `JitJSON[T any]` implements the `io.Reader` interface, enabling integration with `json.Decoder`.

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

### Dynamic type assignment:

```Go
package main

import (
    "fmt"
    "github.com/mcwalrus/go-jitjson"
)

func main() {
	
    // JitJSON.
    var (
		err error
        jit jitjson.AnyJitJSON
	)

    // ... T of int.
	jit, err = jitjson.NewJitJSON[int](1)
	if err != nil {
		panic(err)
	}

    // ... of float64.
    jit, err = jitjson.NewJitJSON[float64](2.0)
	if err != nil {
		panic(err)
	}

    // ... of string.
	jit, err = jitjson.NewJitJSON[string]("another type!")
	if err != nil {
		panic(err)
	}

    // Convert to JitJSON[T] type to unmarshal: 
    v := (jit).(jitjson.JitJSON[string])
    s, err := v.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(s) // Output: another type!
}
```

### Custom decoders:

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

	// Create JitJSON:
    jit, err := jitjson.NewJitJSON[Person](jsonData)
    if err != nil {
        panic(err)
    }

    // Create a json.Decoder:
	dec := json.NewDecoder(jit)
	dec.DisallowUnknownFields()

    // Decode Person:
	var p Person
	err = dec.Decode(&p)
	if err != nil {
        panic(err)
    }

    fmt.Println(p) // Output: {John 30 New York}
}
```

## About

This module is maintained by Max Collier under an MIT License Agreement.
