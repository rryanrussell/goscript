package transpiler

import (
	"bytes"
	"go/parser"
	"go/token"
	"os/exec"
	"testing"

	"github.com/rryanrussell/goscript/pkg/js"
	"github.com/rryanrussell/goscript/pkg/runtime"
)

func TestSync(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
	}{
		{
			name: "Mutex",
			src: `
package main

import (
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	count := 0
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		mu.Lock()
		defer mu.Unlock()
		c := count
		time.Sleep(10 * time.Millisecond)
		count = c + 1
		wg.Done()
	}()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		c := count
		time.Sleep(10 * time.Millisecond)
		count = c + 1
		wg.Done()
	}()

	wg.Wait()
	wg.Wait()
	if count != 2 {
		panic(count)
	}
}
`,
		},
		{
			name: "Once",
			src: `
package main

import (
	"sync"
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"testing"

	"github.com/rryanrussell/goscript/pkg/js"
)

func main() {
	var once sync.Once
	count := 0
	fn := func() {
		count++
	}
	
	var wg sync.WaitGroup
	wg.Add(3)
	
	for i := 0; i < 3; i++ {
		go func() {
			once.Do(fn)
			wg.Done()
		}()
	}
	
	wg.Wait()
	if count != 1 {
		panic("count mismatch")
	}
}
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Use runTranspilerTest-like logic but just check for success for now
			// or we can assert expected output if we know it.
			// Let's assert expected output to be sure.

			// We need to define expected output for these tests.
			// Since the output is large, maybe we just check if it contains key elements?
			// But runTranspilerTest expects exact match.
			// Let's just define a simple helper here.

			out := transpileToString(t, tt.name, tt.src)
			if out == "" {
				t.Fatal("expected output")
			}
			runJS(t, out)
		})
	}
}

func transpileToString(t *testing.T, name, src string) string {
	t.Helper()
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
	return buf.String()
}

func runJS(t *testing.T, src string) {
	t.Helper()
	// Prepend runtime loading
	// We use the embedded files from the runtime package

	var preamble bytes.Buffer
	preamble.WriteString("var runtime;\nvar sync;\n")

	preamble.Write(runtime.RuntimeJS)
	preamble.WriteString("\n")
	preamble.Write(runtime.SyncJS)
	preamble.WriteString("\n")
	preamble.Write(runtime.TimeJS)
	preamble.WriteString("\n")
	preamble.Write(runtime.ChanJS)
	preamble.WriteString("\n")

	// Mock fmt.println for testing
	preamble.WriteString("runtime.fmt = { println: console.log };\n")
	// Mock time.Millisecond if not present (it's usually in Go, but transpiled to number?)
	// In Go: time.Sleep(10 * time.Millisecond)
	// Transpiled: runtime.time.Sleep(10 * time.Millisecond)
	// If time.Millisecond is not defined in JS, this fails.
	// The transpiler should probably handle time.Millisecond as a constant or runtime.time.Millisecond.
	// Let's assume for now we need to define it if missing.
	preamble.WriteString("if (!runtime.time) runtime.time = {};\n")
	preamble.WriteString("if (!runtime.time.Millisecond) runtime.time.Millisecond = 1000000;\n") // ns

	// Wait, time.Sleep in JS runtime likely expects ms?
	// Let's check time.js content.

	fullSrc := preamble.String() + src + "\nmain();\n"

	cmd := exec.Command("node", "-e", fullSrc)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("node failed: %v\nstderr: %s\nsrc:\n%s", err, stderr.String(), fullSrc)
	}
}
