package transpiler

import "testing"

func TestStringsPackage(t *testing.T) {
	src := `
package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "hello world"
	if strings.Contains(s, "hello") {
		fmt.Println("contains hello")
	}
	if strings.HasPrefix(s, "he") {
		fmt.Println("has prefix he")
	}
	if strings.HasSuffix(s, "ld") {
		fmt.Println("has suffix ld")
	}
	upper := strings.ToUpper(s)
	lower := strings.ToLower(upper)
	parts := strings.Split(s, " ")
	joined := strings.Join(parts, "-")

	rep := strings.Replace(s, "l", "L", 2)
	repAll := strings.ReplaceAll(s, "l", "L")
    repOverlap := strings.Replace("banana", "a", "o", 2)
    repInsert := strings.Replace("foo", "o", "oo", 2)
    repNested := strings.Replace("a", "a", "aa", 1)

    // Edge cases
    repEmpty := strings.Replace("abc", "", "X", -1)
    repEmptyLimit := strings.Replace("abc", "", "X", 2)

	trimmed := strings.TrimSpace("  hello  ")
    trimmedSet := strings.Trim("1-2-3", "1-3")
    trimmedEmpty := strings.Trim("abc", "")

    splitEmpty := strings.Split("abc", "")
    splitEmptyStr := strings.Split("", "")
}
`
	expected := `function main() {
    let s = "hello world"
    if (strings.Contains(s, "hello")) {
        fmt.Println("contains hello")
    }

    if (strings.HasPrefix(s, "he")) {
        fmt.Println("has prefix he")
    }

    if (strings.HasSuffix(s, "ld")) {
        fmt.Println("has suffix ld")
    }

    let upper = strings.ToUpper(s)
    let lower = strings.ToLower(upper)
    let parts = strings.Split(s, " ")
    let joined = strings.Join(parts, "-")
    let rep = strings.Replace(s, "l", "L", 2)
    let repAll = strings.ReplaceAll(s, "l", "L")
    let repOverlap = strings.Replace("banana", "a", "o", 2)
    let repInsert = strings.Replace("foo", "o", "oo", 2)
    let repNested = strings.Replace("a", "a", "aa", 1)
    let repEmpty = strings.Replace("abc", "", "X", -1)
    let repEmptyLimit = strings.Replace("abc", "", "X", 2)
    let trimmed = strings.TrimSpace("  hello  ")
    let trimmedSet = strings.Trim("1-2-3", "1-3")
    let trimmedEmpty = strings.Trim("abc", "")
    let splitEmpty = strings.Split("abc", "")
    let splitEmptyStr = strings.Split("", "")
}
`
	runTranspilerTest(t, "StringsPackage", src, expected)
}
