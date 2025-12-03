package parser

import (
	"fmt"
	"reflect"
)

type Node interface {
	Process(Ctx) (any, bool)
}

type Parallel struct {
	Paths []Node
}

type Visit struct {
	Position int
	Node     any
}

type VisitContext struct {
	visits map[Visit]uint64
}

func (p Parallel) Process(c Ctx) (any, bool) {
	var fails []int

	for _, path := range p.Paths {

		visitCtx := c.Visit()

		if visitCtx.visits == nil {
			visitCtx.visits = make(map[Visit]uint64)
		}

		visit := Visit{Position: c.Position(), Node: path}
		previousEvent, visited := visitCtx.visits[visit]

		if visited {
			c.Debug().Enter(path.(*Automaton).Goal, fmt.Sprintf("skip:%d", previousEvent))
			c.Debug().Exit(false)

			continue
		}

		event := c.Debug().Enter(path.(*Automaton).Goal, "")

		visitCtx.visits[visit] = c.Debug().time

		val, ok := try(c, path)

		if ok {
			c.Debug().Failsons(fails)
			c.Debug().Exit(true)
			delete(visitCtx.visits, visit)
			return val, true
		}

		fails = append(fails, event)

		c.Debug().Exit(false)

		delete(visitCtx.visits, visit)
	}

	return nil, false
}

type Optional struct {
	Node
}

func (op Optional) Process(c Ctx) (any, bool) {
	val, ok := try(c, op.Node)
	if !ok {
		c.Debug().Optfail("optfail")
		return nil, true
	}
	return val, true
}

type Repeat struct {
	PayNode             Node
	Delim               Node
	AcceptTrailingDelim bool
	Minimum             int
}

func (r *Repeat) Process(c Ctx) (any, bool) {
	var vals []any

	for {
		val, ok := try(c, r.PayNode)

		if ok {
			vals = append(vals, val)
		} else if len(vals) == 0 || r.Delim == nil || r.AcceptTrailingDelim {
			break
		} else {
			return nil, false
		}

		if r.Delim != nil && !test(c, r.Delim) {
			break
		}
	}

	if len(vals) == 0 {
		c.Debug().Optfail("zero")
	}

	if len(vals) >= r.Minimum {
		return vals, true
	}

	return nil, false
}

type FuncNode struct {
	Func func() (any, bool)
	Name string
}

func (f FuncNode) Hash() any {
	return f.Name
}

func (f FuncNode) Process(c Ctx) (any, bool) {
	return f.Func()
}

type Step struct {
	Node
	Target reflect.StructField
}

type Automaton struct {
	Goal     reflect.Type
	Sequence []Step
}

func (a *Automaton) Process(c Ctx) (any, bool) {
	var goal reflect.Value
	var init bool

	for _, step := range a.Sequence {
		debug := c.Debug().Depth() > 0

		if debug {
			c.Debug().Enter(step.Target.Type, step.Target.Name)
		}

		val, ok := try(c, step)

		if debug {
			c.Debug().Exit(ok)
		}

		if !ok {
			return nil, false
		}

		if val != nil {
			if !init {
				init = true
				goal = reflect.New(a.Goal)
			}

			field := goal.Elem().FieldByIndex(step.Target.Index)

			if field.Type().Kind() == reflect.Slice {
				valSlice := val.([]any)
				slice := reflect.MakeSlice(field.Type(), 0, len(valSlice))
				for _, value := range valSlice {
					slice = reflect.Append(slice, reflect.ValueOf(value))
				}
				val = slice.Interface()
			}

			field.Set(reflect.ValueOf(val))
		}
	}

	return goal.Interface(), true
}
