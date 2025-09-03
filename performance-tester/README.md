# JitJSON Performance Tester

A comprehensive performance testing tool for the [go-jitjson](https://github.com/mcwalrus/go-jitjson) library that automatically generates benchmarks comparing JitJSON against standard `encoding/json`.

## Overview

This tool takes a JSON file containing an array of objects and:

1. ✅ **Validates** that the JSON structure is an array of objects
2. 🏗️ **Generates** Go structs using the `gojson` command-line tool
3. 📊 **Creates** comprehensive benchmark tests comparing:
   - Standard `encoding/json` parsing
   - JitJSON with full parsing
   - JitJSON with partial parsing (configurable percentage)
   - Memory allocation comparisons
4. ✔️ **Verifies** correctness between both approaches
5. 🚀 **Runs** benchmarks automatically (optional)

## Prerequisites

1. **Go 1.18+** (required for generics support)
2. **gojson CLI tool**:
   ```bash
   go install github.com/ChimeraCoder/gojson/gojson@latest
   ```
3. **go-jitjson library** (automatically added to generated go.mod)

## Installation

```bash
# Clone or download the performance tester
cd go-jitjson/performance-tester
```

## Usage

### Basic Usage

```bash
go run main.go -json your-data.json
```

### Advanced Usage

```bash
go run main.go \
  -json data/users.json \
  -struct User \
  -package usertest \
  -output ./benchmarks \
  -parse-percent 0.2 \
  -run
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-json` | Path to JSON file (required) | - |
| `-output` | Output directory for generated files | `benchmark_output` |
| `-struct` | Name for the generated struct | `Item` |
| `-package` | Package name for generated code | `benchmarks` |
| `-parse-percent` | Default percentage for partial parsing (0.0-1.0) | `0.3` |
| `-run` | Run benchmarks after generation | `false` |

## JSON File Requirements

The JSON file **must** contain an array of objects. Examples:

✅ **Valid:**
```json
[
  {"name": "John", "age": 30, "city": "New York"},
  {"name": "Jane", "age": 25, "city": "Boston"}
]
```

❌ **Invalid:**
```json
{"users": [{"name": "John"}]}  // Not an array at root level
[1, 2, 3]                      // Array of primitives, not objects
```

## Generated Files

The tool creates three files in the output directory:

1. **`generated_structs.go`** - Go struct definitions created by gojson
2. **`benchmark_test.go`** - Comprehensive benchmark tests
3. **`go.mod`** - Go module with jitjson dependency

## Benchmark Types

### 1. Standard JSON (`BenchmarkStandardJSON`)
- Baseline performance using `encoding/json`
- Parses all data immediately

### 2. JitJSON Full (`BenchmarkJitJSON`)
- JitJSON with all data parsed
- Shows overhead when parsing everything

### 3. JitJSON Partial (`BenchmarkJitJSONPartial`)
- JitJSON with configurable parse percentage
- Demonstrates performance benefits of deferred parsing
- Use `PARSE_PERCENTAGE` environment variable to adjust

### 4. Memory Benchmarks
- `BenchmarkJitJSONMemory` - JitJSON allocation patterns
- `BenchmarkStandardJSONMemory` - Standard JSON allocation patterns

### 5. Correctness Test (`TestJitJSONCorrectness`)
- Verifies both approaches produce identical results

## Running Benchmarks

### After Generation

```bash
cd benchmark_output
go mod tidy
go test -bench=. -benchmem -v
```

### With Different Parse Percentages

```bash
# Test with 10% parsing
PARSE_PERCENTAGE=0.1 go test -bench=BenchmarkJitJSONPartial -benchmem

# Test with 50% parsing
PARSE_PERCENTAGE=0.5 go test -bench=BenchmarkJitJSONPartial -benchmem

# Test with 90% parsing
PARSE_PERCENTAGE=0.9 go test -bench=BenchmarkJitJSONPartial -benchmem
```

### Specific Benchmarks

```bash
# Only memory benchmarks
go test -bench=Memory -benchmem

# Only JitJSON benchmarks
go test -bench=JitJSON -benchmem

# Run correctness test
go test -run=TestJitJSONCorrectness -v
```

## Understanding Results

### Sample Output

```
BenchmarkStandardJSON-8              100    12345678 ns/op    1234567 B/op    1234 allocs/op
BenchmarkJitJSON-8                    80    15432100 ns/op    1345678 B/op    1345 allocs/op
BenchmarkJitJSONPartial-8            300     4123456 ns/op     456789 B/op     456 allocs/op
BenchmarkJitJSONMemory-8             500     2345678 ns/op     234567 B/op     234 allocs/op
BenchmarkStandardJSONMemory-8        200     5678901 ns/op     567890 B/op     567 allocs/op
```

### Key Metrics

- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)
- **allocs/op**: Number of allocations per operation (lower is better)

### When JitJSON Performs Better

- ✅ **Partial parsing**: When you only need to process a subset of data
- ✅ **Large datasets**: Memory efficiency becomes more important
- ✅ **Streaming scenarios**: When data arrives incrementally
- ✅ **Memory-constrained environments**: Reduced GC pressure

### When Standard JSON Performs Better

- ✅ **Full parsing**: When you need to process all data immediately
- ✅ **Small datasets**: Overhead may outweigh benefits
- ✅ **Simple structures**: Less complex unmarshaling logic

## Example Workflow

1. **Prepare your JSON data:**
   ```bash
   curl -o users.json "https://jsonplaceholder.typicode.com/users"
   ```

2. **Generate and run benchmarks:**
   ```bash
   go run main.go -json users.json -struct User -package usertest -run
   ```

3. **Analyze results and test different scenarios:**
   ```bash
   cd benchmark_output
   PARSE_PERCENTAGE=0.1 go test -bench=BenchmarkJitJSONPartial -benchmem
   PARSE_PERCENTAGE=0.5 go test -bench=BenchmarkJitJSONPartial -benchmem
   PARSE_PERCENTAGE=1.0 go test -bench=BenchmarkJitJSONPartial -benchmem
   ```

## Troubleshooting

### Common Issues

1. **"gojson command not found"**
   ```bash
   go install github.com/ChimeraCoder/gojson/gojson@latest
   ```

2. **"JSON must be an array of objects"**
   - Ensure your JSON file has an array at the root level
   - Each array element must be an object, not a primitive

3. **"benchmark execution failed"**
   - Run `go mod tidy` in the output directory
   - Ensure the JSON file path is accessible from the benchmark directory

### Performance Tips

- Use larger datasets (1000+ objects) for meaningful benchmarks
- Test with different parse percentages to find optimal performance
- Consider your specific use case (streaming vs batch processing)
- Monitor memory usage patterns in addition to execution time

## Contributing

This tool is designed to work with the go-jitjson library. For issues or improvements:

1. Test with various JSON structures
2. Verify benchmark accuracy
3. Report any compatibility issues

## License

This performance tester follows the same MIT license as the go-jitjson library.
