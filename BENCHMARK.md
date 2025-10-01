## Benchmarks

The library provides comprehensive benchmarks for both `encoding/json` (v1) and `encoding/json/v2` implementations. Note jitjson benchmarks perform relative to the size and volume of data the library is applied to. Benefits across applications will be varied with most significant effects in reduced pressure in garbage collection cycles considering different ways data can avoid being parsed.

### Run Benchmarks

An assumed data structure is used for the benefits of jitjson benchmarks:

```bash
export GOEXPERIMENT=jsonv2 
# Run all benchmarks (encoding/json v1)
go test -bench=. -benchmem
# Run all benchmarks (encoding/json v2) - requires Go 1.25+
go test -bench=. -benchmem
```

### Marshaling Benchmarks Worst-Case

These benchmarks perform worst-case analysis by marshaling all objects in the dataset:

```bash
go test -bench='^BenchmarkMarshalWorstCase$' -benchmem
go test -bench='^BenchmarkMarshalWorstCaseV2$' -benchmem
```

### Unmarshaling Benchmarks Worst-Case

These benchmarks perform worst-case analysis by unmarshaling all objects in the dataset:

```bash
go test -bench='^BenchmarkUnmarshalWorstCase$' -benchmem
go test -bench='^BenchmarkUnmarshalWorstCaseV2$' -benchmem
```

### Partial Dataset Parsing Benchmarks

These benchmarks test parsing at a configurable percentage of the total objects available:

```bash
PARSE_PERCENTAGE=0.5 go test -bench='^BenchmarkParsePercentage$' -benchmem
PARSE_PERCENTAGE=0.5 go test -bench='^BenchmarkParsePercentageV2$' -benchmem
```

### Nested JitJSON Structs Benchmarks

Test performance with nested jitjson structures:

```bash
go test -bench='^BenchmarkNestedParse$' -benchmem
go test -bench='^BenchmarkNestedParseV2$' -benchmem
```

### Custom Benchmarking

For your own accurate performance evaluation, consider creating benchmarks with your own data:

```go
// use json or jsonv2
import "encoding/json"

type MyStruct struct { /* ... */ } 

func BenchmarkMyData(b *testing.B) {

    // Read json from file
    data, err := os.ReadFile(os.Getenv("JSON_FILE"))
    if err != nil {
        b.Fatalf("failed to retrieve json-data: %v", err)
    }
    
    // Parse only what you need
    b.Run("jitjson", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var jit []*jitjson.JitJSON[MyStruct]
            json.Unmarshal(data, &jit)
            for _, item := range jit[:len(jit)/2] {
                item.Unmarshal()
            }
        }
    })
    
    // Compare against the standard library
    b.Run("encoding-json", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var result []MyStruct
            json.Unmarshal(data, &result)
        }
    })
}
```