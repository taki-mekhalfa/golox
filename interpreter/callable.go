package interpreter

import (
	"time"

	"github.com/taki-mekhalfa/golox/ast"
)

var (
	clockFn = &clock{}
)

type callable interface {
	call(*Interpreter, []interface{}) interface{}
	// returns the function's arity
	arity() int
}

type clock struct {
}

func (c *clock) arity() int { return 0 }

func (c *clock) call(interpreter *Interpreter, args []interface{}) interface{} {
	return time.Now().Unix()
}

type function struct {
	declaration *ast.Function
}

func (f *function) arity() int { return len(f.declaration.Params) }
func (f *function) call(interpreter *Interpreter, args []interface{}) interface{} {
	// save the current interpreter environment
	previous := interpreter.env

	// create a new environement exclusive to this function call starting-up from the global env
	functionEnv := newEnvironment(interpreter.globals)
	interpreter.env = functionEnv

	// bind function parameters to arguments
	for i, param := range f.declaration.Params {
		functionEnv.define(param.Lexeme, args[i])
	}

	// exectue the function body
	interpreter.Interpret(f.declaration.Body)

	// recover back the interpreter environement
	interpreter.env = previous

	return nil
}
