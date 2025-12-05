package transpiler

import (
	"testing"
)

func TestHasAsync(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "no async operations",
			code: `
				{
					x := 1
					y := 2
					fmt.Println(x + y)
				}
			`,
			expected: false,
		},
		{
			name: "send statement",
			code: `
				{
					ch <- 1
				}
			`,
			expected: true,
		},
		{
			name: "select statement",
			code: `
				{
					select {
					case <-ch:
					default:
					}
				}
			`,
			expected: true,
		},
		{
			name: "unary receive expression",
			code: `
				{
					x := <-ch
				}
			`,
			expected: true,
		},
		{
			name: "receive in expression statement",
			code: `
				{
					<-ch
				}
			`,
			expected: true,
		},
		{
			name: "nested function with async (should be ignored)",
			code: `
				{
					f := func() {
						<-ch
					}
				}
			`,
			expected: false,
		},
		{
			name: "async in if block",
			code: `
				{
					if true {
						<-ch
					}
				}
			`,
			expected: true,
		},
		{
			name: "async in for loop",
			code: `
				{
					for {
						ch <- 1
					}
				}
			`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code snippet into an AST
			body := parseBlock(t, tt.code)

			// Run hasAsync
			result := hasAsync(body)

			if result != tt.expected {
				t.Errorf("hasAsync() = %v, want %v", result, tt.expected)
			}
		})
	}
}
