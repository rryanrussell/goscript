# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

goscript is a minimal Go-to-JavaScript transpiler (~2k LOC) that converts a Go subset to clean, readable JavaScript. This is an educational/experimental project, not intended for production use.

**Key characteristics:**
- Transpiles Go structs, functions, and basic control flow to JavaScript
- Generates human-readable JS output with minimal runtime overhead
- Uses custom JavaScript AST for precise output control
- Supports direct browser API bindings via `.extern.go` files
- Does NOT support: goroutines, channels, interfaces, reflection, most stdlib

## Development Commands

```bash
# Transpile a Go directory to JavaScript (outputs to stdout)
go run main.go <input-directory>

# Example: Transpile the webgl-app example
go run main.go examples/webgl-app > examples/webgl-app/output.js

# Run with custom output
go run main.go ./path/to/go/code > output.js

# Use as library in Go code
# See main.go Transpile() function - takes directory path and io.Writer
```

## Architecture

The transpiler follows a 4-stage pipeline:

### 1. Go AST (stdlib: go/parser, go/ast)
Parses Go source files into Go's built-in AST representation.

### 2. Transpiler (transpiler/ package)
Core transformation logic that converts Go AST nodes to JavaScript AST nodes:
- **transpiler.go**: Main Transpile() entry point and orchestration
- **instance.go**: Transpiler instance that tracks state during transformation
- **context.go**: Compilation context (file set, writer, semantics tracking)

Key transformations:
- Go functions → JS functions
- Go structs → JS classes with constructors
- `len()` → `.length ?? 0`
- `append()` → spread operator `[...(arr ?? []), x]`
- Go slicing `arr[1:3]` → `arr.slice(1, 3)`
- Pointer semantics tracked but mostly erased (JS uses references)

### 3. JavaScript AST (js/ package)
Custom JS AST implementation (~728 LOC, 17 files) providing precise output control:
- **ast.go**: Core AST interface definitions
- **expr.go**: Expression nodes (binary ops, calls, literals)
- **stmt.go**: Statement nodes (assignments, returns, loops)
- **decl.go**: Declaration nodes (functions, variables)
- **class.go**: Class declarations (from Go structs)
- **func.go**: Function declarations
- **block.go**: Block statements
- **call.go**: Function calls
- **lit.go**: Literals
- **loop.go**: Loop constructs
- **module.go**: Top-level module container
- **name.go**: Name/identifier handling
- **opr.go**: Operators
- **token.go**: Token types
- **var.go**: Variable declarations
- **print.go**: Pretty-printer for generating JavaScript output
- **guest.go**: External/guest language interop

### 4. Link (link/ package)
Symbol resolution for external browser APIs:
- **symbol.go**: Symbol table and linking logic
- **builtin.go**: Built-in browser API mappings

External bindings defined in `.extern.go` files map Go type signatures to browser APIs.

## Code Organization

```
goscript/
├── main.go              # CLI entry point, Transpile() function
├── js/                  # JavaScript AST (17 files)
│   └── *.go            # AST nodes and pretty-printing
├── transpiler/          # Go AST → JS AST transformation
│   ├── transpiler.go   # Main transformation logic
│   ├── instance.go     # Transpiler instance state
│   └── context.go      # Compilation context
├── link/                # Symbol resolution
│   ├── symbol.go       # Symbol table
│   └── builtin.go      # Browser API mappings
└── examples/
    └── webgl-app/       # Complete WebGL+WebSocket example
        ├── *.go        # Go source
        └── *.extern.go # Browser API bindings

```

## External Bindings Pattern

Use `.extern.go` files to define browser API type signatures:

```go
// dom.extern.go
type HTMLCanvasElement struct{}
func (c HTMLCanvasElement) GetContext(ctx string) WebGLRenderingContext

type WebGLRenderingContext struct{}
func (gl WebGLRenderingContext) ClearColor(r, g, b, a float64)

var Document DocumentObject
```

The linker maps these to JavaScript equivalents:
- `Document` → `document`
- `GetElementsByTagName` → `getElementsByTagName`
- Go methods → JS methods (PascalCase → camelCase)

## Testing the Transpiler

No formal test suite exists. Manual testing workflow:

1. Write test Go code in a directory
2. Run: `go run main.go ./test-dir > output.js`
3. Create HTML file that includes output.js
4. Open in browser and check console for errors
5. Verify generated JS is readable and correct

For the webgl-app example:
```bash
go run main.go examples/webgl-app > examples/webgl-app/output.js
# Then open HTML file that loads output.js and calls Main()
```

## Supported Go Features

**Supported:**
- Structs (→ JS classes)
- Functions and methods
- Basic types (int, float64, string, bool)
- Slices and arrays
- Maps (→ JS Objects)
- if/else, for loops, switch
- Composite literals
- Method calls and field access
- Type conversions

**NOT Supported:**
- Goroutines, channels
- Interfaces
- Reflection
- Generics (complex)
- defer/panic/recover
- Most stdlib packages (fmt, io, net, etc.)

## Important Design Constraints

1. **Minimal subset only**: Don't add features that bloat the transpiler. Keep it educational.
2. **Clean output**: Generated JS should be human-readable, not obfuscated.
3. **No runtime**: Unlike GopherJS, this has no Go runtime in JS.
4. **Browser-focused**: Targets browser APIs, not Node.js.
5. **Educational priority**: Code clarity > performance > features.

## Module Definition

```go
module github.com/rryanrussell/goscript

go 1.18
```

No external dependencies - uses only Go stdlib (go/parser, go/ast, go/token).

## Common Patterns When Modifying

**Adding a new Go→JS transformation:**
1. Identify the Go AST node type in transpiler/
2. Add case in appropriate transpiler function
3. Create corresponding JS AST node in js/
4. Implement Print() method for the JS node
5. Test with example Go code

**Adding a new external binding:**
1. Create or update `.extern.go` file with Go signatures
2. Add mapping in link/builtin.go if needed
3. Linker automatically maps PascalCase → camelCase

**Debugging transformation issues:**
1. Add debug prints in transpiler/ to see Go AST structure
2. Check JS AST generation in js/ package
3. Verify output in generated JavaScript
4. Compare with expected JS for the Go input

## Project Status

This is a **finished experiment** (circa 2022, extracted from mathfish project). It's released for educational value, not production use.

**Maintenance approach:**
- Accept bug fixes for existing features
- Accept documentation improvements
- Be cautious about feature additions (keep it minimal)
- Redirect production use questions to Go WASM or GopherJS

---

## AI-Assisted Development & Attribution

Parts of this project have been developed with assistance from Claude (AI model by Anthropic) via Claude Code (AI coding agent CLI).

### When to Credit AI in Commits

**If you (Claude) are the sole author of a commit**, include attribution in the commit message:

#### Format for AI-authored commits:

```
<conventional-commit-title>

<detailed description of changes>

---
Co-authored-by: Claude (Sonnet 4.5) <noreply@anthropic.com>
Generated via: Claude Code CLI (claude.ai/code)
Prompt: <brief characterization of the user's request>
```

#### Examples:

**Example 1: Feature implementation**
```
feat: add else clause support to if statements

Implement else and else-if clause handling in the transpiler:
- Extended js.IfStmt to include Else field
- Added transpilation logic for ast.IfStmt.Else
- Updated pretty-printer to output else blocks
- Added test cases for nested if-else chains

---
Co-authored-by: Claude (Sonnet 4.5) <noreply@anthropic.com>
Generated via: Claude Code CLI (claude.ai/code)
Prompt: Implement else clause support for if statements
```

**Example 2: Documentation**
```
docs: create roadmap for Go runtime expansion

Create comprehensive ROADMAP.md with 6-phase plan:
- Phase 1: Foundation (control flow, builtins)
- Phase 2: Stdlib essentials (strings, fmt, errors, math)
- Phase 3: Interfaces and method dispatch
- Phase 4: Defer and panic/recover
- Phase 5: Goroutines via async/await and channels
- Phase 6: Limited reflection support

Includes timeline, runtime size projections, and success metrics.

---
Co-authored-by: Claude (Sonnet 4.5) <noreply@anthropic.com>
Generated via: Claude Code CLI (claude.ai/code)
Prompt: Create ambitious roadmap for expanding Go runtime support in JS
```

**Example 3: Bug fix**
```
fix: handle nil slices in append() transformation

Fixed panic when appending to nil slice by adding null coalescing:
- Changed: [...arr, elem]
- To: [...(arr ?? []), elem]

This matches Go's behavior where append to nil slice creates new slice.

---
Co-authored-by: Claude (Sonnet 4.5) <noreply@anthropic.com>
Generated via: Claude Code CLI (claude.ai/code)
Prompt: Fix nil slice handling in append transformation
```

### When NOT to include AI attribution:

- **Human-authored commits** - If a human wrote the code/docs, credit only the human
- **Pair programming** - If human and AI collaborated significantly, use `Co-authored-by: Human Name <email>` for human, optionally note AI assistance in commit body
- **Minor edits** - Trivial changes (typo fixes, formatting) don't need AI attribution
- **Build/CI commits** - Automated commits from CI/CD don't need attribution

### Characterizing the Prompt

Keep prompt characterization **brief and meaningful**:

**Good examples:**
- "Implement slice expression support with bounds checking"
- "Add strings package runtime implementation"
- "Refactor transpiler to support interface dispatch"
- "Create test suite for control flow transformations"
- "Debug goroutine transformation edge cases"

**Too vague (avoid):**
- "Make changes to code"
- "Fix issues"
- "Update files"

**Too verbose (avoid):**
- "The user requested that I implement support for slice expressions including proper handling of omitted bounds, negative indices, and out-of-range checking, with consideration for both compile-time and runtime approaches, ultimately deciding on..."

### Model Version Tracking

**Current model**: Claude Sonnet 4.5 (`claude-sonnet-4-5-20250929`)

If using a different Claude model version, update the Co-authored-by line:
- `Claude (Opus 4)` - for Opus model
- `Claude (Haiku 4)` - for Haiku model
- `Claude (Sonnet 4.5)` - for Sonnet 4.5 (current)

### Agent Runtime

Always specify "Claude Code CLI" as the agent runtime when applicable. This distinguishes from direct Claude API usage or other interfaces.

---

## AI Development Guidelines

When working on this project as Claude:

1. **Maintain the minimal philosophy** - Don't over-engineer or add unnecessary features
2. **Prioritize readable output** - Generated JS should be clean and understandable
3. **Test thoroughly** - Manually verify transpiled output works in browser
4. **Document decisions** - Explain why you chose compile-time vs runtime approaches
5. **Update ROADMAP.md** - Mark items complete and adjust timeline as needed
6. **Credit appropriately** - Use commit attribution format above for your work
7. **Ask when uncertain** - Use AskUserQuestion for architectural decisions

### Collaboration with Humans

When a human modifies your AI-generated code:
- The human becomes the author (remove AI attribution in subsequent commits)
- Original AI contribution remains in git history
- This is expected and encouraged (AI is a tool, not the project owner)
