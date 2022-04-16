package interpreter

import (
	"time"
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
