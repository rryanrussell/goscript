package transpiler

import (
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/rryanrussell/goscript/pkg/js"
)

func runRuntimeTest(t *testing.T, name, src, expectedOutput string) {
	t.Helper()

	// 1. Transpile Go source
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name+".go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	var buf bytes.Buffer
	ctx := TranspileContext{
		Fset:   fset,
		Writer: &buf,
	}
	fCtx := ModuleContext{TranspileContext: &ctx}
	fCtx.Module.Name = (*js.Ident)(&file.Name.Name)

	for _, decl := range file.Decls {
		Decl(&fCtx, decl)
	}
	fCtx.Module.Print(&fCtx)
	transpiledCode := buf.String()

	// 2. Prepare Runtime
	// Assuming we are running from pkg/transpiler, runtime is at ../runtime
	runtimeDir := "../runtime"
	runtimeFiles := []string{"runtime.js", "chan.js", "fmt.js", "strings.js", "reflect.js"}
	var runtimeCode bytes.Buffer

	// Add global definition if needed (though runtime.js handles it)
	runtimeCode.WriteString("if (typeof globalThis.runtime === 'undefined') { globalThis.runtime = {}; }\n")

	for _, rf := range runtimeFiles {
		path := filepath.Join(runtimeDir, rf)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read runtime file %s: %v", rf, err)
		}
		runtimeCode.Write(content)
		runtimeCode.WriteString("\n")
	}

	// 3. Create Driver Script
	var driverCode bytes.Buffer
	driverCode.Write(runtimeCode.Bytes())
	driverCode.WriteString("\n// Transpiled Code\n")
	driverCode.WriteString(transpiledCode)
	driverCode.WriteString("\n// Entry Point\n")
	driverCode.WriteString("if (typeof main === 'function') { main(); } else { console.error('main function not found'); }")

	// 4. Write to Temp File
	tmpDir := t.TempDir()
	jsPath := filepath.Join(tmpDir, "test.js")
	if err := os.WriteFile(jsPath, driverCode.Bytes(), 0644); err != nil {
		t.Fatalf("Failed to write test JS file: %v", err)
	}

	// 5. Run with Bun
	bunPath, err := exec.LookPath("bun")
	if err != nil {
		t.Skip("bun not found, skipping runtime test")
	}

	cmd := exec.Command(bunPath, jsPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bun execution failed: %v\nOutput:\n%s", err, output)
	}

	// 6. Verify Output
	// Normalize output (trim whitespace)
	got := string(bytes.TrimSpace(output))
	expected := string(bytes.TrimSpace([]byte(expectedOutput)))

	if got != expected {
		t.Errorf("Runtime output mismatch.\nExpected:\n%s\nGot:\n%s", expected, got)
	}
}

func TestRuntimeMakeCheck(t *testing.T) {
	src := `
package main

import (
	"fmt"
	"testing"
)

// Mock testing.T for standalone execution if needed, 
// but here we are just running main() which does the checks.
// Wait, the original make_check.go had Test functions AND a main function.
// The main function printed output. The Test functions used t.Errorf.
// Our runRuntimeTest checks stdout.
// So we should use the main function version that prints.

func main() {
	s := make([]int, 5)
	fmt.Println(len(s))
	s[0] = 10
	fmt.Println(s[0])

	m := make(map[string]int)
	m["a"] = 1
	fmt.Println(m["a"])
}
`
	expected := `5
10
1`
	runRuntimeTest(t, "MakeCheck", src, expected)
}

func TestInterfaceMethods(t *testing.T) {
	src := `
package main

import (
	"fmt"
	"reflect"
)

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
}

func (p *Person) Greet() string {
	return "Hello, " + p.Name
}

type Robot struct {
	ID int
}

func (r Robot) Greet() string {
	return "Beep Boop " + fmt.Sprintf("%d", r.ID)
}

func main() {
	p := &Person{Name: "Alice"}
	var g Greeter = p
	fmt.Println(g.Greet())

	r := Robot{ID: 42}
	g = &r
	fmt.Println(g.Greet())

	fmt.Println("Type of p:", reflect.TypeOf(p).Name())
	fmt.Println("Type of r:", reflect.TypeOf(r).Name())
}
`
	expected := `Hello, Alice
Beep Boop 42
Type of p: main.Person
Type of r: main.Robot`

	runRuntimeTest(t, "InterfaceMethods", src, expected)
}
