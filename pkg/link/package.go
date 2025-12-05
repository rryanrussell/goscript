package link

import (
	_ "embed"
	"io"

	"github.com/rryanrussell/goscript/pkg/runtime"
)

//go:embed fmt.js
var fmtJS []byte

func Include(writer io.Writer) {
	// Include runtime/fmt.js
	if runtime.FmtJS != nil {
		writer.Write(runtime.FmtJS)
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

	// Include runtime/reflect.js
	if runtime.ReflectJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.ReflectJS)
	}

	// Include runtime/errors.js
	if runtime.ErrorsJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.ErrorsJS)
	}

	// Include runtime/strconv.js
	if runtime.StrconvJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.StrconvJS)
	}

	// Include runtime/math.js
	if runtime.MathJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.MathJS)
	}

	// Include runtime/time.js
	if runtime.TimeJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.TimeJS)
	}

	// Include runtime/sort.js
	if runtime.SortJS != nil {
		writer.Write([]byte("\n"))
		writer.Write(runtime.SortJS)
	}
}
