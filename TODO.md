
Benchmark results

Worst case is 70% more allocations with 60% slower performance. Bytes used scale linearly at 10% per order of magnitude.

The breakpoint on performance ns/op for parsing arrays is around 40% of the data being parsed. This scales linearly but is dependent on the size of the data.

The examples




Worst case:
# Benchmark Performance Comparison: JitJSON vs Stdlib

BenchmarkFullParse/JitJSON/Small-12         	   16695	     70135 ns/op	    9024 B/op	     239 allocs/op
BenchmarkFullParse/Stdlib/Small-12          	   26930	     44222 ns/op	    6888 B/op	     143 allocs/op
BenchmarkFullParse/JitJSON/Medium-12        	    1719	    818816 ns/op	   87984 B/op	    2402 allocs/op
BenchmarkFullParse/Stdlib/Medium-12         	    2596	    511558 ns/op	   62952 B/op	    1406 allocs/op
BenchmarkFullParse/JitJSON/Large-12         	     171	   8162065 ns/op	  873846 B/op	   24005 allocs/op
BenchmarkFullParse/Stdlib/Large-12          	     238	   6580736 ns/op	  564203 B/op	   14009 allocs/op

# Define the benchmark results
benchmark_results = {
    'Small': {
        'JitJSON': {'ns/op': 70135, 'B/op': 9024, 'allocs/op': 239},
        'Stdlib': {'ns/op': 44222, 'B/op': 6888, 'allocs/op': 143}
    },
    'Medium': {
        'JitJSON': {'ns/op': 818816, 'B/op': 87984, 'allocs/op': 2402},
        'Stdlib': {'ns/op': 511558, 'B/op': 62952, 'allocs/op': 1406}
    },
    'Large': {
        'JitJSON': {'ns/op': 8162065, 'B/op': 873846, 'allocs/op': 24005},
        'Stdlib': {'ns/op': 6580736, 'B/op': 564203, 'allocs/op': 14009}
    }
}

# Metrics to compare
metrics = ['ns/op', 'B/op', 'allocs/op']

# Function to calculate and display ratios
def calculate_ratios(results):
    print(f"{'Size':<10} | {'Metric':<10} | {'JitJSON':>10} | {'Stdlib':>10} | {'Ratio (J/S)':>12}")
    print("-" * 60)
    for size, data in results.items():
        for metric in metrics:
            jitjson = data['JitJSON'][metric]
            stdlib = data['Stdlib'][metric]
            ratio = jitjson / stdlib if stdlib != 0 else float('inf')
            print(f"{size:<10} | {metric:<10} | {jitjson:>10} | {stdlib:>10} | {ratio:>12.2f}")
    print()

# Execute the comparison
calculate_ratios(benchmark_results)
```







To do:

- Add benchmarks
- Improve documentation
- Add benchmarks
- Add streaming benchmarks
- Add examples


## Benchmarks

Could you provide a set of benchmarks focused on presenting the performance on:

* When only a few elements / object fields are to be parsed from large json files.

* Overhead of the worst case when parsed from large json files.

* Benchmarking described in a dynamic example of using AnyJitJSON.

```json
{
    "nil": null,
    "bool": true,
    "number": 123.45,
    "string": "Hello, World!",
    "slice": [1, "two", false, null],
    "object": {
        "nestedString": "Nested",
        "nestedNumber": 456,
            "deepNestedString": "Deep Nested",
            "deepNestedNumber": 678
        }
    }
}
```





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
    jit = jitjson.New[int](1)
    jit = jitjson.New[float64](2.0)
    jit = jitjson.New[string]("another type!")

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
