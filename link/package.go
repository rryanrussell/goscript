package link

import (
	_ "embed"
	"io"

	"github.com/rryanrussell/goscript/runtime"
)

//go:embed fmt.js
var fmtJS []byte

func Include(writer io.Writer) {
	// Include fmt.js
	writer.Write(fmtJS)
	writer.Write([]byte("\n"))

	// Include runtime/runtime.js
	writer.Write(runtime.RuntimeJS)
	writer.Write([]byte("\n"))

	// Include runtime/strings.js
	writer.Write(runtime.StringsJS)
	writer.Write([]byte("\n"))
}
