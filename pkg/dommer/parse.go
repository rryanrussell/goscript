package dommer

import (
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
	"github.com/rryanrussell/goscript/pkg/dommer/parser"
)

func Parse(dataBytes []byte, debugMax int) []Result {
	parser.DEBUG_MAX = debugMax

	r, err := regexp.Compile("(?m)^(declare|interface)")

	//r, err := regexp.Compile(`(?m)^//start`)

	if err != nil {
		panic(err)
	}

	newlineR, err := regexp.Compile("\n")

	if err != nil {
		panic(err)
	}

	matches := r.FindAllIndex(dataBytes, -1)

	var jobs [16][2]int
	var tasksPerJob = len(matches) / len(jobs)

	for i := 0; i < len(jobs)-1; i++ {
		jobs[i][0] = matches[tasksPerJob*i][0]
		jobs[i][1] = matches[tasksPerJob*(i+1)][0]
	}

	jobs[len(jobs)-1][0] = matches[tasksPerJob*(len(jobs)-1)][0]
	jobs[len(jobs)-1][1] = len(dataBytes)

	start := time.Now()

	results := make([]Result, len(jobs))
	var wg sync.WaitGroup

	node := parser.NewParser()

	for i := 0; i < len(jobs); i++ {
		wg.Add(1)
		go func(id, start, end int) {
			defer wg.Done()
			var result Result

			result.LineStart = len(newlineR.FindAllIndex(dataBytes[:start], -1))

			startTime := time.Now()
			scanner := parser.NewTokenizer(dataBytes[start:end])

			ctx := parser.NewCtx(scanner)
			ids, ok := node.Process(ctx)

			if ids != nil {
				result.Value = ids.(*ast.DtsModule)
			}

			result.Ok = ok
			result.Duration = time.Since(startTime)
			result.EventCount = ctx.Debug().EventCount()
			result.LineEnd = len(newlineR.FindAllIndex(dataBytes[:start+int(scanner.ReadPosition())], -1))

			if !result.Ok || len(result.Value.Decls) == 0 {
				ctx.Debug().Dump()
				panic(fmt.Errorf("failure encountered at line %d", result.LineStart))
			}

			results[id] = result
		}(i, jobs[i][0], jobs[i][1])
	}

	wg.Wait()

	endProc := time.Now()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration > results[j].Duration
	})

	var features int

	for i := 0; i < len(results); i++ {
		var items int
		r := results[i]

		if r.Value != nil {
			items = len(r.Value.Decls)
		}

		features += items

		if i < 10 || !r.Ok || items == 0 {
			color := ast.Green

			if !r.Ok || items == 0 {
				color = ast.Red
			}

			var mspl, epl float64

			if r.LineEnd != r.LineStart {
				mspl = float64(r.Duration.Milliseconds()) / float64(r.LineEnd-r.LineStart)
				epl = float64(r.EventCount) / float64(r.LineEnd-r.LineStart)
			}

			fmt.Printf("%dms %s%d items%s lines %d-%d (%.2f ms/l) (%.2f ev/l)\n", r.Duration.Milliseconds(), color, items, ast.Reset, r.LineStart, r.LineEnd, mspl, epl)
		}
	}

	fmt.Printf("total: %dms, tasks: %d, features: %d\n", endProc.Sub(start).Milliseconds(), len(results), features)

	return results
}
