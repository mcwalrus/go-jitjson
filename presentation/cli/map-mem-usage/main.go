package main

import (
	"flag"
	"fmt"
	"runtime"
)

func printMemoryUsageKB(label string, m *runtime.MemStats) {
	fmt.Printf("%s:\n", label)
	fmt.Printf("  Allocated Memory: %.2f KB\n", float64(m.Alloc)/1024)
	fmt.Printf("  Total Allocated: %.2f KB\n", float64(m.TotalAlloc)/1024)
	fmt.Printf("  Heap Memory: %.2f KB\n", float64(m.HeapAlloc)/1024)
	fmt.Printf("  System Memory: %.2f KB\n", float64(m.HeapSys)/1024)
	fmt.Printf("  Malloc Calls: %d\n", m.Mallocs)
	fmt.Printf("  Free Calls: %d\n", m.Frees)
	fmt.Println()
}

func printMemoryUsageMB(label string, m *runtime.MemStats) {
	fmt.Printf("%s:\n", label)
	fmt.Printf("  Allocated Memory: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("  Total Allocated: %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("  Heap Memory: %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("  System Memory: %.2f MB\n", float64(m.HeapSys)/1024/1024)
	fmt.Printf("  Malloc Calls: %d\n", m.Mallocs)
	fmt.Printf("  Free Calls: %d\n", m.Frees)
	fmt.Println()
}

func main() {
	// Parse command line arguments
	numMaps := flag.Int("maps", 1, "number of maps to create")
	numElements := flag.Int("elements", 1000, "number of key-value pairs per map")
	capacity := flag.Int("capacity", 0, "capacity of the map")
	flag.Parse()

	var m runtime.MemStats

	// 1. Memory before any maps are initialized
	runtime.GC()
	runtime.ReadMemStats(&m)
	printMemoryUsageKB("Before map initialization", &m)

	// 2. Create and initialize maps
	maps := make([]map[int]any, *numMaps)
	for i := range maps {
		if *capacity > 0 {
			maps[i] = make(map[int]any, *capacity)
		} else {
			maps[i] = make(map[int]any)
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&m)
	printMemoryUsageKB("After map initialization", &m)

	// 3. Populate maps with elements
	emptyStruct := "<string>"
	for i := range maps {
		for j := 0; j < *numElements; j++ {
			maps[i][j] = emptyStruct
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&m)
	printMemoryUsageMB("After populating maps", &m)

	// Keep reference to prevent GC
	fmt.Printf("Created %d maps with %d elements each\n", len(maps), *numElements)
}
