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
	if fmtJS != nil {
		writer.Write(fmtJS)
	}

	// Include runtime/strings.js
	if runtime.StringsJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.StringsJS)
	}

	// Include runtime/runtime.js
	if runtime.RuntimeJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.RuntimeJS)
	}
}
