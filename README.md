# compare-json-go

Go port of [`compare-json`](https://github.com/unitstack/compare-json) — **structured JSON comparison**: find what changed between two JSON values, with control over how keys, values, and arrays are matched.

> Online playground: **[comparejson.com](https://comparejson.com)**

## Why

Most JSON diff tools either output unstructured text or hide the parts that matter (where a key was added, whether a type changed, whether two arrays differ in order or in content). `compare-json-go` returns a structured list of differences — every entry carries a path, the side it belongs to (`base` / `contrast` / `both`), and the kind of change (`added`, `deleted`, `typeChanged`, `valueChanged`) — so you can render, filter, or program against it.

## Features

- **Deep comparison** of objects, arrays, and primitives.
- **Three array comparison strategies**: `byIndex` (default), `lcs` (minimal diff via Longest Common Subsequence), and `unordered` (multiset match).
- **Case-insensitive** key and/or value matching.
- **Numeric-string equality** — optionally treat `"1"` and `1` as equal.
- **Path tracking** with both segment-array and dot-notation forms.
- **Zero external dependencies** — the library and CLI use only the Go standard library.
- **CLI** with table or JSON output, reading from inline strings or files.

## Installation

```bash
go get github.com/unitstack/compare-json-go
```

## Quick Start

### Library

```go
package main

import (
	"encoding/json"
	"fmt"

	comparejson "github.com/unitstack/compare-json-go"
)

func main() {
	var base, contrast interface{}
	json.Unmarshal([]byte(`{"name":"Alice","age":30}`), &base)
	json.Unmarshal([]byte(`{"name":"Bob","age":"30","email":"bob@test.com"}`), &contrast)

	diffs := comparejson.CompareJSON(base, contrast, &comparejson.CompareOptions{
		ArrayCompareMethod: comparejson.ArrayCompareByIndex, // or ArrayCompareLCS, ArrayCompareUnordered
	})
	for _, d := range diffs {
		fmt.Println(d.PathString, d.PathBelongsTo, d.DiffType)
	}
	// name   both     valueChanged
	// age    both     typeChanged
	// email  contrast added
}
```

### CLI

```bash
go run ./cmd/compare-json base.json contrast.json
```

```
┌──────────────┬──────────────┐
│ Key          │ Change Type  │
├──────────────┼──────────────┤
│ (Base) a     │ valueChanged │
│ (Base) b     │ deleted      │
│ (Contrast) c │ added        │
└──────────────┴──────────────┘
```

```bash
# Array strategies and case-insensitive matching
compare-json -a lcs -k -v base.json contrast.json

# Machine-readable JSON output
compare-json --json-export base.json contrast.json

# Write the report to a file
compare-json -o diff.txt base.json contrast.json
```

Note: with the stdlib `flag` package, flags must come **before** the positional arguments.

## Options

| Flag | Description |
|------|-------------|
| `-a, --array-compare-method` | Array comparison strategy: `byIndex` (default), `lcs`, `unordered` |
| `-k, --key-case-insensitive` | Compare object keys case-insensitively |
| `-v, --value-case-insensitive` | Compare string values case-insensitively |
| `--numeric-string-equals-number` | Treat numeric strings as equal to numbers |
| `-j, --json-export` | Output the differences as JSON |
| `-o, --output <file>` | Write output to a file instead of stdout |

## Differences from the TypeScript reference

This port intentionally uses native Go semantics instead of replicating JavaScript
quirks. See [COMPATIBILITY.md](./COMPATIBILITY.md) for the full policy. In short:

- Numeric strings are parsed with `strconv.ParseFloat` (no JS `Number()` quirks like
  `"" → 0`, whitespace trimming, or hex).
- Object keys are compared in **lexicographic order** (Go maps have no insertion order),
  so difference entries may be ordered differently than the TS output; content is the same.
- Key existence checks look at own keys only (the TS reference's `in`-operator
  prototype-chain behavior is a bug and is not reproduced).
- JSON is decoded into `float64`, so integers beyond 2^53 lose precision (same as JS).

## Development

```bash
# build
go build ./...

# run unit tests and CLI end-to-end tests
go test ./...
```

The repo layout:

```
├── compare.go, object.go, array.go, types.go, utils.go   # library (package comparejson)
├── cmd/compare-json/                                     # CLI
│   ├── main.go
│   └── main_test.go                                      # CLI e2e tests (builds the real binary)
└── example/                                              # runnable demo
```

## License

MIT
