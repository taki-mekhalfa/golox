package interpreter

type class struct {
	name string
}

// String implements fmt.Stringer
func (c *class) String() string {
	return c.name + " class"
}

func (c *class) arity() int { return 0 }
func (c *class) call(interpreter *Interpreter, args []interface{}) (ret interface{}) {
	return &instance{klass: c}
}

type instance struct {
	klass *class
}

// String implements fmt.Stringer
func (ins *instance) String() string {
	return ins.klass.name + " instance"
}
