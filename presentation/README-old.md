## Key Features

- 🪶 Zero dependencies
- 🧩 Dynamic type parsing of JSON
- 🔗 Integrates with encoding/json interface types
- 💾 Reduce memory when working with large JSON datasets 
- 🚀 Improves performance by avoiding unnecessary parsing of JSON
- 🛠️ Customisable parsing

[![CI/CD](https://github.com/mcwalrus/go-jitjson/.github/workflows/golang-ci.yml/badge.svg)](https://github.com/mcwalrus/go-jitjson/.github/workflows/golang-ci.yml)
[![CodeQL](https://github.com/mcwalrus/go-jitjson/.github/workflows/codeql.yml/badge.svg)](https://github.com/mcwalrus/go-jitjson/.github/workflows/codeql.yml)
[![Release](https://img.shields.io/github/release/mcwalrus/go-jitjson.svg)](https://github.com/mcwalrus/go-jitjson/releases/latest)

Changes to v1.3.0:

* Support custom parsers with JSONParser interface.
* Support opt-in encoding/json/v2 parser for experimental Go 1.25 builds.
* Support json.MarshalerTo and json.UnmarshalerFrom with JitJSON[T] type.
* New methods Parser, SetParser for JitJSON[T] and AnyJitJSON[T] types.
* Updated JitJSON[T] method Set to SetValue (breaking change).
* Provided json/v2 usage example.

Because of library usage is not yet adopted, I have included a breaking change on a minor version release.

Future versions will not include breaking changes.




- 🪶 Zero dependencies
- 🧩 Dynamic type parsing of JSON
- 🔗 Integrates with encoding/json interface types
- 💾 Reduce memory when working with large JSON datasets 
- 🚀 Improves performance by avoiding unnecessary parsing of JSON
- 🛠️ Customisable parsing


```Golang


type JITJSON struct {
    data   []byte
    keys   []string
    values []interface{}
}

// Marshal implementation stores json key value pairs to struct
func (jit *JITJSON) Marshal() error {
    jit.values = ...
    jit.keys = ...
    return nil
}



func (jit *JITJSON) GetInt(key string) (int, bool) {
    return nil
}

func (jit *JITJSON) GetInt64(key string) (int64, bool) {
    return nil
}

func (jit *JITJSON) GetStr(key string) (string, bool) {
    return nil
}

func (jit *JITJSON) GetInt64Array(key string) ([]int, bool) {
    return nil
}


func (jit *JITJSON) ToVesselEventFishing() *VesselEventFishing {
    return &VesselEventFishing{ /* ... */ }
}

func (jit *JITJSON) ToVesselEventMissingAIS() *VesselEventMissingAIS {
    return &VesselEventMissingAIS{ /* ... */ }
}

func (jit *JITJSON) ToVesselEventAnomolousMovement() *VesselEventAnomolousMovement {
    return &VesselEventAnomolousMovements{ /* ... */ }
}



```

---

## Marshal

```Go
package main

import (
    "fmt"
	"encoding/json"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email,omitempty"` // not included when empty
	Address string `json:"-"`               // This field will be ignored
}

func main() {
    person := Person{
        Name:    "Alice Johnson",
        Age:     30,
        Email:   "", 
        Address: "123 Main St",
    }

    jsonData, err := json.Marshal(person)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s\n", string(jsonData)) // Output: {"name":"Alice Johnson","age":30}
}
```

---

## Unmarshal

```Go
package main

import (
    "fmt"
	"encoding/json"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email,omitempty"` // not included when empty
	Address string `json:"-"`               // This field will be ignored
}

func main() {
    jsonData := []byte(`{
        "name":    "Augus Fitzgerald",
        "email":   "ang.gulls@hotmail.com",
        "address": "22 plaza palace, VIC 3001"
    }`)

    var person Person
    err := json.Unmarshal(jsonData, &person)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", person) // Output: {Augus Fitzgerald 0 ang.gulls@hotmail.com }
}
```

---

## Valid JSON

```Go
// valid
true 
1000.123 
"hello W0R1d" 
null 
[1, 2, 3] 
{"hello": "there!"} 

// invalid
My Guy!
{"hello":}
11000.2002.2100
```

## json.Marshal

```Go
https://pkg.go.dev/encoding/json#Marshaler
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
```

---

## Example

```Go
package main

import (
    "fmt"
	"encoding/json"
)

type Animal struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Type    string `json:"animal_type"`
}

func (p Animal) MarshalJSON() ([]byte, error) {
    if p.Type == "Dog" {
        return "I'm actually a cat. Meow!"
    }
    return fmt.Sprintf("Woof Woof Woof! said the %s", p.Type)
}

// verify interface:
var _ json.Marshaler = (nil)(*Animal)

func main() {
    animal1 := Person{
        Name:    "Lizzy Lizz",
        Age:     11,
        Type:    "Cat"
    }
    animal2 := Person{
        Name:    "Roger Rog",
        Age:     13,
        Type:    "Dog"
    }

    jsonData1, err := json.Marshal(animal1)
    if err != nil {
        panic(err)
    }

    jsonData2, err = json.Marshal(animal2)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s\n", string(jsonData1)) // Output: "Woof Woof Woof! said the Cat"
    fmt.Printf("%s\n", string(jsonData2)) // Output: "Dog is actually a cat. Meow!"
}
```

---

## json.Unmarshal

```Go
https://pkg.go.dev/encoding/json#Marshaler
type Marshaler interface {
	UnmarshalJSON(data []byte) error
}
```

---

## Example

```Go
package main

import (
	"encoding/json"
	"fmt"
)

type Animal struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Type string `json:"animal_type"`
}

func (a *Animal) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if name, ok := raw["name"].(string); ok {
		a.Name = name
	}
	if age, ok := raw["age"].(float64); ok {
		a.Age = int(age)
	}

	if typ, ok := raw["animal_type"].(string); ok {
		if typ == "Dog" {
			a.Type = "Actually a Cat!"
		} else if typ == "Cat" {
			a.Type = "Actually a Dog!"
		} else {
			a.Type = "I don't know!?!"
		}
	}
	return nil
}

// verify interface
var _ json.Unmarshaler = (*Animal)(nil)

func main() {
	jsonDog := `{"name":"Rover","age":5,"animal_type":"Dog"}`
	jsonCat := `{"name":"Whiskers","age":3,"animal_type":"Cat"}`
	jsonDuck := `{"name":"Donald","age":7,"animal_type":"Duck"}`

	var a1, a2, a3 Animal
	_ = json.Unmarshal([]byte(jsonDog), &a1)
	_ = json.Unmarshal([]byte(jsonCat), &a2)
	_ = json.Unmarshal([]byte(jsonDuck), &a3)

	fmt.Printf("Decoded Dog: %+v\n", a1)
	fmt.Printf("Decoded Cat: %+v\n", a2)
	fmt.Printf("Decoded Duck: %+v\n", a3)
}
```

