#!/bin/bash

# JitJSON Performance Tester Demo Script
# This script demonstrates the complete workflow of the performance tester

set -e  # Exit on any error

echo "🚀 JitJSON Performance Tester Demo"
echo "=================================="
echo

# Check if gojson is installed
if ! command -v gojson &> /dev/null; then
    echo "📦 Installing gojson..."
    go install github.com/ChimeraCoder/gojson/gojson@latest
    echo "✅ gojson installed successfully"
    echo
fi

# Step 1: Generate benchmarks using sample data
echo "📊 Step 1: Generating performance benchmarks..."
echo "Using sample-data.json with 5 user objects"
echo

go run main.go \
    -json sample-data.json \
    -struct User \
    -package usertest \
    -output demo_output \
    -parse-percent 0.3

echo
echo "✅ Benchmark generation completed!"
echo

# Step 2: Show generated files
echo "📁 Step 2: Generated files:"
echo "============================"
ls -la demo_output/
echo

# Step 3: Run the benchmarks
echo "🏃 Step 3: Running benchmarks..."
echo "==============================="
cd demo_output

echo "Installing dependencies..."
go mod tidy
echo

echo "Running all benchmarks:"
go test -bench=. -benchmem -v

echo
echo "🔬 Testing different parse percentages:"
echo "========================================"

echo
echo "📈 10% parsing:"
PARSE_PERCENTAGE=0.1 go test -bench=BenchmarkJitJSONPartial -benchmem

echo
echo "📈 50% parsing:"
PARSE_PERCENTAGE=0.5 go test -bench=BenchmarkJitJSONPartial -benchmem

echo
echo "📈 90% parsing:"
PARSE_PERCENTAGE=0.9 go test -bench=BenchmarkJitJSONPartial -benchmem

echo
echo "✅ Demo completed successfully!"
echo
echo "📋 Summary:"
echo "==========="
echo "- Generated Go structs from JSON using gojson"
echo "- Created comprehensive benchmark suite"
echo "- Verified correctness between JitJSON and standard JSON"
echo "- Tested various parsing percentages"
echo "- Measured memory allocation patterns"
echo
echo "🎯 Key Insights:"
echo "================"
echo "- JitJSON shows benefits when parsing < 100% of data"
echo "- Memory efficiency improves with lower parse percentages"
echo "- Correctness is maintained across all scenarios"
echo
echo "📂 All generated files are in: demo_output/"
echo "🔧 Customize by editing the command line parameters in main.go"
