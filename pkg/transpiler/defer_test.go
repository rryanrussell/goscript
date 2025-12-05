package transpiler

import (
	"testing"
)

func TestHasEffects(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		hasDefer   bool
		hasRecover bool
	}{
		{
			name: "pure expressions",
			code: `
				{
					x := 1
					y := 2
					z := x + y
				}
			`,
			hasDefer:   false,
			hasRecover: false,
		},
		{
			name: "defer statement",
			code: `
				{
					defer fmt.Println("done")
				}
			`,
			hasDefer:   true,
			hasRecover: false,
		},
		{
			name: "recover call",
			code: `
				{
					if r := recover(); r != nil {
						fmt.Println("recovered", r)
					}
				}
			`,
			hasDefer:   false,
			hasRecover: true,
		},
		{
			name: "defer and recover",
			code: `
				{
					defer func() {
						if r := recover(); r != nil {
							fmt.Println("recovered")
						}
					}()
				}
			`,
			hasDefer:   true,
			hasRecover: false, // recover is inside the deferred func, not the block itself
		},
		{
			name: "recover in block",
			code: `
				{
					recover()
				}
			`,
			hasDefer:   false,
			hasRecover: true,
		},
		{
			name: "nested function with defer (should be ignored)",
			code: `
				{
					f := func() {
						defer fmt.Println("nested")
					}
				}
			`,
			hasDefer:   false,
			hasRecover: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the code snippet into an AST
			body := parseBlock(t, tt.code)

			// Run hasEffects
			gotDefer, gotRecover := hasEffects(body)

			if gotDefer != tt.hasDefer {
				t.Errorf("hasEffects() defer = %v, want %v", gotDefer, tt.hasDefer)
			}
			if gotRecover != tt.hasRecover {
				t.Errorf("hasEffects() recover = %v, want %v", gotRecover, tt.hasRecover)
			}
		})
	}
}
