package main

import (
	"fmt"
	"testing"
)

func TestMakeSlice(t *testing.T) {
	s := make([]int, 5)
	if len(s) != 5 {
		t.Errorf("expected len 5, got %d", len(s))
	}
	s[0] = 10
	if s[0] != 10 {
		t.Errorf("expected s[0] to be 10, got %d", s[0])
	}
}

func TestMakeMap(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	if m["a"] != 1 {
		t.Errorf("expected m['a'] to be 1, got %d", m["a"])
	}
}

func main() {
	// This main function won't be run by go test, but it is useful if we run this file directly with transpiler.
	s := make([]int, 5)
	fmt.Println(len(s))
	s[0] = 10
	fmt.Println(s[0])

	m := make(map[string]int)
	m["a"] = 1
	fmt.Println(m["a"])
}
