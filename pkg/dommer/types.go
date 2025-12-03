package dommer

import (
	"time"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

type Result struct {
	Value      *ast.DtsModule
	LineStart  int
	LineEnd    int
	Duration   time.Duration
	Ok         bool
	EventCount uint64
}
