package comparejson

import (
	"sort"
	"strings"
)

func compareObject(base, contrast map[string]interface{}, pathSegments []string, options *CompareOptions) []Difference {
	var differences []Difference
	matchedContrastKeys := make(map[string]bool)

	var contrastKeyMap map[string][]string
	if options.KeyCaseInsensitive {
		contrastKeyMap = buildCaseInsensitiveMap(contrast)
	}

	for _, key := range sortedKeys(base) {
		baseValue := base[key]
		var matchedKey string
		var exists bool

		if options.KeyCaseInsensitive {
			matchedKey, exists = findUnmatchedKey(contrastKeyMap, strings.ToLower(key), matchedContrastKeys)
		} else {
			matchedKey = key
			_, exists = contrast[key]
		}

		newPath := append(append([]string{}, pathSegments...), key)

		if !exists {
			differences = append(differences, Difference{
				PathSegments:  newPath,
				PathString:    pathSegmentsToString(newPath),
				PathBelongsTo: PathBelongsToBase,
				DiffType:      DiffTypeDeleted,
			})
		} else {
			matchedContrastKeys[matchedKey] = true
			contrastValue := contrast[matchedKey]
			differences = append(differences, compareValue(baseValue, contrastValue, newPath, options)...)
		}
	}

	for _, key := range sortedKeys(contrast) {
		if matchedContrastKeys[key] {
			continue
		}
		hasMatch := false
		if options.KeyCaseInsensitive {
			for sk := range base {
				if strings.ToLower(sk) == strings.ToLower(key) {
					hasMatch = true
					break
				}
			}
		} else {
			_, hasMatch = base[key]
		}
		if !hasMatch {
			newPath := append(append([]string{}, pathSegments...), key)
			differences = append(differences, Difference{
				PathSegments:  newPath,
				PathString:    pathSegmentsToString(newPath),
				PathBelongsTo: PathBelongsToContrast,
				DiffType:      DiffTypeAdded,
			})
		}
	}

	return differences
}

func buildCaseInsensitiveMap(obj map[string]interface{}) map[string][]string {
	m := make(map[string][]string)
	// Iterate in sorted key order so that candidate lists are deterministic;
	// Go map iteration order is random and would otherwise make
	// case-insensitive key pairing non-deterministic.
	for _, k := range sortedKeys(obj) {
		lower := strings.ToLower(k)
		m[lower] = append(m[lower], k)
	}
	return m
}

func findUnmatchedKey(m map[string][]string, lowerKey string, matched map[string]bool) (string, bool) {
	candidates, ok := m[lowerKey]
	if !ok {
		return "", false
	}
	for _, k := range candidates {
		if !matched[k] {
			return k, true
		}
	}
	return "", false
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
