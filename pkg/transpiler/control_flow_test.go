package transpiler

import "testing"

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "BasicDefer",
			src: `
package main

import "fmt"

func main() {
	defer fmt.Println("deferred")
	fmt.Println("normal")
}
`,
			expected: `function main() {
    let __gs_deferred = [  ]
    try  {
        __gs_deferred.push(() =>  {
            runtime.fmt.println("deferred")
        }
)
        runtime.fmt.println("normal")
    }
 finally  {
        runtime.runDefers(__gs_deferred)
    }

}
`,
		},
		{
			name: "PanicRecover",
			src: `
package main

func main() {
	defer func() {
		if r := recover(); r != nil {
			println("recovered", r)
		}
	}()
	panic("oops")
}
`,
			expected: `function main() {
    let __gs_deferred = [  ]
    try  {
        __gs_deferred.push(() =>  {
            (function () {
                let __gs_panicValue = null
                try  {
                     {
                        let r = runtime.recover()
                        if (r!=nil) {
                            println("recovered", r)
                        }

                    }

                }
 catch (__e)  {
                    runtime.__gs_panicValue = __e
                }

            })()
        }
)
        runtime.panic("oops")
    }
 finally  {
        runtime.runDefers(__gs_deferred)
    }

}
`,
		},
		{
			name: "TypeSwitchInDefer",
			src: `
package main

func main() {
	var x interface{} = 10
	defer func() {
		switch v := x.(type) {
		case int:
			println(v)
		}
	}()
}
`,
			expected: `function main() {
    let __gs_deferred = [  ]
    try  {
        let x = 10
        __gs_deferred.push(() =>  {
            (function () {
                 {
                    let _tag = x
                    if (typeof _tag=='number') {
                        let v = _tag
                        println(v)
                    }

                }

            })()
        }
)
    }
 finally  {
        runtime.runDefers(__gs_deferred)
    }

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
