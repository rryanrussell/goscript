package parser

import (
	"fmt"
	"strings"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

type debugMsg struct {
	depth    int
	ok       bool
	node     any
	name     string
	parentId int

	timein  uint64
	timeout uint64

	optfail string
	failson bool

	logs []string

	matchNode  Node
	matchToken ast.Token
}

func (msg debugMsg) String() string {
	indent := strings.Repeat("  ", msg.depth)
	color := ast.Red
	txt := "reject"
	if msg.ok {
		color = ast.Green
		txt = "accept"
	}
	if len(msg.optfail) > 0 {
		color = ast.Gray
		txt = msg.optfail
	}

	nameAndSpace := msg.name

	if len(nameAndSpace) > 0 {
		nameAndSpace = nameAndSpace + " "
	}

	str := fmt.Sprintf("%s%s%v %s%s%s", indent, nameAndSpace, msg.node, color, txt, ast.Reset)

	if msg.matchNode != nil {
		str = str + fmt.Sprintf(" %+v <- %s", msg.matchNode, msg.matchToken)
	}

	str = str + fmt.Sprintf(" %d-%d", msg.timein, msg.timeout)

	for _, log := range msg.logs {
		str = str + " " + log
	}

	return str
}

type Debug struct {
	msgs      []debugMsg
	msgIds    Stack[int]
	nextMsgId int
	time      uint64
	idOffset  int
}

var DEBUG_MAX int

const (
	MAX_DEPTH    = 8
	MIN_COLLAPSE = 5
)

func (d *Debug) Enter(node any, name string) int {
	depth := d.msgIds.Len()

	if depth >= MAX_DEPTH {
		d.msgIds.Push(-1)
		return -1
	}

	if len(d.msgs) >= 2*DEBUG_MAX && d.msgIds.Len() > 0 {
		d.idOffset += DEBUG_MAX
		d.msgs = d.msgs[DEBUG_MAX:]
	}

	d.time++
	parentId, _ := d.msgIds.Peek()
	d.msgs = append(d.msgs, debugMsg{
		depth:    d.msgIds.Len(),
		name:     name,
		node:     node,
		parentId: parentId,
		timein:   d.time,
	})
	msgId := d.nextMsgId
	d.nextMsgId++
	d.msgIds.Push(msgId)
	return msgId
}

func (d *Debug) EventCount() uint64 {
	return d.time
}

func (d *Debug) Depth() int {
	return d.msgIds.Len()
}

func (d *Debug) Exit(ok bool) {
	d.time++
	id, _ := d.msgIds.Pop()

	id -= d.idOffset

	if id < 0 {
		return
	}

	d.msgs[id].ok = ok
	d.msgs[id].timeout = d.time
}

func (d *Debug) Match(match Node, tok ast.Token) {
	id, _ := d.msgIds.Peek()
	id -= d.idOffset

	if id < 0 {
		return
	}

	d.msgs[id].matchNode = match
	d.msgs[id].matchToken = tok
}

func (d *Debug) Optfail(optfail string) {
	id, _ := d.msgIds.Peek()
	id -= d.idOffset

	if id < 0 {
		return
	}
	d.msgs[id].optfail = optfail
}

func (d *Debug) Failsons(ids []int) {
	for _, id := range ids {
		id -= d.idOffset

		if id < 0 {
			continue
		}

		d.msgs[id].failson = true
	}
}

// func debugLog(c Ctx, event int, msg string) {
// 	event -= idOffset

// 	if event < 0 {
// 		return
// 	}

// 	msgs[event].logs = append(msgs[event].logs, msg)
// }

func (d *Debug) Dump() {
	msgSlice := make([]debugMsg, 0, DEBUG_MAX)

	for i := 0; i < len(d.msgs) && len(msgSlice) < DEBUG_MAX; i++ {
		// var sonskip int
		// for ; msgs[len(msgs)-i-1].failson; i++ {
		// 	sonskip++
		// }
		// if sonskip > 0 {
		// 	defer println("sonskip", sonskip)
		// }
		msgSlice = append(msgSlice, d.msgs[len(d.msgs)-i-1])
	}

	failson := -1
	var skipcount int

	for i := 0; i < len(msgSlice); i++ {
		msg := msgSlice[len(msgSlice)-i-1]
		color := ast.Reset

		if failson >= 0 && msg.depth <= failson {
			failson = -1
		}

		if msg.failson && (failson == -1 || msg.depth < failson) && msg.depth > 4 {
			failson = msg.depth
		}

		if failson >= 0 {
			if msg.depth > failson {
				skipcount++
				continue
			}

			color = ast.Gray
		}

		if skipcount > 0 {
			prefixLength := len(fmt.Sprintf("%d %d", msg.timein, msg.depth))
			indent := strings.Repeat(".", prefixLength)
			msgIndent := strings.Repeat("  ", msg.depth)
			fmt.Printf("%s%s  %s%d failsons skipped %s\n", ast.Gray, indent, msgIndent, skipcount, ast.Reset)
			skipcount = 0
		}

		fmt.Printf("%s%d %d  %s%s\n", color, msg.timein, msg.depth, msg, ast.Reset)
	}
}
