package interpreter

import (
	"time"

	"github.com/taki-mekhalfa/golox/ast"
)

var (
	// clockFn is a built-in function that returns the current time in unix seconds
	clockFn = &clock{}
)

type callable interface {
	call(*Interpreter, []interface{}) interface{}
	// returns the function's arity
	arity() int
}

type clock struct{}

func (c *clock) arity() int { return 0 }

func (c *clock) call(interpreter *Interpreter, args []interface{}) interface{} {
	return float64(time.Now().Unix())
}

type return_ struct {
	value interface{}
}

type function struct {
	closure     *environment
	declaration *ast.Function
}

func (f *function) arity() int { return len(f.declaration.Params) }
func (f *function) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	// save the current interpreter environment
	previous := interpreter.env

	defer func() {
		// recover back the interpreter environement
		interpreter.env = previous

		err := recover()
		if err == nil {
			// if the function returns nothing or does not have a returns statement
			// we return <nil>
			return
		}
		if ret_, ok := err.(return_); ok {
			ret = ret_.value
			return
		}
		// if the error is not nil or a return value, panic again to not silently hide the error
		panic(err)
	}()

	// create a new environement exclusive to this function call starting-up from the global env
	functionEnv := newEnvironment(f.closure)
	interpreter.env = functionEnv

	// bind function parameters to arguments
	for i, param := range f.declaration.Params {
		functionEnv.define(param.Lexeme, args[i])
	}

	// exectue the function body.
	// we don't use a block as a function body because a block will create
	// a new scope which is not what we want.
	// we want the arguments to be in the same scope as the function body.
	for _, stmt := range f.declaration.Body {
		stmt.Accept(interpreter)
	}

	return
}
