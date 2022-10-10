# Go JIT-JSON

Go JitJSON provides 'just-in-time' compilation for encoding / decoding json to types for any generic type excluding pointer types and interfaces. JitJSON uses `encoding/json` directly to perform json marshalling and unmarshalling, which allows for easy integration with existing types which provide json tags or implement (json.Marshaler)[https://pkg.go.dev/encoding/json#Marshaler] and (json.Unmarshaler)[https://pkg.go.dev/encoding/json#Unmarshaler] interfaces.

JitJSON cannot be applied to pointer types and interfaces as:
- Pointer types are stored as **type by JitJSON where **type cannot be parsed by `json.Unmarshal`.
- Interfaces do not perform introspection perform `json.Unmarshal` on the underlying interface value.

### Examples

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
}
```