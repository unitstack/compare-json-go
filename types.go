package comparejson

type DiffType string

const (
	DiffTypeAdded        DiffType = "added"
	DiffTypeDeleted      DiffType = "deleted"
	DiffTypeTypeChanged  DiffType = "typeChanged"
	DiffTypeValueChanged DiffType = "valueChanged"
)

type PathBelongsTo string

const (
	PathBelongsToBase     PathBelongsTo = "base"
	PathBelongsToContrast PathBelongsTo = "contrast"
	PathBelongsToBoth     PathBelongsTo = "both"
)

type ArrayCompareMethod string

const (
	ArrayCompareByIndex   ArrayCompareMethod = "byIndex"
	ArrayCompareLCS       ArrayCompareMethod = "lcs"
	ArrayCompareUnordered ArrayCompareMethod = "unordered"
)

type CompareOptions struct {
	ArrayCompareMethod        ArrayCompareMethod
	KeyCaseInsensitive        bool
	ValueCaseInsensitive      bool
	NumericStringEqualsNumber bool
}

type Difference struct {
	PathSegments  []string      `json:"pathSegments"`
	PathString    string        `json:"pathString"`
	PathBelongsTo PathBelongsTo `json:"pathBelongsTo"`
	DiffType      DiffType      `json:"diffType"`
}
