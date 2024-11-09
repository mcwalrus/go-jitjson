input = """
BenchmarkPartialParse/JitJSON/Small-12         	   24825	     42233 ns/op	    4088 B/op	     108 allocs/op
BenchmarkPartialParse/Stdlib/Small-12          	   26164	     42821 ns/op	    6888 B/op	     143 allocs/op
BenchmarkPartialParse/JitJSON/Medium-12        	    3014	    385069 ns/op	   38552 B/op	    1029 allocs/op
BenchmarkPartialParse/Stdlib/Medium-12         	    2988	    408855 ns/op	   62952 B/op	    1406 allocs/op
BenchmarkPartialParse/JitJSON/Large-12         	     306	   3853306 ns/op	  379372 B/op	   10212 allocs/op
BenchmarkPartialParse/Stdlib/Large-12          	     295	   4090471 ns/op	  564204 B/op	   14009 allocs/op
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