package tests

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dop251/goja"
	"github.com/rryanrussell/goscript/js"
	"github.com/rryanrussell/goscript/link"
	x "github.com/rryanrussell/goscript/transpiler"
)

// RunTranspiled executes the given Go source code using the Goscript transpiler
// and runs the resulting JavaScript in a Goja runtime.
// It returns the Goja runtime instance and any error encountered.
func RunTranspiled(source string) (*goja.Runtime, error) {
	// 1. Transpile Go to JS
	jsCode, err := TranspileString(source)
	if err != nil {
		return nil, fmt.Errorf("transpilation failed: %w", err)
	}

	// 2. Run JS in Goja
	vm := goja.New()

	// Add console.log support
	vm.Set("console", map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			var args []interface{}
			for _, arg := range call.Arguments {
				args = append(args, arg.Export())
			}
			fmt.Println(args...)
			return goja.Undefined()
		},
	})

	_, err = vm.RunString(jsCode)
	if err != nil {
		fmt.Println("JS Code:\n", jsCode)
		return nil, fmt.Errorf("runtime execution failed: %w", err)
	}

	return vm, nil
}

// TranspileString transpiles a single Go source file (as a string) to JavaScript.
// This is a simplified version of main.Transpile adapted for strings.
func TranspileString(source string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.go", source, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	ctx := x.TranspileContext{
		Fset:   fset,
		Writer: &out,
	}

	fCtx := x.ModuleContext{TranspileContext: &ctx}
	extern := link.GetBuiltInSymbols()

	// Symbolize
	symbols := link.Symbolize(&ctx, file)
	for k, s := range symbols {
		extern[k] = s
	}

	fCtx.Module.Name = (*js.Ident)(&file.Name.Name)

	for _, n := range file.Decls {
		x.Decl(&fCtx, n)
	}

	fCtx.Exit()

	linkCtx := link.LinkContext{
		Extern:    extern,
		Semantics: fCtx.Semantics,
	}

	// TODO: Handle warnings if needed
	// if len(linkCtx.Warns) > 0 {
	// 	fmt.Printf("linkCtx.Warns: %+v\n", linkCtx.Warns)
	// }

	link.LinkExternals(&linkCtx, &fCtx.Module)

	// Capture the runtime includes
	var runtimeOut bytes.Buffer
	link.Include(&runtimeOut)

	// Write runtime first
	out.WriteString(runtimeOut.String())

	// Then write the transpiled module
	fCtx.Module.Print(&fCtx)

	return out.String(), nil
}
