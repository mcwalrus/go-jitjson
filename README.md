# Go JIT-JSON

Go JitJSON provides 'just-in-time' compilation to encode / decode json.

## Usage

JitJSON can be applied for any generic types excluding pointers of types and interfaces. JitJSON uses `encoding/json` to perform json parsing, which allows for easy integration with existing types which already implement json tags or [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) or [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) methods.

### Tradeoffs

Using `jitjson` will incur an overhead to perform json validation and reflection on types. This benefit of the library is gained when json parsing may not be required. 

## Examples

### Unmarshalling:
```
import 	"github.com/mcwalrus/go-jitjson"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
    data := []byte(`
        {
            "name": "Willy Wonka",
            "age":  42
        }
    `)

    jit, err := jitjson.NewJitJSON[Person](data)
    if err != nil {
        panic(err)
    }

    // performs 'just-in-time' compilation with jit.Unmarshal().
    person, err := jit.Unmarshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(person)
    // output: {Name:Willy Wonka Age:42}
}
```

### Marshalling:
```
import 	"github.com/mcwalrus/go-jitjson"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
    person := Person{
        Name: "Charlie Bucket",
        Age:  12,
    }

    jit, err := jitjson.NewJitJSON[Person](nil)
    if err != nil {
        panic(err)
    }

    // performs 'just-in-time' compilation with jit.Marshal().
    jit.Set(person)
    data, err := jit.Marshal()
    if err != nil {
        panic(err)
    }

    fmt.Println(string(data))
    // output: {"name":"Charlie Bucket","age":12}
}
```

## Unhandled cases

JitJSON cannot be applied to pointer types and interfaces as:
- Pointer types are stored as `**type` by JitJSON, where `**type` cannot be parsed by `json.Unmarshal`.
- Interfaces do not perform introspection perform `json.Unmarshal` on the underlying interface value.

These maybe implemented on request.
