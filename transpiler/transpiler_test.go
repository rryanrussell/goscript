package transpiler

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
    "github.com/rryanrussell/goscript/js"
)

func TestIfElse(t *testing.T) {
	src := `
package main

import "fmt"

func main() {
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	} else {
		fmt.Println("x is not greater than 5")
	}

	if x > 20 {
		fmt.Println("x is greater than 20")
	} else if x > 15 {
		fmt.Println("x is greater than 15")
	} else {
		fmt.Println("x is small")
	}
}
`
	expected := `function main() {
    let x = 10
    if (x>5) {
        fmt.Println("x is greater than 5")
    }
 else  {
        fmt.Println("x is not greater than 5")
    }

    if (x>20) {
        fmt.Println("x is greater than 20")
    }
 else if (x>15) {
        fmt.Println("x is greater than 15")
    }
 else  {
        fmt.Println("x is small")
    }

}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	var buf bytes.Buffer
	ctx := TranspileContext{
		Fset:   fset,
		Writer: &buf,
	}
	fCtx := ModuleContext{TranspileContext: &ctx}

    // We need to initialize the context properly, similar to main.go
    fCtx.Module.Name = (*js.Ident)(&file.Name.Name)

	for _, decl := range file.Decls {
		Decl(&fCtx, decl)
	}

	fCtx.Module.Print(&fCtx)

    // Normalize newlines for comparison
    got := buf.String()
    if got != expected {
        t.Errorf("Expected:\n%s\nGot:\n%s", expected, got)
    }
}
