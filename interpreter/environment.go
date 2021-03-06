package interpreter

import "github.com/taki-mekhalfa/golox/ast"

func (i *Interpreter) Resolve(expr ast.Expr, distance int) {
	i.scopeDists[expr] = distance
}

func (i *Interpreter) lookUp(expr ast.Expr, name string) (interface{}, bool) {
	dist, ok := i.scopeDists[expr]
	if ok {
		env := i.env
		for i := 0; i < dist; i++ {
			env = env.parent
		}
		return env.get(name)
	}
	return i.globals.get(name)
}

type environment struct {
	values map[string]interface{}
	parent *environment
}

func newEnvironment(parent *environment) *environment {
	return &environment{
		values: map[string]interface{}{},
		parent: parent,
	}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) assign(name string, value interface{}) {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return
	}
	e.parent.assign(name, value)
}

func (e *environment) get(name string) (interface{}, bool) {
	v, ok := e.values[name]
	if !ok && e.parent != nil {
		// check the parent environement if the symbol is not defined in the current one
		return e.parent.get(name)
	}
	return v, ok
}
