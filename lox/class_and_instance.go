package lox

type loxClass struct {
	name        string
	declaration sClass
	closure     *environment
}

type loxClassInstance struct {
	klass  loxClass
	fields map[string]any
}

var _ callable = loxClass{} // assert interface adherence

func (c loxClass) String() string {
	return c.name
}

func (c loxClass) arity() int {
	return 0
}

/*
calling a class instntiates it, and returns an instance of it
*/
func (c loxClass) call(i interpreter, arguments []any) (any, error) {
	instance := loxClassInstance{klass: c, fields: make(map[string]any)}
	return instance, nil
}

func (i loxClassInstance) String() string {
	return i.klass.name + " instance"
}

func (i loxClassInstance) get(name token) any {
	val, ok := i.fields[name.lexeme]
	if !ok {
		logRuntimeError(name.line, "Undefined property '"+name.lexeme+"'.")
	}
	return val
}

func (i loxClassInstance) set(name token, val any) any {
	i.fields[name.lexeme] = val
	return val
}
