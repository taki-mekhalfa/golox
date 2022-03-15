package interpreter

type environment struct {
	values map[string]interface{}
	parent *environment
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
		return e.parent.get(name)
	}
	return v, ok
}
