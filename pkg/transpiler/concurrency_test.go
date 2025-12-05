package transpiler

import (
	"testing"
)

func TestConcurrency(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name: "GoStmt",
			src: `
package main

func main() {
	go func() {
		println("hello")
	}()
}
`,
			expected: `function main() {
  runtime.go(() => {
    (function () {
      println("hello");
    })();
  });
}
`,
		},
		{
			name: "ChannelSendRecv",
			src: `
package main

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	x := <-ch
}
`,
			expected: `async function main() {
  let ch = runtime.makeChan();
  runtime.go(() => {
    (async function () {
      await ch.send(1);
    })();
  });
  let x = await ch.recv();
}
`,
		},
		{
			name: "SelectStmt",
			src: `
package main

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	
	select {
	case x := <-c1:
		println(x)
	case c2 <- 1:
		println("sent")
	default:
		println("default")
	}
}
`,
			expected: `async function main() {
    let c1 = runtime.makeChan()
    let c2 = runtime.makeChan()
    await runtime.select([ { op: 'recv', chan: c1, handler: (_val) =>  {
        let x = _val
        println(x)
    }
 }, { op: 'send', chan: c2, val: 1, handler: () =>  {
        println("sent")
    }
 }, { op: 'default', handler: () =>  {
        println("default")
    }
 } ])
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
