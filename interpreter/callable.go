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

type return_ interface{}

type function struct {
	declaration *ast.Function
}

func (f *function) arity() int { return len(f.declaration.Params) }
func (f *function) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	// save the current interpreter environment
	previous := interpreter.env

	defer func() {
		err := recover()

		var ok bool
		ret, ok = err.(return_)
		if err == nil || ok {
			// recover back the interpreter environement
			interpreter.env = previous
			return
		}
		// if the error is not nil or a return value, panic again to not silently hide the error
		panic(err)
	}()

	// create a new environement exclusive to this function call starting-up from the global env
	functionEnv := newEnvironment(interpreter.globals)
	interpreter.env = functionEnv

	// bind function parameters to arguments
	for i, param := range f.declaration.Params {
		functionEnv.define(param.Lexeme, args[i])
	}

	// exectue the function body
	for _, stmt := range f.declaration.Body {
		stmt.Accept(interpreter)
	}

	return
}
