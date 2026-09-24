package comparejson

import (
	"fmt"
	"reflect"
	"regexp"
)

func getValueType(v interface{}) string {
	if v == nil {
		return "null"
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return "string"
	case reflect.Float64, reflect.Int, reflect.Int64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "unknown"
	}
}

func pathSegmentsToString(segments []string) string {
	if len(segments) == 0 {
		return ""
	}
	result := segments[0]
	indexPattern := regexp.MustCompile(`^\[\d+\]$`)
	for i := 1; i < len(segments); i++ {
		seg := segments[i]
		if indexPattern.MatchString(seg) {
			result += seg
		} else {
			result += "." + seg
		}
	}
	return result
}

func formatIndexStr(index int) string {
	return fmt.Sprintf("[%d]", index)
}
