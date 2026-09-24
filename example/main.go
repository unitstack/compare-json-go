package main

import (
	"fmt"

	comparejson "github.com/unitstack/compare-json-go"
)

func main() {
	base := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"tags": []interface{}{"go", "json"},
	}

	contrast := map[string]interface{}{
		"name": "John",
		"age":  float64(31),
		"tags": []interface{}{"go", "yaml"},
	}

	options := &comparejson.CompareOptions{
		ArrayCompareMethod: comparejson.ArrayCompareByIndex,
	}

	diffs := comparejson.CompareJSON(base, contrast, options)

	fmt.Printf("Found %d differences:\n", len(diffs))
	for _, d := range diffs {
		fmt.Printf("  %s (%s): %s\n", d.PathString, d.PathBelongsTo, d.DiffType)
	}
}
