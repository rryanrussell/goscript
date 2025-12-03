package link

import (
	"fmt"
	"io"
	"os"
)

func Include(writer io.Writer) {
	// Include fmt.js
	// Try link/fmt.js (from root) first
	data, err := os.ReadFile("link/fmt.js")
	if err != nil {
		// Fallback to goscript/link/fmt.js
		data, err = os.ReadFile("goscript/link/fmt.js")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to include fmt.js: %v\n", err)
		}
	}
	if data != nil {
		writer.Write(data)
	}

	// Include runtime/strings.js
	data, err = os.ReadFile("runtime/strings.js")
	if err == nil {
		writer.Write([]byte("\n"))
		writer.Write(data)
	} else {
		// Try goscript/runtime/strings.js fallback
		data, err = os.ReadFile("goscript/runtime/strings.js")
		if err == nil {
			writer.Write([]byte("\n"))
			writer.Write(data)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: failed to include runtime/strings.js: %v\n", err)
		}
	}
}
