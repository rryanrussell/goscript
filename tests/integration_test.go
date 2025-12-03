package tests

import (
	"testing"
)

func TestBasicVariables(t *testing.T) {
	src := `package main

var x = 10
var y = 20
var z = x + y

func main() {
}
`
	vm, err := RunTranspiled(src)
	if err != nil {
		t.Fatalf("RunTranspiled failed: %v", err)
	}

	// In the current transpiler, package-level variables are usually global in the generated JS
	// or part of a package object. Let's see how they are generated.
	// Based on experience with this codebase or typical transpilation, they might be top-level.

	// Note: Currently top-level vars use 'let', so they are not attached to the global object in Goja by default
	// unless we are in the same scope. However, RunString executes in the global scope but `let`
	// declarations are block scoped to the script execution context in some engines or just not on globalThis.
	// But in Goja's RunString, `let` variables defined at the top level should be accessible if we query them
	// from the same runtime context, BUT `vm.Get("x")` retrieves from the global object (globalThis).
	// `let` variables do NOT become properties of the global object.
	// So we might need to change how we verify them, or how they are declared (var vs let).
	// For now, let's try to verify by running a small script that returns them.

	val, err := vm.RunString("[x, y, z]")
	if err != nil {
		t.Fatalf("Failed to inspect variables: %v", err)
	}
	export := val.Export().([]interface{})

	if export[0].(int64) != 10 {
		t.Errorf("Expected x to be 10, got %v", export[0])
	}
	if export[1].(int64) != 20 {
		t.Errorf("Expected y to be 20, got %v", export[1])
	}
	if export[2].(int64) != 30 {
		t.Errorf("Expected z to be 30, got %v", export[2])
	}
}

func TestFunctionCall(t *testing.T) {
	src := `package main

func add(a, b int) int {
	return a + b
}

var result = add(5, 7)

func main() {}
`
	vm, err := RunTranspiled(src)
	if err != nil {
		t.Fatalf("RunTranspiled failed: %v", err)
	}

	// Since `result` is a variable, it should be in the global scope.
	// Note: If the transpiler outputs `let result = ...` at the top level,
	// Goja might treat it as block-scoped if the script is not run as a module or if there's implicit wrapping.
	// However, `RunString` usually executes in the global scope.
	// Let's check if the variable is defined.
	val, err := vm.RunString("result")
	if err != nil {
		t.Fatalf("Failed to inspect result: %v", err)
	}
	if val.ToInteger() != 12 {
		t.Errorf("Expected result to be 12, got %v", val)
	}
}
