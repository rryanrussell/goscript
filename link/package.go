package link

import (
	"io"
	"os"
)

func Include(writer io.Writer) {
	data, _ := os.ReadFile("goscript/link/fmt.js")
	writer.Write(data)
}
