package js

import "fmt"

// SetType generates a runtime call to set the type of an object.
// runtime.SetType(obj, "typeName")
func SetType(obj Expr, typeName string) Stmt {
	return &ExprStmt{
		X: &Call{
			Func: &Selector{
				X:   IdentP("runtime"),
				Sel: IdentP("SetType"),
			},
			Args: []Expr{
				obj,
				&BasicLit{Value: fmt.Sprintf("\"%s\"", typeName)},
			},
		},
	}
}

// Panic generates a runtime call to panic with a message.
// runtime.panic("msg")
func Panic(msg string) Stmt {
	return &ExprStmt{
		X: &Call{
			Func: &Selector{
				X:   IdentP("runtime"),
				Sel: IdentP("panic"),
			},
			Args: []Expr{
				&BasicLit{Value: fmt.Sprintf("\"%s\"", msg)},
			},
		},
	}
}
