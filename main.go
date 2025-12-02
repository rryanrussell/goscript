package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/rryanrussell/goscript/js"
	"github.com/rryanrussell/goscript/link"
	x "github.com/rryanrussell/goscript/transpiler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goscript <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]
	if err := Transpile(dir, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Transpile(dir string, out io.Writer) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)

	if err != nil {
		return err
	}

	if len(pkgs) != 1 {
		return fmt.Errorf("unsupported package count: expected %d, found %d", 1, len(pkgs))
	}

	ctx := x.TranspileContext{
		Fset:   fset,
		Writer: out,
	}

	fCtx := x.ModuleContext{TranspileContext: &ctx}

	extern := link.GetBuiltInSymbols()

	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			symbols := link.Symbolize(&ctx, file)
			for k, s := range symbols {
				extern[k] = s
			}

			if strings.HasSuffix(filename, ".extern.go") {
				continue
			}

			fCtx.Module.Name = (*js.Ident)(&file.Name.Name)

			for _, n := range file.Decls {
				x.Decl(&fCtx, n)
			}
		}
	}

	fCtx.Exit()

	linkCtx := link.LinkContext{
		Extern:    extern,
		Semantics: fCtx.Semantics,
	}

	if len(linkCtx.Warns) > 0 {
		fmt.Printf("linkCtx.Warns: %+v\n", linkCtx.Warns)
	}

	link.LinkExternals(&linkCtx, &fCtx.Module)

	link.Include(ctx.Writer)

	fCtx.Module.Print(&fCtx)

	return nil
}
