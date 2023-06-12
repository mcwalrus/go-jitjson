# JIT-JSON

Go JitJSON provides an implementation to encode or decode json encodings or values in _just-in-time_ fashion within runtime.

## Usage

Go JitJSON is a lightweight wrapper over [encoding/json](https://pkg.go.dev/encoding/json).

I recommend to only use the module in cases where you may not need to encode or decode all encodings / values in memory.

## Reference

See module documentation at: [pkg.go.dev/github.com/mcwalrus/go-jitjson](https://pkg.go.dev/github.com/mcwalrus/go-jitjson).

At an overview, the type `jitjson.JitJSON` can contain a json encoding or value. Type can perform encoding of a value or decoding of an encoding in just-in-time fashion by use of methods `Marshal()` and `Unmarshal()`. Type also implements [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) and [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) interfaces where type can be used directly with module methods `json.Marshal()` and `json.Unmarshal()`.

## Examples

Encode json:
```Go
// store value.
var value = `"json encoded"`
jit, err := jitjson.NewJitJSON[string](value)
if err != nil {
    panic(err)
}

// just-in-time 'encode'.
var data []byte
data, err = jit.Marshal()
if err != nil {
    panic(err)
}

// TODO: use encoding 'data'.
```

Decode json:
```Go
// store encoding.
var data = []byte(`"json encoded"`)
jit, err := jitjson.NewJitJSON[string](data)
if err != nil {
    panic(err)
}

// just-in-time 'decode'.
var value string
value, err = jit.Unmarshal()
if err != nil {
    panic(err)
}

// TODO: use value 'value'.
```

`encoding/json` module usage:
```Go
var (
    data = []byte(`"..."`)
    jit  jitjson.NewJitJSON[string]
)

// store encoding of 'jit'.
err := json.Unmarshal(data, &jit)
if err != nil {
    panic(err)
}

// return encoding of 'jit'.
var data []byte
data, err = json.Marshal(jit)
if err != nil {
    panic(err)
}
```

## Advanced usage

If your retrieving an encoding or value from multiple places, jitjson.JitJSON caches encoding or values once returned:

Encode json:
```Go
...

// just-in-time 'encode'.
data, err := jit.Marshal()
if err != nil {
    panic(err)
}

// first attempt to marshal caches encoding from jit - returns immediately.
data, err = jit.Marshal()
if err != nil {
    panic(err)
}
```

Decode json:
```Go
...

// just-in-time 'decode'.
var value string
value, err = jit.Unmarshal()
if err != nil {
    panic(err)
}

// first attempt to unmarshal caches value from jit - returns immediately.
value, err = jit.Unmarshal()
if err != nil {
    panic(err)
}
```

In some use cases, you may want to encode or decode indexes/elements of arrays/slices, or properties/values of objects/maps in just-in-time fashion. See `jit_json_chuck_test.go` for examples of such cases.

## Contribute

Feedback is always appreciated for my projects. Please submit or reach out to me through the issue tracker, or personally at collierwm@outlook.com.

Cheers!
