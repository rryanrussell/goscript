# Dommer Roadmap

This document outlines the roadmap for expanding `dommer` to support more of the TypeScript `.d.ts` specification, specifically targeting `lib.dom.d.ts`.

## Current Status

- **Parser**: The parser (`pkg/dommer/parser`) and AST (`pkg/dommer/ast`) appear to support a wide range of TypeScript constructs, including interfaces, globals, functions, namespaces, and various type expressions.
- **Generator**: The generator (`pkg/dommer/generator.go`) is currently the bottleneck. It contains a hardcoded allowlist that restricts output to only a few specific WebGL-related interfaces and types.

## Roadmap

### Phase 1: Unlocking the Generator

The immediate goal is to remove the artificial limits in the generator and allow it to attempt processing all declarations.

- [x] **Remove Allowlist**: Remove the `switch decl.Ident.Name.Value` checks in `Generate` function in `pkg/dommer/generator.go`.
- [x] **Default Handlers**: Ensure there are default handlers (even if they just log a warning or skip) for all `DtsDecl` types so the generator doesn't panic on unknown declarations.

### Phase 2: Global Scope Support

Support for top-level declarations that are currently ignored.

- [x] **Global Variables**: Implement `GlobalVarDecl` handling.
    - Map `declare var foo: Type` to `var Foo Type`.
- [x] **Global Constants**: Implement `GlobalConstDecl` handling.
    - Map `declare const foo: Type` to `var Foo Type` (since Go consts are limited).
- [x] **Global Functions**: Implement `FunctionDecl` handling.
    - Map `declare function foo(args): Ret` to `func Foo(args) Ret`.

### Phase 2.5: Generator Robustness

Fix issues discovered after unlocking the generator.

- [x] **Keyword Escaping**: Escape Go keywords in parameter names and identifiers (e.g., `range` -> `range0`, `select` -> `select0`).
- [ ] **Numeric Literals**: Handle numeric literal types or constants in interfaces (e.g., `TIMEOUT_IGNORED -1`).
- [x] **Invalid Identifiers**: Handle identifiers with invalid characters (e.g., `-`) and conflicts between var/type names.

### Phase 3: Advanced Type Support

Improve mapping of TypeScript types to Go types.

- [ ] **Unions**: Improve `TypeUnionExpr` handling.
    - Current: Handles `T | null`.
    - Goal: Handle `A | B` (likely as `any` or a generated interface).
- [ ] **Intersections**: Implement `TypeIntersectionExpr`.
    - Map `A & B` to a struct embedding both or an interface.
- [ ] **Function Types**: Handle `LambdaTypeExpr` (arrow functions) in arguments/variables.
- [ ] **Arrays**: Verify `TypeArrayExpr` works for all types.

### Phase 4: Interface Features

Support richer interface definitions.

- [ ] **Indexers**: Implement `IndexerDecl` (e.g., `[key: string]: number`).
    - Could map to a method `Get(key string) number` or similar.
- [ ] **Call Signatures**: Implement `CallableDecl` (interfaces that are callable).
- [ ] **Getters/Setters**: Implement `GetterDecl` and `SetterDecl`.

### Phase 5: Namespaces

- [ ] **Namespace Support**: Implement `NamespaceDecl`.
    - TypeScript namespaces are often used for grouping. In Go, these might need to be flattened with prefixes or mapped to separate packages (though separate packages might be complex for a single `.d.ts` file).

### Phase 6: Go Idioms

- [ ] **Interfaces vs Structs**: Currently, most things are generated as structs. Consider generating Go interfaces for TypeScript interfaces to allow for better polymorphism and mocking.
- [ ] **Enums**: Improve enum generation to use Go's `iota` or string constants (partially implemented).
