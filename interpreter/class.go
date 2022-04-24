package interpreter

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/token"
)

type class struct {
	name    string
	methods map[string]*function
}

func newClass(name string) *class {
	return &class{name: name, methods: make(map[string]*function)}
}

// String implements fmt.Stringer
func (c *class) String() string {
	return c.name + " class"
}

func (c *class) arity() int { return 0 }
func (c *class) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	return newInstance(c)
}

type instance struct {
	klass      *class
	properties map[string]interface{}
}

func newInstance(klass *class) *instance {
	return &instance{klass: klass, properties: make(map[string]interface{})}
}

// String implements fmt.Stringer
func (ins *instance) String() string {
	return ins.klass.name + " instance"
}

func (ins *instance) get(t token.Token) interface{} {
	// properties shadow methods.
	// this is a subtle but important semantic point
	if property, ok := ins.properties[t.Lexeme]; ok {
		return property
	}
	if method, ok := ins.klass.methods[t.Lexeme]; ok {
		return method
	}
	panic(runtimeError{
		token: t,
		msg:   fmt.Sprintf("Undefined property '%s'.", t.Lexeme),
	})
}

func (ins *instance) set(t token.Token, value interface{}) {
	ins.properties[t.Lexeme] = value
}
