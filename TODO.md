
## Go docs

Go docs is hosted at: ...

## AnyJitJSON

I need to think about the design of AnyJitJSON.

I like the redesign, providing more of a type-safe approach to JSON marshalling / unmarshalling.

However, I've deviated from the initial aim of providing multiple types that you can format into.

Ideally, I want it to be this simple:

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

## Benchmarks

Could you provide a set of benchmarks focused on presenting the performance on:

* When only a few elements / object fields are to be parsed from large json files.

* Overhead of the worst case when parsed from large json files.

* Benchmarking described in a dynamic example of using AnyJitJSON. Use a switch statement to process different types.
