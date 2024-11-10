input = """
BenchmarkNestedParseWorstCase/JitJSON/Small-12         	   35372	     34224 ns/op	    4696 B/op	     121 allocs/op
BenchmarkNestedParseWorstCase/Stdlib/Small-12          	  119982	     10292 ns/op	    1568 B/op	      40 allocs/op
BenchmarkNestedParseWorstCase/JitJSON/Medium-12        	     645	   1854512 ns/op	  147912 B/op	    1392 allocs/op
BenchmarkNestedParseWorstCase/Stdlib/Medium-12         	   15350	     76064 ns/op	   10928 B/op	     226 allocs/op
BenchmarkNestedParseWorstCase/JitJSON/Large-12         	       6	 170335688 ns/op	11966714 B/op	   17152 allocs/op
BenchmarkNestedParseWorstCase/Stdlib/Large-12          	    1497	    770239 ns/op	  100816 B/op	    2033 allocs/op
"""

def parse_benchmark_input(input_str):
    # Initialize results dictionary
    benchmark_results = {}
    benchmark_name = None
    
    # Process each line
    for line in input_str.strip().split('\n'):

        # Skip non-benchmark lines
        if not line.startswith('Benchmark'):
            continue
            
        # Parse benchmark line
        # Format: BenchmarkFullParse/Parser/Size-N  iterations  time  memory  allocs
        parts = line.split()
        if not parts or len(parts) != 8:
            continue
            
        # Extract benchmark name from first part
        full_name = parts[0]
        name_parts = full_name.split('/')
        
        if len(name_parts) != 3:
            continue
            
        benchmark_name = name_parts[0].replace('Benchmark', '')
        parser = name_parts[1]
        size = name_parts[2].split('-')[0]
        
        # Parse metrics (removing units)
        iterations = int(parts[1])
        ns_op = int(parts[2])
        b_op = int(parts[4])
        allocs_op = int(parts[6])
        
        # Build nested structure
        if size not in benchmark_results:
            benchmark_results[size] = {}
        
        benchmark_results[size][parser] = {
            'iterations': iterations,
            'ns/op': ns_op,
            'B/op': b_op,
            'allocs/op': allocs_op
        }
    
    return benchmark_name, benchmark_results

# Metrics to compare
metrics = ['ns/op', 'B/op', 'allocs/op'] # 'iterations'

# Function to calculate and display ratios
def calculate_ratios(benchmark_name, results):
    print(f"Benchmark: {benchmark_name}")
    print(f"{'Size':<10} | {'Metric':<10} | {'JitJSON':>10} | {'Stdlib':>10} | {'Ratio (J/S)':>12}")
    print("-" * 60)
    
    for size, data in results.items():
        for metric in metrics:
            jitjson = data['JitJSON'][metric]
            stdlib = data['Stdlib'][metric]
            ratio = jitjson / stdlib if stdlib != 0 else float('inf')
            print(f"{size:<10} | {metric:<10} | {jitjson:>10} | {stdlib:>10} | {ratio:>12.2f}")
    print()

# Parse input and execute comparison
benchmark_name, benchmark_results = parse_benchmark_input(input)
calculate_ratios(benchmark_name, benchmark_results)