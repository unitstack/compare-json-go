package comparejson

import (
	"reflect"
	"strconv"
	"strings"
)

func CompareJSON(base, contrast interface{}, options *CompareOptions) []Difference {
	if options == nil {
		options = &CompareOptions{ArrayCompareMethod: ArrayCompareByIndex}
	}
	differences := compareValue(base, contrast, []string{}, options)
	if differences == nil {
		return []Difference{}
	}
	return differences
}

func compareValue(base, contrast interface{}, pathSegments []string, options *CompareOptions) []Difference {
	baseType := getValueType(base)
	contrastType := getValueType(contrast)

	if baseType == "object" && contrastType == "object" {
		return compareObject(base.(map[string]interface{}), contrast.(map[string]interface{}), pathSegments, options)
	}

	if baseType == "array" && contrastType == "array" {
		baseArr := toSlice(base)
		contrastArr := toSlice(contrast)
		switch options.ArrayCompareMethod {
		case ArrayCompareLCS:
			return compareArrayLCS(baseArr, contrastArr, pathSegments, options)
		case ArrayCompareUnordered:
			return compareArrayUnordered(baseArr, contrastArr, pathSegments, options)
		default:
			return compareArray(baseArr, contrastArr, pathSegments, options)
		}
	}

	if baseType != contrastType {
		if options.NumericStringEqualsNumber {
			if isNumericStringMatch(base, contrast, baseType, contrastType) {
				return []Difference{}
			}
		}
		return []Difference{{
			PathSegments:  pathSegments,
			PathString:    pathSegmentsToString(pathSegments),
			PathBelongsTo: PathBelongsToBoth,
			DiffType:      DiffTypeTypeChanged,
		}}
	}

	if !valuesEqual(base, contrast, baseType, options) {
		return []Difference{{
			PathSegments:  pathSegments,
			PathString:    pathSegmentsToString(pathSegments),
			PathBelongsTo: PathBelongsToBoth,
			DiffType:      DiffTypeValueChanged,
		}}
	}

	return []Difference{}
}

func toSlice(v interface{}) []interface{} {
	rv := reflect.ValueOf(v)
	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result
}

func isNumericStringMatch(base, contrast interface{}, baseType, contrastType string) bool {
	if baseType == "string" && contrastType == "number" {
		if num, err := strconv.ParseFloat(base.(string), 64); err == nil {
			return num == contrast.(float64)
		}
	}
	if baseType == "number" && contrastType == "string" {
		if num, err := strconv.ParseFloat(contrast.(string), 64); err == nil {
			return base.(float64) == num
		}
	}
	return false
}

func valuesEqual(base, contrast interface{}, valueType string, options *CompareOptions) bool {
	if valueType == "string" && options.ValueCaseInsensitive {
		return strings.ToLower(base.(string)) == strings.ToLower(contrast.(string))
	}
	return reflect.DeepEqual(base, contrast)
}
