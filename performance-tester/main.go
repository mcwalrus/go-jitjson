package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// Config holds the configuration for the performance tester
type Config struct {
	JSONFile     string
	OutputDir    string
	StructName   string
	PackageName  string
	ParsePercent float64
}

// validateJSONStructure checks if the JSON file contains an array of objects
func validateJSONStructure(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// First, check if it's valid JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Check if it's an array
	array, ok := jsonData.([]interface{})
	if !ok {
		return fmt.Errorf("JSON must be an array of objects, got %T", jsonData)
	}

	if len(array) == 0 {
		return fmt.Errorf("JSON array cannot be empty")
	}

	// Check if all elements are objects
	for i, item := range array {
		if _, ok := item.(map[string]interface{}); !ok {
			return fmt.Errorf("element at index %d is not an object, got %T", i, item)
		}
	}

	fmt.Printf("✓ JSON validation passed: found array with %d objects\n", len(array))
	return nil
}

// generateStructsWithGoJSON uses gojson command line tool to generate structs
func generateStructsWithGoJSON(config Config) error {
	// Check if gojson is installed
	if _, err := exec.LookPath("gojson"); err != nil {
		return fmt.Errorf("gojson command not found. Install it with: go install github.com/ChimeraCoder/gojson/gojson@latest")
	}

	outputFile := filepath.Join(config.OutputDir, "generated_structs.go")

	// Run gojson command
	cmd := exec.Command("gojson",
		"-name", config.StructName,
		"-pkg", config.PackageName,
		"-input", config.JSONFile,
		"-o", outputFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gojson failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ Generated structs in %s\n", outputFile)
	return nil
}

// BenchmarkTemplate contains the template for generating benchmark code
const BenchmarkTemplate = `package {{.PackageName}}

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

var testData []byte

func init() {
	var err error
	testData, err = os.ReadFile("{{.JSONFile}}")
	if err != nil {
		panic(err)
	}
}

// shouldParseIterator creates a function that returns true for a given percentage of calls
func shouldParseIterator(parsePercent float64) func() bool {
	var record float64 = 0
	return func() bool {
		record += parsePercent
		if record >= 1 {
			record = record - 1
			return true
		} else {
			return false
		}
	}
}

// BenchmarkStandardJSON benchmarks standard encoding/json parsing
func BenchmarkStandardJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data []{{.StructName}}
		err := json.Unmarshal(testData, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJitJSON benchmarks jitjson with full parsing
func BenchmarkJitJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data []*jitjson.JitJSON[{{.StructName}}]
		err := json.Unmarshal(testData, &data)
		if err != nil {
			b.Fatal(err)
		}

		// Parse all items
		for _, item := range data {
			_, err := item.Unmarshal()
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// BenchmarkJitJSONPartial benchmarks jitjson with partial parsing
func BenchmarkJitJSONPartial(b *testing.B) {
	parsePercent := {{.ParsePercent}}
	if envPercent := os.Getenv("PARSE_PERCENTAGE"); envPercent != "" {
		if p, err := strconv.ParseFloat(envPercent, 64); err == nil {
			parsePercent = p
		}
	}

	b.Logf("Parsing %.1f%% of data", parsePercent*100)

	for i := 0; i < b.N; i++ {
		shouldParse := shouldParseIterator(parsePercent)
		
		var data []*jitjson.JitJSON[{{.StructName}}]
		err := json.Unmarshal(testData, &data)
		if err != nil {
			b.Fatal(err)
		}

		// Parse only selected items
		for _, item := range data {
			if shouldParse() {
				_, err := item.Unmarshal()
				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}
}

// BenchmarkJitJSONMemory benchmarks memory allocation patterns
func BenchmarkJitJSONMemory(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		var data []*jitjson.JitJSON[{{.StructName}}]
		err := json.Unmarshal(testData, &data)
		if err != nil {
			b.Fatal(err)
		}
		// Don't parse - just measure allocation overhead
	}
}

// BenchmarkStandardJSONMemory benchmarks standard JSON memory allocation
func BenchmarkStandardJSONMemory(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		var data []{{.StructName}}
		err := json.Unmarshal(testData, &data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestJitJSONCorrectness ensures jitjson produces same results as standard JSON
func TestJitJSONCorrectness(t *testing.T) {
	// Standard JSON parsing
	var standardData []{{.StructName}}
	err := json.Unmarshal(testData, &standardData)
	if err != nil {
		t.Fatal(err)
	}

	// JitJSON parsing
	var jitData []*jitjson.JitJSON[{{.StructName}}]
	err = json.Unmarshal(testData, &jitData)
	if err != nil {
		t.Fatal(err)
	}

	if len(standardData) != len(jitData) {
		t.Fatalf("Length mismatch: standard=%d, jit=%d", len(standardData), len(jitData))
	}

	// Compare first few items to verify correctness
	maxCheck := 10
	if len(standardData) < maxCheck {
		maxCheck = len(standardData)
	}

	for i := 0; i < maxCheck; i++ {
		jitItem, err := jitData[i].Unmarshal()
		if err != nil {
			t.Fatal(err)
		}

		// Marshal both to JSON for comparison
		standardJSON, err := json.Marshal(standardData[i])
		if err != nil {
			t.Fatal(err)
		}

		jitJSON, err := json.Marshal(jitItem)
		if err != nil {
			t.Fatal(err)
		}

		if string(standardJSON) != string(jitJSON) {
			t.Fatalf("Item %d mismatch:\nStandard: %s\nJit: %s", i, standardJSON, jitJSON)
		}
	}

	t.Logf("✓ Correctness test passed for %d items", maxCheck)
}
`

// generateBenchmarkCode creates the benchmark test file
func generateBenchmarkCode(config Config) error {
	tmpl, err := template.New("benchmark").Parse(BenchmarkTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse benchmark template: %w", err)
	}

	outputFile := filepath.Join(config.OutputDir, "benchmark_test.go")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create benchmark file: %w", err)
	}
	defer file.Close()

	// Get absolute path for JSON file in template
	absJSONPath, err := filepath.Abs(config.JSONFile)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	templateData := struct {
		PackageName  string
		StructName   string
		JSONFile     string
		ParsePercent float64
	}{
		PackageName:  config.PackageName,
		StructName:   config.StructName,
		JSONFile:     absJSONPath,
		ParsePercent: config.ParsePercent,
	}

	if err := tmpl.Execute(file, templateData); err != nil {
		return fmt.Errorf("failed to execute benchmark template: %w", err)
	}

	fmt.Printf("✓ Generated benchmark code in %s\n", outputFile)
	return nil
}

// createGoMod creates a go.mod file for the benchmark project
func createGoMod(config Config) error {
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/mcwalrus/go-jitjson v1.0.0
)
`, config.PackageName)

	goModPath := filepath.Join(config.OutputDir, "go.mod")
	if err := ioutil.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	fmt.Printf("✓ Created go.mod in %s\n", goModPath)
	return nil
}

// runBenchmarks executes the generated benchmarks
func runBenchmarks(config Config) error {
	fmt.Println("\n🚀 Running benchmarks...")

	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "-v")
	cmd.Dir = config.OutputDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("benchmark execution failed: %w", err)
	}

	return nil
}

func main() {
	var config Config

	flag.StringVar(&config.JSONFile, "json", "", "Path to JSON file containing array of objects (required)")
	flag.StringVar(&config.OutputDir, "output", "benchmark_output", "Output directory for generated files")
	flag.StringVar(&config.StructName, "struct", "Item", "Name for the generated struct")
	flag.StringVar(&config.PackageName, "package", "benchmarks", "Package name for generated code")
	flag.Float64Var(&config.ParsePercent, "parse-percent", 0.3, "Default percentage of data to parse in partial benchmarks (0.0-1.0)")

	runBench := flag.Bool("run", false, "Run benchmarks after generation")
	flag.Parse()

	if config.JSONFile == "" {
		fmt.Println("Usage: go run main.go -json <path-to-json-file> [options]")
		fmt.Println("\nThis program generates performance benchmarks for jitjson library.")
		fmt.Println("It requires a JSON file containing an array of objects.")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  go run main.go -json data.json -struct User -package usertest -run")
		os.Exit(1)
	}

	if config.ParsePercent < 0 || config.ParsePercent > 1 {
		log.Fatal("parse-percent must be between 0.0 and 1.0")
	}

	fmt.Printf("🔧 JitJSON Performance Tester\n")
	fmt.Printf("JSON File: %s\n", config.JSONFile)
	fmt.Printf("Output Dir: %s\n", config.OutputDir)
	fmt.Printf("Struct Name: %s\n", config.StructName)
	fmt.Printf("Package: %s\n", config.PackageName)
	fmt.Printf("Parse Percent: %.1f%%\n\n", config.ParsePercent*100)

	// Step 1: Validate JSON structure
	if err := validateJSONStructure(config.JSONFile); err != nil {
		log.Fatal("JSON validation failed:", err)
	}

	// Step 2: Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	// Step 3: Generate structs using gojson
	if err := generateStructsWithGoJSON(config); err != nil {
		log.Fatal("Struct generation failed:", err)
	}

	// Step 4: Create go.mod
	if err := createGoMod(config); err != nil {
		log.Fatal("Go module creation failed:", err)
	}

	// Step 5: Generate benchmark code
	if err := generateBenchmarkCode(config); err != nil {
		log.Fatal("Benchmark generation failed:", err)
	}

	fmt.Println("\n✅ Performance test suite generated successfully!")
	fmt.Printf("📁 Files created in: %s\n", config.OutputDir)
	fmt.Println("   - generated_structs.go (struct definitions)")
	fmt.Println("   - benchmark_test.go (benchmark tests)")
	fmt.Println("   - go.mod (module definition)")

	if *runBench {
		if err := runBenchmarks(config); err != nil {
			log.Fatal("Benchmark execution failed:", err)
		}
	} else {
		fmt.Println("\n📊 To run benchmarks:")
		fmt.Printf("   cd %s\n", config.OutputDir)
		fmt.Println("   go mod tidy")
		fmt.Println("   go test -bench=. -benchmem")
		fmt.Println("\n🔬 To run with different parse percentages:")
		fmt.Println("   PARSE_PERCENTAGE=0.1 go test -bench=BenchmarkJitJSONPartial -benchmem")
	}
}
