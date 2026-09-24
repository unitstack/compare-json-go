package comparejson_test

import (
	"testing"

	comparejson "github.com/unitstack/compare-json-go"
)

func TestPrimitiveValueComparison(t *testing.T) {
	t.Run("should return empty array for same values", func(t *testing.T) {
		cases := []struct{ base, contrast interface{} }{
			{"test", "test"},
			{float64(123), float64(123)},
			{true, true},
			{nil, nil},
		}
		for _, c := range cases {
			diffs := comparejson.CompareJSON(c.base, c.contrast, nil)
			if len(diffs) != 0 {
				t.Errorf("Expected 0 differences for %v == %v, got %d", c.base, c.contrast, len(diffs))
			}
		}
	})

	t.Run("should detect value changes", func(t *testing.T) {
		diffs := comparejson.CompareJSON("old", "new", nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{}, "", "both", comparejson.DiffTypeValueChanged)
	})

	t.Run("should detect type changes", func(t *testing.T) {
		diffs := comparejson.CompareJSON("string", float64(123), nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{}, "", "both", comparejson.DiffTypeTypeChanged)
	})
}

func TestObjectComparison(t *testing.T) {
	t.Run("should detect deleted keys", func(t *testing.T) {
		base := map[string]interface{}{"a": float64(1), "b": float64(2)}
		contrast := map[string]interface{}{"a": float64(1)}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"b"}, "b", "base", comparejson.DiffTypeDeleted)
	})

	t.Run("should detect added keys", func(t *testing.T) {
		base := map[string]interface{}{"a": float64(1)}
		contrast := map[string]interface{}{"a": float64(1), "b": float64(2)}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"b"}, "b", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should detect value changes in nested objects", func(t *testing.T) {
		base := map[string]interface{}{"nested": map[string]interface{}{"value": "old"}}
		contrast := map[string]interface{}{"nested": map[string]interface{}{"value": "new"}}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"nested", "value"}, "nested.value", "both", comparejson.DiffTypeValueChanged)
	})

	t.Run("should handle multiple differences in objects", func(t *testing.T) {
		base := map[string]interface{}{"a": float64(1), "b": "old", "c": true}
		contrast := map[string]interface{}{"a": float64(1), "b": "new", "d": false}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 3 {
			t.Fatalf("Expected 3 differences, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"b"}, "b", "both", comparejson.DiffTypeValueChanged)
		assertDiff(t, diffs[1], []string{"c"}, "c", "base", comparejson.DiffTypeDeleted)
		assertDiff(t, diffs[2], []string{"d"}, "d", "contrast", comparejson.DiffTypeAdded)
	})
}

func TestArrayComparison(t *testing.T) {
	t.Run("should detect deleted elements", func(t *testing.T) {
		base := []interface{}{float64(1), float64(2), float64(3)}
		contrast := []interface{}{float64(1), float64(2)}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[2]"}, "[2]", "base", comparejson.DiffTypeDeleted)
	})

	t.Run("should detect added elements", func(t *testing.T) {
		base := []interface{}{float64(1), float64(2)}
		contrast := []interface{}{float64(1), float64(2), float64(3)}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[2]"}, "[2]", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should detect value changes in array elements", func(t *testing.T) {
		base := []interface{}{float64(1), "old", true}
		contrast := []interface{}{float64(1), "new", true}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[1]"}, "[1]", "both", comparejson.DiffTypeValueChanged)
	})

	t.Run("should handle nested arrays", func(t *testing.T) {
		base := []interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), float64(4)}}
		contrast := []interface{}{[]interface{}{float64(1), float64(2)}, []interface{}{float64(3), "5"}}
		diffs := comparejson.CompareJSON(base, contrast, nil)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[1]", "[1]"}, "[1][1]", "both", comparejson.DiffTypeTypeChanged)
	})
}

func TestLCSArrayComparison(t *testing.T) {
	opts := &comparejson.CompareOptions{ArrayCompareMethod: comparejson.ArrayCompareLCS}

	t.Run("should correctly diff shifted arrays", func(t *testing.T) {
		base := []interface{}{"a", "b", "c"}
		contrast := []interface{}{"b", "c", "d"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 2 {
			t.Fatalf("Expected 2 differences, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[0]"}, "[0]", "base", comparejson.DiffTypeDeleted)
		assertDiff(t, diffs[1], []string{"[2]"}, "[2]", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should return empty for identical arrays", func(t *testing.T) {
		base := []interface{}{float64(1), float64(2), float64(3)}
		contrast := []interface{}{float64(1), float64(2), float64(3)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should handle empty arrays", func(t *testing.T) {
		diffs := comparejson.CompareJSON([]interface{}{}, []interface{}{}, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should handle no common elements", func(t *testing.T) {
		base := []interface{}{float64(1), float64(2), float64(3)}
		contrast := []interface{}{float64(4), float64(5), float64(6)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 6 {
			t.Fatalf("Expected 6 differences, got %d", len(diffs))
		}
	})

	t.Run("should handle objects in arrays", func(t *testing.T) {
		base := []interface{}{
			map[string]interface{}{"id": float64(1)},
			map[string]interface{}{"id": float64(2)},
			map[string]interface{}{"id": float64(3)},
		}
		contrast := []interface{}{
			map[string]interface{}{"id": float64(2)},
			map[string]interface{}{"id": float64(3)},
			map[string]interface{}{"id": float64(4)},
		}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 2 {
			t.Fatalf("Expected 2 differences, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[0]"}, "[0]", "base", comparejson.DiffTypeDeleted)
		assertDiff(t, diffs[1], []string{"[2]"}, "[2]", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should handle LCS where base has extra middle element", func(t *testing.T) {
		base := []interface{}{"a", "x", "b"}
		contrast := []interface{}{"a", "b"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[1]"}, "[1]", "base", comparejson.DiffTypeDeleted)
	})
}

func TestUnorderedArrayComparison(t *testing.T) {
	opts := &comparejson.CompareOptions{ArrayCompareMethod: comparejson.ArrayCompareUnordered}

	t.Run("should treat reordered arrays as equal", func(t *testing.T) {
		base := []interface{}{"a", "b", "c"}
		contrast := []interface{}{"c", "b", "a"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should detect added and deleted elements ignoring order", func(t *testing.T) {
		base := []interface{}{"a", "b", "c"}
		contrast := []interface{}{"a", "c", "d"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 2 {
			t.Fatalf("Expected 2 differences, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[1]"}, "[1]", "base", comparejson.DiffTypeDeleted)
		assertDiff(t, diffs[1], []string{"[2]"}, "[2]", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should handle duplicate values correctly", func(t *testing.T) {
		base := []interface{}{"a", "a", "b"}
		contrast := []interface{}{"a", "b", "b"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 2 {
			t.Fatalf("Expected 2 differences, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"[1]"}, "[1]", "base", comparejson.DiffTypeDeleted)
		assertDiff(t, diffs[1], []string{"[2]"}, "[2]", "contrast", comparejson.DiffTypeAdded)
	})

	t.Run("should handle objects ignoring order", func(t *testing.T) {
		base := []interface{}{
			map[string]interface{}{"id": float64(1)},
			map[string]interface{}{"id": float64(2)},
			map[string]interface{}{"id": float64(3)},
		}
		contrast := []interface{}{
			map[string]interface{}{"id": float64(3)},
			map[string]interface{}{"id": float64(1)},
			map[string]interface{}{"id": float64(2)},
		}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})
}

func TestValueCaseInsensitive(t *testing.T) {
	opts := &comparejson.CompareOptions{ValueCaseInsensitive: true}

	t.Run("should treat different-case strings as equal", func(t *testing.T) {
		diffs := comparejson.CompareJSON("Hello", "hello", opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should still detect different string values", func(t *testing.T) {
		diffs := comparejson.CompareJSON("hello", "world", opts)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{}, "", "both", comparejson.DiffTypeValueChanged)
	})

	t.Run("should ignore case in object values", func(t *testing.T) {
		base := map[string]interface{}{"name": "Alice"}
		contrast := map[string]interface{}{"name": "alice"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should ignore case in array string elements", func(t *testing.T) {
		base := []interface{}{"Hello", "World"}
		contrast := []interface{}{"hello", "world"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})
}

func TestKeyCaseInsensitive(t *testing.T) {
	opts := &comparejson.CompareOptions{KeyCaseInsensitive: true}

	t.Run("should match keys with different cases", func(t *testing.T) {
		base := map[string]interface{}{"Name": "Alice"}
		contrast := map[string]interface{}{"name": "Alice"}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should match nested keys with different cases", func(t *testing.T) {
		base := map[string]interface{}{"User": map[string]interface{}{"Name": "Alice"}}
		contrast := map[string]interface{}{"user": map[string]interface{}{"name": "Alice"}}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should detect deleted keys case-insensitively", func(t *testing.T) {
		base := map[string]interface{}{"a": float64(1), "B": float64(2)}
		contrast := map[string]interface{}{"A": float64(1)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"B"}, "B", "base", comparejson.DiffTypeDeleted)
	})

	t.Run("should detect added keys case-insensitively", func(t *testing.T) {
		base := map[string]interface{}{"a": float64(1)}
		contrast := map[string]interface{}{"A": float64(1), "b": float64(2)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{"b"}, "b", "contrast", comparejson.DiffTypeAdded)
	})
}

func TestNumericStringEqualsNumber(t *testing.T) {
	opts := &comparejson.CompareOptions{NumericStringEqualsNumber: true}

	t.Run("should treat numeric string as equal to number", func(t *testing.T) {
		diffs := comparejson.CompareJSON("123", float64(123), opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should treat number as equal to numeric string", func(t *testing.T) {
		diffs := comparejson.CompareJSON(float64(456), "456", opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should not treat non-numeric strings as equal to numbers", func(t *testing.T) {
		diffs := comparejson.CompareJSON("abc", float64(123), opts)
		if len(diffs) != 1 {
			t.Fatalf("Expected 1 difference, got %d", len(diffs))
		}
		assertDiff(t, diffs[0], []string{}, "", "both", comparejson.DiffTypeTypeChanged)
	})

	t.Run("should handle floating point numeric strings", func(t *testing.T) {
		diffs := comparejson.CompareJSON("3.14", float64(3.14), opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should handle negative numeric strings", func(t *testing.T) {
		diffs := comparejson.CompareJSON("-42", float64(-42), opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should work in nested objects", func(t *testing.T) {
		base := map[string]interface{}{"count": "100"}
		contrast := map[string]interface{}{"count": float64(100)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should work in arrays", func(t *testing.T) {
		base := []interface{}{"1", "2", "3"}
		contrast := []interface{}{float64(1), float64(2), float64(3)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})

	t.Run("should work with unordered array method", func(t *testing.T) {
		opts := &comparejson.CompareOptions{
			NumericStringEqualsNumber: true,
			ArrayCompareMethod:        comparejson.ArrayCompareUnordered,
		}
		base := []interface{}{"1", "2", "3"}
		contrast := []interface{}{float64(3), float64(2), float64(1)}
		diffs := comparejson.CompareJSON(base, contrast, opts)
		if len(diffs) != 0 {
			t.Fatalf("Expected 0 differences, got %d", len(diffs))
		}
	})
}

func assertDiff(t *testing.T, diff comparejson.Difference, expectedPath []string, expectedPathStr string, expectedBelongsTo string, expectedType comparejson.DiffType) {
	t.Helper()
	if diff.PathString != expectedPathStr {
		t.Errorf("Expected pathString %q, got %q", expectedPathStr, diff.PathString)
	}
	if string(diff.PathBelongsTo) != expectedBelongsTo {
		t.Errorf("Expected pathBelongsTo %q, got %q", expectedBelongsTo, diff.PathBelongsTo)
	}
	if diff.DiffType != expectedType {
		t.Errorf("Expected diffType %q, got %q", expectedType, diff.DiffType)
	}
	if len(diff.PathSegments) != len(expectedPath) {
		t.Errorf("Expected pathSegments length %d, got %d", len(expectedPath), len(diff.PathSegments))
		return
	}
	for i, seg := range diff.PathSegments {
		if seg != expectedPath[i] {
			t.Errorf("Expected pathSegments[%d] = %q, got %q", i, expectedPath[i], seg)
		}
	}
}
