# goscript Development Roadmap

## Vision

Expand goscript from a minimal transpiler to a more complete Go→JavaScript compiler while maintaining **clean, readable output** as the north star. Accept modest runtime overhead (~10-50KB) when necessary for core features.

## Philosophy

- **Coherent JavaScript preferred**: Compile-time transformations over runtime emulation
- **Balanced runtime**: Accept runtime helpers for features that can't be cleanly compiled
- **Educational clarity**: Keep code readable; this is a learning tool
- **No WASM**: Pure JavaScript output (browsers can optimize it)

---

## Phase 1: Foundation (Complete Current Gaps) [Small]

**Status**: Currently ~70% of basic Go supported

### Missing Control Flow
- [x] **Else/Else-if clauses** (js/stmt.go, transpiler/transpiler.go)
  - **Status**: Implemented
  - **Approach**: Compile-time - extend IfStmt in JS AST

- [x] **Switch statements** (new: js/switch.go, transpiler case)
  - **Status**: Implemented (basic value switch)
  - **Approach**: Compile-time - transform to if/else chain or JS switch

- [x] **Break/Continue** (js/stmt.go)
  - **Status**: Implemented
  - **Approach**: Compile-time - direct JS equivalents

- [x] **Full Range Loop Support**
  - Currently `for i, v := range arr` is supported
  - `for i := range arr` and `for _, v := range arr` work
  - **Approach**: Compile-time transformation

- [ ] **Labeled statements & goto** (low priority)
  - **Approach**: Integrated via "State Machine Frame" for backward jumps; simple labels for forward jumps.

### Missing Built-in Operations
- [x] **Slice expressions** `arr[1:3]` (transpiler.go line ~200)
  - **Status**: Implemented (2-index slicing)
  - **Approach**: Compile-time → `arr.slice(1, 3)`

- [x] **`copy()` builtin** (runtime or compile-time)
  - **Status**: Implemented
  - **Approach**: Runtime function `runtime.copy(dst, src)`

- [x] **`make()` builtin** for slices/maps
  - **Status**: Implemented
  - **Approach**: Runtime function
  - `make([]int, 5)` → `runtime.makeSlice(5, 0)` (returns array of length 5)
  - `make(map[K]V)` → `runtime.makeMap()` (returns JS object)

- [x] **Better `append()` support**
  - **Status**: Implemented
  - **Approach**: Compile-time expansion with spread operator
  - `append(arr, a, b, c)` → `[...(arr ?? []), a, b, c]`
  - `append(arr, other...)` → `[...(arr ?? []), ...(other ?? [])]`

- [ ] **Refactor `call` function in `transpiler/transpiler.go`**
  - **Goal**: Split built-in handling (append, make, len, etc.) into separate functions to improve readability and maintainability.

### Type System Basics
- [ ] **Type switches** (requires runtime type info)
  - Currently `switch x.(type)` is not handled
  - **Approach**: Runtime - attach `__goType` property to values
  - Compile switch to if/else checking `__goType`

- [ ] **Better type assertions** (currently tracked but not enforced)
  - **Approach**: Runtime - check `__goType`, panic on mismatch

- [ ] **Map operations** (currently no special handling)
  - **Approach**: Runtime - use JS Map with helpers
  - `m[k]` → `runtime.mapGet(m, k, defaultValue)`
  - `m[k] = v` → `runtime.mapSet(m, k, v)`

### Multiple Return Values
- [ ] **Better unpacking** (currently packs to array)
  - **Approach**: Compile-time - detect context
  - `a, b := foo()` → `const [a, b] = foo()`
  - `a, ok := m[k]` → Special case for map access

---

## Phase 2: Standard Library Essentials [Medium]

**Goal**: Provide JS implementations of critical Go stdlib packages

### strings package (runtime: runtime/strings.js)
- [x] `strings.Contains(s, substr)` → `s.includes(substr)`
- [x] `strings.HasPrefix(s, prefix)` → `s.startsWith(prefix)`
- [x] `strings.HasSuffix(s, suffix)` → `s.endsWith(suffix)`
- [x] `strings.Split(s, sep)` → `s.split(sep)`
- [x] `strings.Join(arr, sep)` → `arr.join(sep)`
- [x] `strings.ToUpper/ToLower` → `s.toUpperCase()` / `s.toLowerCase()`
- [x] `strings.Trim, TrimSpace` → regex or manual
- [x] `strings.Replace` → `s.replaceAll(old, new)`

### fmt package (expand link/fmt.js)
- [x] `fmt.Println` (already exists)
- [x] `fmt.Errorf` (basic version exists)
- [ ] `fmt.Sprintf` - proper format string parsing
- [ ] `fmt.Printf` → console.log with formatting
- [ ] Support for `%v`, `%s`, `%d`, `%f`, `%t`, `%x`, `%p` verbs

### errors package (runtime: runtime/errors.js)
- [ ] `errors.New(msg)` → `new Error(msg)`
- [ ] `errors.Is(err, target)` → error chain checking
- [ ] `errors.As(err, target)` → type assertion for errors
- [ ] `fmt.Errorf` with `%w` for error wrapping

### strconv package (runtime: runtime/strconv.js)
- [ ] `strconv.Atoi(s)` → `parseInt(s, 10)`
- [ ] `strconv.Itoa(n)` → `String(n)`
- [ ] `strconv.ParseFloat(s, bits)` → `parseFloat(s)`
- [ ] `strconv.FormatInt/FormatFloat` → String() with radix

### math package (runtime: runtime/math.js)
- [ ] Constants: `math.Pi`, `math.E`, etc. → `Math.PI`, `Math.E`
- [ ] Functions: `math.Sqrt`, `math.Pow`, `math.Sin`, etc. → `Math.*`
- [ ] `math.Floor/Ceil/Round` → `Math.floor()` etc.
- [ ] `math.Max/Min` → `Math.max/min`
- [ ] `math.Abs` → `Math.abs`

### time package (runtime: runtime/time.js)
- [ ] `time.Now()` → `new Date()`
- [ ] `time.Since(t)` → `Date.now() - t`
- [ ] `time.Sleep(d)` → `await runtime.sleep(ms)` (requires async)
- [ ] `time.Duration` type → number (milliseconds)
- [ ] Basic duration parsing

### sort package (runtime: runtime/sort.js)
- [ ] `sort.Ints(arr)` → `arr.sort((a,b) => a - b)`
- [ ] `sort.Strings(arr)` → `arr.sort()`
- [ ] `sort.Sort(data)` → custom comparator
- [ ] `sort.Slice(arr, less)` → `arr.sort(less)`

---

## Phase 3: Interfaces & Method Dispatch [Medium]

**Goal**: Enable polymorphism via interface types

### Runtime Type Information
- [ ] **Attach type metadata** to all values
  - Structs get `__goType: "pkg.StructName"`
  - Functions get `__goType: "func(...)"`

### Interface Implementation
- [ ] **Interface types** (transpiler support)
  - Parse `type I interface { Method() }` in Go AST
  - Generate JS class for interface wrapper

- [ ] **Interface assignments** (runtime checks)
  - `var i I = structValue` → check if struct has required methods
  - Wrap struct in interface proxy object

- [ ] **Method dispatch** (runtime)
  - `i.Method()` → lookup method in proxy, forward to underlying value

- [ ] **Type assertions** (runtime)
  - `v := i.(ConcreteType)` → unwrap interface, check type, panic if wrong
  - `v, ok := i.(ConcreteType)` → return (value, false) on mismatch

**Approach**:
```javascript
// Runtime support
class Interface {
  constructor(value, methods) {
    this.__value = value;
    this.__methods = methods; // map of method name → bound function
  }
  // Proxy all interface methods
}

function assertInterface(value, requiredMethods) {
  // Check if value has all required methods
  // Return Interface wrapper
}
```



---

## Phase 4: Defer & Panic/Recover [Medium]

**New Approach**: Use **Universal Frame Template** to wrap functions detecting these effects.

### Defer Statement
- [ ] **Frame Selection**
  - Scan for `defer` keyword
  - Select "Bracket Frame" (try/finally)
  - Hoist `__deferred` array in prologue

**Example**:
```javascript
function foo() {
  const __deferred = [];
  try {
     __deferred.push(cleanup);
     // ...
  } finally {
     runDefers(__deferred);
  }
}
```

### Panic/Recover
- [ ] **Frame Selection**
  - Scan for `recover`
  - Select "Exception Frame" (try/catch) or merge with Bracket Frame
  - Catch block handles `panic` value, sets it for `recover()` to find



---

## Phase 5: Goroutines & Channels [Large]

**New Approach**: Treat `go`, channel ops, and `select` as **Yield Effects** triggering **Coroutine Frames**.

### Goroutines & Channels
- [ ] **Effect Analysis**
  - Scan for `go`, `ch <-`, `<-ch`, `select`
  - Mark function as "Yielding" (requires generator/async)

- [ ] **Frame Selection**
  - Wrap body in "Coroutine Frame" (async generator or state machine)
  - All channel ops become `yield*` calls to runtime

**Example**:
```javascript
function* worker() {
  yield* runtime.send(ch, 1);
}
```

### Select Statement
- [ ] **Runtime Implementation**
  - `runtime.select` manages the race logic
  - Transpiles to `yield* runtime.select([...cases])`

**Example**:
```go
select {
case v := <-ch1:
  handle(v)
case ch2 <- val:
  // sent
default:
  // default
}
```

**Transpile to**:
```javascript
const result = await runtime.select([
  { type: 'recv', chan: ch1, handler: (v) => handle(v) },
  { type: 'send', chan: ch2, value: val, handler: () => {} },
  { type: 'default', handler: () => {} }
]);
```



---

## Phase 6: Reflection (Limited) [Large]

**Goal**: Basic reflection for JSON marshaling, etc.

- [ ] **Type metadata** (expand runtime type info)
  - `reflect.TypeOf(v)` → return type descriptor
  - `reflect.ValueOf(v)` → return value wrapper

- [ ] **Struct field iteration**
  - `t.NumField()`, `t.Field(i)` → introspect struct

- [ ] **Field access by name**
  - `v.FieldByName("X")` → get/set field

- [ ] **Method calls by name**
  - `v.MethodByName("Foo")` → get method

**Approach**:
- Generate metadata at compile time for each struct
- Store in `__goReflect` property
- Runtime functions access metadata



---

## Phase 7: Advanced Features [Future]

### Generics (Type Parameters)
- [ ] **Simple generics** for functions
  - Compile-time monomorphization (generate version per type)
  - Or runtime with type erasure (like TypeScript)

- [ ] **Generic types**
  - `type Stack[T any] struct { items []T }`
  - Generate JS class per instantiation

### Method Values
- [ ] Currently limited support
- [ ] Bind `this` correctly for method values
- [ ] `f := obj.Method` → `f = obj.Method.bind(obj)`

### Closures (Improved)
- [ ] Basic closures work (FuncLit)
- [ ] Ensure captured variables work correctly
- [ ] Handle closure over loop variables (common Go gotcha)

---

## Runtime Size Projections

| Phase | Feature Set | Estimated Runtime Size |
|-------|-------------|------------------------|
| **Phase 1** | Control flow, builtins, maps | ~5-8 KB |
| **Phase 2** | Stdlib (strings, fmt, errors, math) | ~15-25 KB |
| **Phase 3** | Interfaces & type assertions | ~25-40 KB |
| **Phase 4** | Defer, panic/recover | ~30-45 KB |
| **Phase 5** | Goroutines, channels, select | **50-75 KB** |
| **Phase 6** | Reflection (limited) | ~65-95 KB |

**Target**: Stay under 100KB total runtime (minified)

---

## Code Organization

### New files to create:

```
goscript/
├── runtime/
│   ├── runtime.js          # Core runtime (make, copy, type checks)
│   ├── strings.js          # strings package
│   ├── fmt.js              # fmt package (move from link/fmt.js)
│   ├── errors.js           # errors package
│   ├── strconv.js          # strconv package
│   ├── math.js             # math package
│   ├── time.js             # time package
│   ├── sort.js             # sort package
│   ├── interface.js        # Interface dispatch
│   ├── defer.js            # Defer/panic/recover
│   ├── chan.js             # Channels & goroutines
│   └── reflect.js          # Reflection (limited)
├── transpiler/
│   ├── transpiler.go       # Existing - add new AST cases
│   ├── async.go            # NEW - goroutine→async transformation
│   └── defer.go            # NEW - defer→try/finally transformation
└── link/
    └── runtime.go          # NEW - include runtime/*.js files
```

---

## Testing Strategy

### Current: Goja Runtime Integration

We have integrated **Goja** (a pure Go JavaScript runtime) to allow direct testing of transpiled code within the Go test suite.

### Proposed:
1. **Unit tests** for transpiler (transpiler_test.go)
   - Test each Go construct → expected JS output
   - Skipped tests can be used as TODOs for partially supported or planned features.

2. **Integration tests** (tests/)
   - **Status**: Active (Goja-based)
   - We use `tests/runner.go` to transpile Go code and execute it immediately in Goja.
   - We can inspect variables and return values in the JS runtime to verify correctness.

3. **Browser tests** (tests/browser/)
   - Load transpiled JS in headless browser
   - Verify DOM manipulation works

4. **Benchmark suite**
   - Measure transpilation time
   - Measure runtime performance vs native JS
   - Track runtime size growth

---

## Success Metrics

| Metric | Current | Phase 1 | Phase 3 | Phase 5 |
|--------|---------|---------|---------|---------|
| **Go coverage** | ~15% | ~40% | ~65% | ~85% |
| **Runtime size** | ~1 KB | ~8 KB | ~40 KB | ~75 KB |
| **Stdlib packages** | 1 (partial fmt) | 5 | 8 | 12 |
| **Example apps** | 1 (webgl) | 3 | 5 | 8 |
| **Test coverage** | 0% | 30% | 50% | 70% |

---

## Non-Goals

**Will NOT support:**
- Full standard library (os, net/http, database/sql, etc.)
- CGo or system calls
- Assembly code
- Build tags / conditional compilation
- Go modules / dependency management (use browser imports)
- Full compatibility with GopherJS (different goals)

**Use WASM for:**
- Production applications requiring full Go
- High-performance computing
- Existing Go codebases without modification

---

## Timeline Summary

| Phase | Effort | Key Deliverable |
|-------|----------|-----------------|
| **Phase 1** | Small | Complete basic Go support |
| **Phase 2** | Medium | Essential stdlib packages |
| **Phase 3** | Medium | Interfaces working |
| **Phase 4** | Medium | Defer/panic/recover |
| **Phase 5** | Large | Async/await goroutines |
| **Phase 6+** | Large | Reflection, advanced features |

**Total**: Significant effort to reach ~85% Go coverage with clean JS output

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to contribute to this roadmap.

**Priority areas for contributors:**
- Phase 1 features (else clauses, switch, slice expressions)
- Stdlib package implementations (strings, math, etc.)
- Test suite development
- Documentation improvements
- Example applications

---

## Credits

This roadmap was developed with assistance from:
- **Claude Sonnet 4.5** (claude-sonnet-4-5-20250929) - AI model by Anthropic
- **Claude Code** - AI-powered coding agent CLI (claude.ai/code)

Original transpiler implementation (2022) by Ryan Russell.

Roadmap planning session (December 2025) utilized Claude Code's planning agents to analyze the existing codebase and design expansion strategies.
