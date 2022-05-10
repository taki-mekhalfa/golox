package interpreter

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/token"
)

const (
	init_ = "init"
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

func (c *class) arity() int {
	if initializer, ok := c.methods[init_]; ok {
		return initializer.arity()
	}
	return 0
}

func (c *class) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	instance := newInstance(c)

	// check if the user did provide an initializer,
	// if so, call it before returning the instance.
	if initializer, ok := c.methods[init_]; ok {
		method := &function{declaration: initializer.declaration, closure: newEnvironment(initializer.closure)}
		method.closure.define("this", instance)
		method.call(interpreter, args)
	}

	// we always return the instance, even if the user had a 'return;'
	// in the init functions (which would mean a nil).
	// we don't allow returning a value inside the initializer.
	return instance
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
	// don't allow code to access the `init` function
	if method, ok := ins.klass.methods[t.Lexeme]; t.Lexeme != init_ && ok {
		method := &function{declaration: method.declaration, closure: newEnvironment(method.closure)}
		method.closure.define("this", ins)
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
