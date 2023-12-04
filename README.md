# Go-JitJSON

go-jitjson provides a just-in-time (JIT) approach to JSON encoding and decoding in Go. It's designed to be a lightweight wrapper of the [encoding/json](https://pkg.go.dev/encoding/json) module. The library provides a type `JitJSON[T any]`, which can hold either a JSON encoding or a value of any type `T`. The `JITInterface` interface type enables dynamic type assignment for `JitJSON[T any]`, where `T` can be any type. 

See [API reference](https://pkg.go.dev/github.com/mcwalrus/go-jitjson) for more documentation.

## Usage

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

### Dynamic type-set with JitJSON:

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
        jit jitjson.JITInterface
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


## About

This module is maintained by Max Collier under an MIT License Agreement.
