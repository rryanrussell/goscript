package transpiler

import "testing"

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "AppendMultiple",
			src: `
package main

func main() {
	arr := []int{1}
	arr = append(arr, 2, 3)
}
`,
			expected: `function main() {
    let arr = [ 1 ]
    arr = [ ...(arr??[  ]), 2, 3 ]
}
`,
		},
		{
			name: "AppendSpread",
			src: `
package main

func main() {
	arr := []int{1}
	other := []int{2, 3}
	arr = append(arr, other...)
}
`,
			expected: `function main() {
    let arr = [ 1 ]
    let other = [ 2, 3 ]
    arr = [ ...(arr??[  ]), ...(other??[  ]) ]
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
