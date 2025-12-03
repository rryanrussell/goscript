package transpiler

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"

	"github.com/rryanrussell/goscript/js"
)

func runTranspilerTest(t *testing.T, name, src, expected string) {
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

	// We need to initialize the context properly, similar to main.go
	fCtx.Module.Name = (*js.Ident)(&file.Name.Name)

	for _, decl := range file.Decls {
		Decl(&fCtx, decl)
	}

	fCtx.Module.Print(&fCtx)

	got := buf.String()
	if got != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, got)
	}
}

func TestTranspiler(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "IfElse",
			src: `
package main

import "fmt"

func main() {
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	} else {
		fmt.Println("x is not greater than 5")
	}

	if x > 20 {
		fmt.Println("x is greater than 20")
	} else if x > 15 {
		fmt.Println("x is greater than 15")
	} else {
		fmt.Println("x is small")
	}
}
`,
			expected: `function main() {
    let x = 10
    if (x>5) {
        fmt.Println("x is greater than 5")
    }
 else  {
        fmt.Println("x is not greater than 5")
    }

    if (x>20) {
        fmt.Println("x is greater than 20")
    }
 else if (x>15) {
        fmt.Println("x is greater than 15")
    }
 else  {
        fmt.Println("x is small")
    }

}
`,
		},
		{
			name: "Variables",
			src: `
package main

func main() {
	var a int
	b := 2
}
`,
			expected: `function main() {
    let a = {  }
    let b = 2
}
`,
		},
		{
			name: "Loops",
			src: `
package main

func main() {
	for i := 0; i < 10; i++ {
	}

	x := 0
	for x < 5 {
		x++
	}
}
`,
			expected: `function main() {
    for (let i = 0;i<10;++i) {
    }

    let x = 0
    while (x<5) {
        ++x
    }

}
`,
		},
		{
			name: "Arrays",
			src: `
package main

func main() {
	arr := []int{1, 2, 3}
	x := arr[0]
	arr = append(arr, 4)
	l := len(arr)
}
`,
			expected: `function main() {
    let arr = [ 1, 2, 3 ]
    let x = arr[0]
    arr = [ ...(arr??[  ]), 4 ]
    let l = (arr?.length??0)
}
`,
		},
		{
			name: "Functions",
			src: `
package main

func add(a, b int) int {
	return a + b
}

func main() {
	res := add(1, 2)
}
`,
			expected: `function add(a, b) {
    return a+b
}

function main() {
    let res = add(1, 2)
}
`,
		},
		{
			name: "Structs",
			src: `
package main

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{Name: "Alice", Age: 30}
	n := p.Name
}
`,
			expected: `class Person {
    Name=null
    Age=null

    constructor(name, age) {
        this.Name = name
        this.Age = age
    }
}

function main() {
    let p = new Person("Alice", 30)
    let n = p.Name
}
`,
		},
		{
			name: "Range",
			src: `
package main

func main() {
	arr := []int{1, 2, 3}
	for i, v := range arr {
		_ = i
		_ = v
	}
	for _, v := range arr {
		_ = v
	}
	for i := range arr {
		_ = i
	}
}
`,
			expected: `function main() {
    let arr = [ 1, 2, 3 ]
    for (let [i, v] of arr.entries()) {
        _ = i
        _ = v
    }

    for (const v of arr) {
        _ = v
    }

    for (const i in arr) {
        _ = i
    }

}
`,
		},
		{
			name: "Pointers",
			src: `
package main

func main() {
	x := 10
	p := &x
	y := *p
}
`,
			expected: `function main() {
    let x = 10
    let p = x
    let y = p
}
`,
		},
		{
			name: "TypeAssertion",
			src: `
package main

func main() {
	var i interface{} = 10
	s := i.(int)
}
`,
			expected: `function main() {
    let i = {  }
    let s = i
}
`,
		},
		{
			name: "FuncLit",
			src: `
package main

func main() {
	f := func(x int) int {
		return x * x
	}
}
`,
			expected: `function main() {
    let f = function (x) {
        return x*x
    }
}
`,
		},
		{
			name: "SliceExpr",
			src: `
package main

func main() {
	arr := []int{1, 2, 3, 4, 5}
	s1 := arr[1:3]
	s2 := arr[:2]
	s3 := arr[2:]
	s4 := arr[:]
}
`,
			expected: `function main() {
    let arr = [ 1, 2, 3, 4, 5 ]
    let s1 = arr.slice(1, 3)
    let s2 = arr.slice(0, 2)
    let s3 = arr.slice(2)
    let s4 = arr.slice()
}
`,
		},
		{
			name: "IfInit",
			src: `
package main

func main() {
	if x := 10; x > 5 {
		println(x)
	}
}
`,
			expected: `function main() {
     {
        let x = 10
        if (x>5) {
            println(x)
        }

    }

}
`,
		},
		{
			name: "SwitchBasic",
			src: `
package main

func main() {
	x := 10
	switch x {
	case 1:
		println("one")
	case 2, 3:
		println("two or three")
	default:
		println("other")
	}
}
`,
			expected: `function main() {
    let x = 10
    switch (x) {
        case 1:
            println("one")
            break
        case 2:
        case 3:
            println("two or three")
            break
        default:
            println("other")
            break
    }
}
`,
		},
		{
			name: "SwitchNoTag",
			src: `
package main

func main() {
	x := 10
	switch {
	case x > 5:
		println("greater than 5")
	default:
		println("not greater")
	}
}
`,
			expected: `function main() {
    let x = 10
    switch (true) {
        case x>5:
            println("greater than 5")
            break
        default:
            println("not greater")
            break
    }
}
`,
		},
		{
			name: "SwitchInit",
			src: `
package main

func main() {
	switch x := 10; x {
	case 10:
		println(10)
	}
}
`,
			expected: `function main() {
     {
        let x = 10
        switch (x) {
            case 10:
                println(10)
                break
        }
    }

}
`,
		},
		{
			name: "SwitchFallthrough",
			src: `
package main

func main() {
	x := 1
	switch x {
	case 1:
		println("one")
		fallthrough
	case 2:
		println("two")
	}
}
`,
			expected: `function main() {
    let x = 1
    switch (x) {
        case 1:
            println("one")
        case 2:
            println("two")
            break
    }
}
`,
		},
		{
			name: "BreakContinue",
			src: `
package main

func main() {
	for i := 0; i < 10; i++ {
		if i == 5 {
			break
		}
		if i == 2 {
			continue
		}
	}
}
`,
			expected: `function main() {
    for (let i = 0;i<10;++i) {
        if (i==5) {
            break
        }

        if (i==2) {
            continue
        }

    }

}
`,
		},
		{
			name: "Copy",
			src: `
package main

func main() {
	a := []int{1, 2, 3}
	b := make([]int, 3)
	n := copy(b, a)
	copy(b, []int{4, 5})
}
`,
			expected: `function main() {
    let a = [ 1, 2, 3 ]
    let b = make(__U__, 3)
    let n = runtime.copy(b, a)
    runtime.copy(b, [ 4, 5 ])
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTranspilerTest(t, tt.name, tt.src, tt.expected)
		})
	}
}
