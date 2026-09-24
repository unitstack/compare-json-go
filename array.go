package comparejson

func compareArray(base, contrast []interface{}, pathSegments []string, options *CompareOptions) []Difference {
	var differences []Difference
	minLen := len(base)
	if len(contrast) < minLen {
		minLen = len(contrast)
	}

	for i := 0; i < minLen; i++ {
		newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
		differences = append(differences, compareValue(base[i], contrast[i], newPath, options)...)
	}

	for i := minLen; i < len(base); i++ {
		newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
		differences = append(differences, Difference{
			PathSegments:  newPath,
			PathString:    pathSegmentsToString(newPath),
			PathBelongsTo: PathBelongsToBase,
			DiffType:      DiffTypeDeleted,
		})
	}

	for i := minLen; i < len(contrast); i++ {
		newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
		differences = append(differences, Difference{
			PathSegments:  newPath,
			PathString:    pathSegmentsToString(newPath),
			PathBelongsTo: PathBelongsToContrast,
			DiffType:      DiffTypeAdded,
		})
	}

	return differences
}

func compareArrayLCS(base, contrast []interface{}, pathSegments []string, options *CompareOptions) []Difference {
	var differences []Difference
	baseMatched, contrastMatched := computeLCS(base, contrast, options)

	baseInLCS := make(map[int]bool)
	contrastInLCS := make(map[int]bool)
	for _, idx := range baseMatched {
		baseInLCS[idx] = true
	}
	for _, idx := range contrastMatched {
		contrastInLCS[idx] = true
	}

	for i := range base {
		if !baseInLCS[i] {
			newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
			differences = append(differences, Difference{
				PathSegments:  newPath,
				PathString:    pathSegmentsToString(newPath),
				PathBelongsTo: PathBelongsToBase,
				DiffType:      DiffTypeDeleted,
			})
		}
	}

	for i := range contrast {
		if !contrastInLCS[i] {
			newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
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

func compareArrayUnordered(base, contrast []interface{}, pathSegments []string, options *CompareOptions) []Difference {
	var differences []Difference
	contrastMatched := make([]bool, len(contrast))

	for i, bv := range base {
		found := false
		for j, cv := range contrast {
			if !contrastMatched[j] && deepEqual(bv, cv, options) {
				contrastMatched[j] = true
				found = true
				break
			}
		}
		if !found {
			newPath := append(append([]string{}, pathSegments...), formatIndexStr(i))
			differences = append(differences, Difference{
				PathSegments:  newPath,
				PathString:    pathSegmentsToString(newPath),
				PathBelongsTo: PathBelongsToBase,
				DiffType:      DiffTypeDeleted,
			})
		}
	}

	for j := range contrast {
		if !contrastMatched[j] {
			newPath := append(append([]string{}, pathSegments...), formatIndexStr(j))
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

func deepEqual(a, b interface{}, options *CompareOptions) bool {
	diffs := compareValue(a, b, []string{}, options)
	return len(diffs) == 0
}

func computeLCS(base, contrast []interface{}, options *CompareOptions) ([]int, []int) {
	m, n := len(base), len(contrast)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if deepEqual(base[i-1], contrast[j-1], options) {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	var baseMatchedIdx, contrastMatchedIdx []int
	i, j := m, n
	for i > 0 && j > 0 {
		if deepEqual(base[i-1], contrast[j-1], options) {
			baseMatchedIdx = append([]int{i - 1}, baseMatchedIdx...)
			contrastMatchedIdx = append([]int{j - 1}, contrastMatchedIdx...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return baseMatchedIdx, contrastMatchedIdx
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
