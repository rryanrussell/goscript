# goscript - Minimal Go-to-JavaScript Transpiler

A minimal Go-to-JavaScript transpiler focused on clean, readable output for browser UI code.

⚠️ **Project Status: Experimental / Educational**

This is a *finished experiment* demonstrating a minimal transpiler design. It's not intended to replace [GopherJS](https://github.com/gopherjs/gopherjs) or WebAssembly, but explores a different design point in the Go→JS compilation space.

## Design Philosophy

**goscript** targets a specific niche:
- ✅ **Small Go subset** - Structs, functions, basic control flow
- ✅ **Clean JavaScript output** - Human-readable, minimal runtime
- ✅ **Direct DOM bindings** - WebGL, WebSocket via `.extern.go` files
- ✅ **No runtime overhead** - Unlike full Go compilers

**Not supported:**
- ❌ Goroutines, channels, interfaces
- ❌ Full Go standard library
- ❌ Reflection, complex generics
- ❌ Production support or maintenance guarantees

## Use This If

- 📚 **Learning compiler implementation** - Clean ~2k LOC transpiler
- 🔬 **Research/experiments** - Comparing transpiler approaches
- 📦 **Size-constrained environments** - Minimal JS output matters
- 🎓 **Teaching** - Real working example for courses

## Don't Use This If

- 🏭 **Need production support** → Use [Go WebAssembly](https://github.com/golang/go/wiki/WebAssembly)
- 🔧 **Need full Go language** → Use [GopherJS](https://github.com/gopherjs/gopherjs)
- 🚀 **Building serious web apps** → Use TypeScript or standard Go WASM

## Quick Example

**Input (Go):**
```go
package app

type Vec2 struct {
    X float64
    Y float64
}

func Main() {
    pos := Vec2{X: 10, Y: 20}
    canvas := Document.GetElementsByTagName("canvas")[0]
    gl := canvas.GetContext("webgl")
    gl.ClearColor(0, 0, 0, 1)
}
```

**Output (JavaScript):**
```javascript
class Vec2 {
    constructor() {
        this.X = 0;
        this.Y = 0;
    }
}

function Main() {
    const pos = new Vec2();
    pos.X = 10;
    pos.Y = 20;
    const canvas = document.getElementsByTagName("canvas")[0];
    const gl = canvas.getContext("webgl");
    gl.clearColor(0, 0, 0, 1);
}
```

## How It Works

1. **Parse Go code** using `go/parser` and `go/ast`
2. **Transform to JS AST** - Custom JavaScript AST representation
3. **Link external symbols** - Map Go types to browser APIs
4. **Generate JavaScript** - Pretty-print readable output

### External Bindings

Define browser APIs with `.extern.go` files:

```go
// dom.extern.go
type HTMLCanvasElement struct{}
func (c HTMLCanvasElement) GetContext(ctx string) WebGLRenderingContext

type WebGLRenderingContext struct{}
func (gl WebGLRenderingContext) ClearColor(r, g, b, a float64)
```

The transpiler recognizes these as external symbols and generates appropriate JS calls.

### Dommer: Automatic Bindings

`dommer` is a tool included in `goscript` that generates Go extern definitions from TypeScript declaration files (`.d.ts`).

- **Input:** TypeScript definition files (e.g., `lib.dom.d.ts`)
- **Output:** Go extern files (e.g., `dom.extern.go`)
- **Purpose:** Automates the creation of type-safe browser API bindings.

```bash
go run cmd/dommer/main.go -input lib/lib.dom.d.ts -output app/dom.extern.go
```

## Architecture

```
┌─────────────┐
│  Go Source  │
└──────┬──────┘
       │ go/parser
       ▼
┌─────────────┐
│   Go AST    │
└──────┬──────┘
       │ transpiler/
       ▼
┌─────────────┐
│   JS AST    │  (js/ package)
└──────┬──────┘
       │ link/
       ▼
┌─────────────┐
│ JavaScript  │
└─────────────┘
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for deep dive.

## Example: WebGL Application

See [examples/webgl-app/](examples/webgl-app/) for a complete WebGL+WebSocket demo:
- Real-time rendering with WebGL
- Binary WebSocket communication
- Go structs mapped to JS objects
- ~150 LOC of Go → clean JavaScript

## Installation

```bash
go get github.com/rryanrussell/goscript
```

## Usage

```go
package main

import (
    "os"
    "goscript"
)

func main() {
    // Transpile directory to stdout
    if err := goscript.Transpile("./myapp", os.Stdout); err != nil {
        panic(err)
    }
}
```

Or use as CLI:

```bash
go run github.com/rryanrussell/goscript ./myapp > output.js
```

## Comparison to Other Tools

| Tool | Approach | Output Size | Go Coverage | Status |
|------|----------|-------------|-------------|--------|
| **goscript** | Minimal subset | Small | ~15% | Experimental |
| **GopherJS** | Full compiler | Large | ~95% | Maintained |
| **Go WASM** | Native WASM | Medium | 100% | Official |
| **Joy** | Idiomatic subset | Small | ~30% | Abandoned (2017) |

**goscript** is closest to [Joy](https://mat.tm/joy) in philosophy but more recent and focused on direct DOM access.

## Limitations

**Language Features:**
- No goroutines or channels
- No interfaces or reflection
- No defer/panic/recover
- Limited standard library (no fmt, io, etc.)
- Basic type system only

**Output:**
- No source maps (yet)
- No optimization passes
- No minification
- Requires manual extern definitions (automated via `dommer`)

**Stability:**
- No semantic versioning
- Breaking changes possible
- Minimal test coverage
- Educational project, not production-ready

See [ROADMAP.md](ROADMAP.md) for future plans and known limitations to be addressed.

## Contributing

**This is a finished experiment**, but contributions welcome for:

✅ **Documentation improvements**
✅ **Bug fixes**
✅ **Example additions**
✅ **Educational materials**

⚠️ **Feature additions** will be evaluated carefully to maintain the "minimal" philosophy.

**For major features**, consider forking! This project intentionally stays small.

### Development

```bash
# Clone
git clone https://github.com/rryanrussell/goscript
cd goscript

# Run example
cd examples/webgl-app
go run ../../main.go . > output.js

# Inspect JS AST
cd js
go doc
```

## Project History

Developed in 2022 as part of a genetic algorithm framework for training neural networks. The transpiler enabled writing browser UI code in Go while keeping the JS output minimal and readable. Extracted and released standalone in 2025 for educational value.

## Why This Exists

**Original need:** Write browser UI in Go without 2MB WASM runtime overhead.

**Broader value:** Demonstrates transpiler design for minimal language subsets.

**Lesson learned:** WASM is the right choice for most cases, but exploring alternatives teaches you a lot about both compilers and the constraints they navigate.

## Alternatives You Should Consider

- **[Go WebAssembly](https://github.com/golang/go/wiki/WebAssembly)** - Official, full Go support
- **[GopherJS](https://github.com/gopherjs/gopherjs)** - Mature, extensive compatibility
- **[TinyGo WASM](https://tinygo.org/docs/guides/webassembly/)** - Smaller WASM builds
- **TypeScript** - If you don't need Go specifically

## License

MIT License - see [LICENSE](LICENSE)

## Author

Ryan Russell ([@rryanrussell](https://github.com/rryanrussell))

Built with curiosity about minimal transpilers and a desire for clean JavaScript output.

---

**Remember:** This is an educational project showing one approach to Go→JS transpilation. It's not trying to be GopherJS. It's not trying to replace WASM. It's exploring a specific design point and hopefully teaching something along the way.

If you build something cool with it, let me know! If you use it to learn compilers, even better.
