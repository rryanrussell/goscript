package transpiler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// parseBlock parses a code snippet (which should be the body of a function)
// and returns the corresponding *ast.BlockStmt.
// The snippet is wrapped in `package main; func main() { ... }` before parsing.
func parseBlock(t *testing.T, code string) *ast.BlockStmt {
	t.Helper()
	src := "package main; func main() " + strings.TrimSpace(code)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}
	return f.Decls[0].(*ast.FuncDecl).Body
}
