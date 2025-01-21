package lox

type loxClass struct {
	name       string
	methods    map[string]loxFunction
	superclass *loxClass
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
	initializer, ok := c.findMethod("init")
	if ok {
		return initializer.arity()
	}
	return 0
}

/*
calling a class instntiates it, and returns an instance of it
*/
func (c loxClass) call(i interpreter, arguments []any) (any, error) {
	instance := loxClassInstance{klass: c, fields: make(map[string]any)}
	initializer, ok := c.findMethod("init")
	if ok {
		// constructors are special, when the instance is created, they're automatically called
		// with the arguments passed to the class
		initializer.bind(instance).call(i, arguments)
	}
	return instance, nil
}

func (c loxClass) findMethod(name string) (loxFunction, bool) {
	method, ok := c.methods[name]
	if ok {
		return method, ok
	}
	if c.superclass == nil {
		return method, false
	}
	method, ok = c.superclass.findMethod(name)
	return method, ok
}

func (i loxClassInstance) String() string {
	return i.klass.name + " instance"
}

func (i loxClassInstance) get(name token) any {
	val, ok := i.fields[name.lexeme]
	if ok {
		return val
	}
	method, ok := i.klass.findMethod(name.lexeme)
	if ok {
		return method.bind(i)
	}

	logRuntimeError(name.line, "Undefined property '"+name.lexeme+"'.")
	return nil
}

func (i loxClassInstance) set(name token, val any) any {
	i.fields[name.lexeme] = val
	return val
}
