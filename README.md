# Go JIT-JSON

Go JitJSON provides 'just-in-time' compilation to encode / decode json.

JitJSON can be applied for any generic types excluding pointers of types and interfaces. JitJSON uses `encoding/json` directly to perform json marshalling and unmarshalling. This allows for easy integration with existing types which either implement json struct tags or marshalling interfaces (json.Marshaler)[https://pkg.go.dev/encoding/json#Marshaler] or (json.Unmarshaler)[https://pkg.go.dev/encoding/json#Unmarshaler].

Using `jitjson` will incur an overhead to perform json validation and reflection on types. This benefit of this library will be gained when marshalling / unmarshalling is not performed for every type or json encoding unless required. 

JitJSON cannot be applied to pointer types and interfaces as:
- Pointer types are stored as **type by JitJSON where **type cannot be parsed by `json.Unmarshal`.
- Interfaces do not perform introspection perform `json.Unmarshal` on the underlying interface value.

## Examples

Just-in-time unmarshalling:
```
import 	"github.com/MaxCollier/go-jitjson"

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

Just-in-time marshalling:
```
import 	"github.com/MaxCollier/go-jitjson"

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