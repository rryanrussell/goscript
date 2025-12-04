# goscript Architecture

Deep dive into how goscript transforms Go code to JavaScript.

> **Note:** This document describes the **current implementation** (circa 2022). For future expansion plans including stdlib support, interfaces, goroutines, and more, see [ROADMAP.md](ROADMAP.md).

## Overview

goscript is a **source-to-source transpiler** that converts a subset of Go to readable JavaScript. Unlike full compilers (GopherJS), it targets a minimal language subset with clean output.

**Current status:** ~2k LOC transpiler supporting ~15% of Go (structs, functions, basic control flow)

## Pipeline

```
Go Source
    ↓
[go/parser] Parse Go into Go AST
    ↓
[transpiler] Transform Go AST → JS AST
    ↓
[link] Resolve external symbols
    ↓
[js] Pretty-print JavaScript
    ↓
JavaScript Output
```

## Components

### 1. Go AST (stdlib)

Uses Go's built-in `go/parser` and `go/ast` packages:

```go
fset := token.NewFileSet()
pkgs, _ := parser.ParseDir(fset, dir, nil, parser.ParseComments)
```

This gives us a fully-parsed Go AST to work with.

### 2. JS AST (js/ package)

Custom JavaScript AST representation in `js/` (~728 LOC):

**Key types:**
- `Expr` - Expression nodes (identifiers, literals, calls, etc.)
- `Stmt` - Statement nodes (assignments, returns, loops, etc.)
- `Decl` - Declaration nodes (functions, classes, variables)
- `Module` - Top-level module container

**Design choice:** Custom AST (vs reusing existing JS parser) allows:
- Precise control over output formatting
- Go-centric operations (e.g., `append()` → spread operator)
- Minimal dependencies

**Example nodes:**

```go
// js/expr.go
type Binary struct {
    X  Expr
    Op Token  // +, -, ??, etc.
    Y  Expr
}

type Call struct {
    Fun  Expr
    Args []Expr
}

// js/stmt.go
type Return struct {
    X Expr
}

type Assign struct {
    Left  []Expr
    Op    Token  // =, :=
    Right []Expr
}
```

### 3. Transpiler (transpiler/ package)

Core transformation logic (~12KB):

**Main entry point:**
```go
func Transpile(dir string, out io.Writer) error {
    // 1. Parse Go
    pkgs := parser.ParseDir(dir)

    // 2. Transform each declaration
    for _, file := range pkg.Files {
        for _, decl := range file.Decls {
            Decl(&context, decl)  // Go AST → JS AST
        }
    }

    // 3. Link externals
    link.LinkExternals(&context, &module)

    // 4. Print JavaScript
    module.Print(&context)
}
```

**Key transformations:**

#### Function Declarations
```go
// Go AST: *ast.FuncDecl
func Foo(x int) int {
    return x + 1
}

// → JS AST: js.FuncDecl
function Foo(x) {
    return x + 1;
}
```

#### Struct Types
```go
// Go AST: *ast.StructType
type Vec2 struct {
    X float64
    Y float64
}

// → JS AST: js.ClassDecl
class Vec2 {
    constructor() {
        this.X = 0;
        this.Y = 0;
    }
}
```

#### Special Go Functions
```go
// len() → .length ?? 0
len(arr)  →  arr?.length ?? 0

// append() → spread operator
append(arr, x)  →  [...(arr ?? []), x]

// Go slicing → JS slice()
arr[1:3]  →  arr.slice(1, 3)
```

#### Pointer Semantics
```go
// & and * are tracked but mostly erased (JS has reference semantics)
func modify(ptr *Vec2) {
    ptr.X = 10
}

// →
function modify(ptr) {
    ptr.X = 10;  // Just regular reference
}
```

### 4. Link (link/ package)

Resolves external symbols defined in `.extern.go` files:

**Example extern definition:**
```go
// dom.extern.go
type HTMLCanvasElement struct{}
func (c HTMLCanvasElement) GetContext(ctx string) WebGLRenderingContext

var Document DocumentObject
```

**Symbol table:**
```go
type Symbol struct {
    Name     string
    Type     string
    External bool  // Maps to browser API
}

// Built-in symbols (link/builtin.go)
func GetBuiltInSymbols() map[string]Symbol {
    return map[string]Symbol{
        "Document":           {External: true},
        "HTMLCanvasElement":  {External: true},
        // ... etc
    }
}
```

**Linking process:**
1. Scan `.extern.go` files for type/function signatures
2. Build symbol table mapping Go names → JS names
3. Replace references with JS equivalents:
   - `Document` → `document`
   - `GetElementsByTagName` → `getElementsByTagName`
   - `WebGLRenderingContext` → (native browser type)

### 5. Pretty Printing (js/*.go)

Each JS AST node implements `Print()`:

```go
// js/func.go
func (f *FuncDecl) Print(ctx *Context) {
    fmt.Fprintf(ctx.Writer, "function %s(", f.Name)
    // ... print params
    fmt.Fprintf(ctx.Writer, ") {\n")
    // ... print body
    fmt.Fprintf(ctx.Writer, "}\n")
}
```

**Formatting features:**
- Proper indentation
- Semicolon insertion
- Whitespace control
- No source maps (limitation)

## Type System

**Runtime Type System (Since 2025):**

The system uses a `runtime` namespace (injected or defined globally) to handle type operations.

*   **Struct Metadata:**
    *   Structs are transpiled to Classes.
    *   Constructors call `runtime.SetType(this, "pkg.StructName")` to attach type metadata.
*   **Type Switches:**
    *   Transpiled to `if-else` chains.
    *   Checks `typeof x` for primitives.
    *   Checks `x?.__goType` for structs.
*   **Type Assertions:**
    *   Transpiled to IIFEs using `runtime.panic` on failure.
    *   Comma-ok form returns `[value, boolean]` tuple.

**Go types → JS equivalents:**

| Go Type | JavaScript |
|---------|------------|
| `int`, `float64` | `number` |
| `string` | `string` |
| `bool` | `boolean` |
| `[]T` | `Array` |
| `struct` | `class` (with `__goType` metadata) |
| `map[K]V` | `Object` or `Map` |
| `func` | `function` |
| `interface{}` | `any` |

**Limitations:**
- No type checking (assumes valid Go)
- No generic types
- No interface dispatch
- No reflection

## Compilation Context

Tracks state during transpilation:

```go
type TranspileContext struct {
    Fset   *token.FileSet     // Go source positions
    Writer io.Writer           // Output destination
}

type ModuleContext struct {
    *TranspileContext
    Module    js.Module         // Accumulated JS AST
    Semantics map[js.Node]Semantics  // Node metadata
}

type Semantics struct {
    PassByReference bool  // Track pointer types
}
```

## Example Transformation

**Input Go:**
```go
package app

type State struct {
    Count int
}

func increment(s *State) {
    s.Count = s.Count + 1
}

func Main() {
    state := State{Count: 0}
    increment(&state)
}
```

**Go AST (simplified):**
```
*ast.File
  Decls:
    *ast.GenDecl (type State)
      Specs:
        *ast.TypeSpec
          Type: *ast.StructType
            Fields: [Count int]

    *ast.FuncDecl (increment)
      Type: *ast.FuncType
        Params: [s *State]
      Body: *ast.BlockStmt
        List: [s.Count = s.Count + 1]

    *ast.FuncDecl (Main)
      Body: *ast.BlockStmt
        List:
          - state := State{Count: 0}
          - increment(&state)
```

**JS AST (simplified):**
```
js.Module
  Decls:
    js.ClassDecl (State)
      Fields: [Count: 0]

    js.FuncDecl (increment)
      Params: [s]
      Body:
        js.Assign: s.Count = s.Count + 1

    js.FuncDecl (Main)
      Body:
        js.VarDecl: state = new State()
        js.Assign: state.Count = 0
        js.Call: increment(state)
```

**Output JavaScript:**
```javascript
class State {
    constructor() {
        runtime.SetType(this, "app.State");
        this.Count = 0;
    }
}

function increment(s) {
    s.Count = s.Count + 1;
}

function Main() {
    const state = new State();
    state.Count = 0;
    increment(state);
}
```

## Design Decisions

### Why Custom JS AST?

**Pros:**
- Full control over output formatting
- Can add Go-specific transformations easily
- No dependency on JS parser implementation
- Easier to debug and understand

**Cons:**
- More code to maintain
- Need to handle JS semantics correctly
- Can't round-trip JS code

**Verdict:** Worth it for educational clarity and output control.

### Why No Goroutines?

JavaScript has no native threading. Options:
1. Transpile to async/await (complex)
2. Use Web Workers (heavyweight)
3. Don't support (chosen)

**Verdict:** Out of scope for minimal transpiler.

### Why .extern.go Pattern?

External definitions keep type signatures in Go while mapping to browser APIs:

```go
// In Go:
gl.ClearColor(0, 0, 0, 1)

// Maps to JS:
gl.clearColor(0, 0, 0, 1)
```

**Benefits:**
- Type-safe in Go land
- Clear boundary between Go and JS
- Easy to add new bindings

**Limitation:** Manual work to define APIs.

### Why No Standard Library?

Supporting `fmt`, `io`, `net`, etc. would require:
- Reimplementing huge portions in JS
- Or bundling entire Go runtime

**Verdict:** Use browser APIs directly instead.

## Performance

**Transpilation time:** Fast (<100ms for small projects)
- Single-pass transformation
- No optimization passes
- Minimal analysis

**Runtime performance:** Depends on output quality
- No runtime overhead (unlike full compilers)
- Generated code is straightforward JS
- Browser optimizes it like any JS

**Output size:**
- ~1:1.5 Go LOC : JS LOC ratio
- No large runtime bundle
- Human-readable (not minified)

## Future Improvements

> **See [ROADMAP.md](ROADMAP.md) for the comprehensive expansion plan.**

- Else clauses, switch statements, break/continue
- Slice expressions, better make/copy/append
- Improved type assertions and map operations
- strings, fmt, errors, strconv, math, time, sort packages
- Interfaces & method dispatch
- Defer & panic/recover
- Goroutines → async/await, Channels → Promise queues
- Reflection (limited), generics

## Lessons Learned

1. **Minimal subsets are useful** - Don't need full language for browser UI
2. **Clean output matters** - Readable JS helps debugging
3. **Custom AST is powerful** - Full control over transformations
4. **WASM is usually better** - But experiments teach you a lot
5. **Runtime Abstraction** - Keeping runtime logic separate from AST generation simplifies maintenance.

## References

**Similar projects:**
- [GopherJS](https://github.com/gopherjs/gopherjs) - Full Go compiler
- [Joy](https://mat.tm/joy) - Similar minimal approach (abandoned)
- [TinyGo](https://tinygo.org) - Small WASM builds

**Go compiler resources:**
- [go/ast documentation](https://pkg.go.dev/go/ast)
- [How Go builds work](https://golang.org/doc/install/source)

**JavaScript AST formats:**
- [ESTree spec](https://github.com/estree/estree) - JS AST standard
- [Babel AST](https://babeljs.io/docs/en/babel-types) - Similar approach

---

This architecture favors **simplicity and readability** over completeness. It's a teaching tool as much as a working transpiler.
