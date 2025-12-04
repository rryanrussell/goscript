package transpiler

import "testing"

func TestTypeSystem(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "TypeSwitch",
			src: `
package main

type Person struct {
    Name string
}

func main() {
    var x interface{} = Person{Name: "Bob"}
    switch v := x.(type) {
    case int:
        println("int", v)
    case string:
        println("string", v)
    case Person:
        println("Person", v.Name)
    default:
        println("unknown", v)
    }
}
`,
			expected: `class Person {
    Name=null

    constructor(name) {
        this.__goType = "main.Person"
        this.Name = name
    }
}

function main() {
    let x = new Person("Bob")
     {
        let _tag = x
        if (typeof _tag=='number') {
            let v = _tag
            println("int", v)
        }
 else if (typeof _tag=='string') {
            let v = _tag
            println("string", v)
        }
 else if (_tag?.__goType=="main.Person") {
            let v = _tag
            println("Person", v.Name)
        }
 else  {
            let v = _tag
            println("unknown", v)
        }

    }

}
`,
		},
		{
			name: "TypeAssertionPanic",
			src: `
package main

func main() {
	var x interface{} = 10
	i := x.(int)
}
`,
			expected: `function main() {
    let x = 10
    let i = (function () {
        let _v = x
        if (typeof _v=='number') {
            return _v
        }

        throw new Error("interface conversion: interface is not int")
    })()
}
`,
		},
		{
			name: "TypeAssertionCommaOk",
			src: `
package main

func main() {
	var x interface{} = 10
	i, ok := x.(int)
}
`,
			expected: `function main() {
    let x = 10
    let [ i, ok ] = (function () {
        let _v = x
        if (typeof _v=='number') {
            return [ _v, true ]
        }

        return [ null, false ]
    })()
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
