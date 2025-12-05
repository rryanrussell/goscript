package transpiler

import (
	"testing"
)

func TestStdlib(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "FmtPrintf",
			src: `
package main

import "fmt"

func main() {
	fmt.Printf("Hello %s\n", "World")
}
`,
			expected: `function main() {
    runtime.fmt.printf("Hello %s\n", "World")
}
`,
		},
		{
			name: "FmtSprintf",
			src: `
package main

import "fmt"

func main() {
	s := fmt.Sprintf("Value: %d", 42)
}
`,
			expected: `function main() {
    let s = runtime.fmt.sprintf("Value: %d", 42)
}
`,
		},
		{
			name: "StringsContains",
			src: `
package main

import "strings"

func main() {
	b := strings.Contains("hello", "ell")
}
`,
			expected: `function main() {
    let b = runtime.strings.Contains("hello", "ell")
}
`,
		},
		{
			name: "StringsSplit",
			src: `
package main

import "strings"

func main() {
	parts := strings.Split("a,b,c", ",")
}
`,
			expected: `function main() {
    let parts = runtime.strings.Split("a,b,c", ",")
}
`,
		},
		{
			name: "ErrorsNew",
			src: `
package main

import "errors"

func main() {
	err := errors.New("failed")
}
`,
			expected: `function main() {
    let err = runtime.errors.New("failed")
}
`,
		},
		{
			name: "StrconvAtoi",
			src: `
package main

import "strconv"

func main() {
	i, err := strconv.Atoi("123")
}
`,
			expected: `function main() {
    let [ i, err ] = runtime.strconv.Atoi("123")
}
`,
		},
		{
			name: "MathAbs",
			src: `
package main

import "math"

func main() {
	f := math.Abs(-1.5)
}
`,
			expected: `function main() {
    let f = runtime.math.Abs(-1.5)
}
`,
		},
		{
			name: "TimeNow",
			src: `
package main

import "time"

func main() {
	t := time.Now()
}
`,
			expected: `function main() {
    let t = runtime.time.Now()
}
`,
		},
		{
			name: "SortInts",
			src: `
package main

import "sort"

func main() {
	sort.Ints([]int{3, 1, 2})
}
`,
			expected: `function main() {
    runtime.sort.Ints([ 3, 1, 2 ])
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTranspilerTest(t, tt.name, tt.src, tt.expected)
		})
	}
}
