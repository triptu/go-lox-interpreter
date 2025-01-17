package lox

import "errors"

/*
for the internal state of the interpreter, we need to keep track of the
variables declared in the program. This is done by creating an environment
which is a map of variables to their values.
This is recursive in nature following the scope chain.
*/

type environment struct {
	outer *environment
	vars  map[string]any
}

func newEnvironment() *environment {
	return &environment{
		vars: make(map[string]any),
	}
}

func newChildEnvironment(outer *environment) *environment {
	return &environment{
		outer: outer,
		vars:  make(map[string]any),
	}
}

// when the variable is defined for the first time in current scope
func (e *environment) define(name string, val any) {
	e.vars[name] = val
}

func (e *environment) get(name string) (any, error) {
	if val, ok := e.vars[name]; ok {
		return val, nil
	} else if e.outer != nil {
		return e.outer.get(name)
	} else {
		return nil, errors.New("not found")
	}
}

func (e *environment) set(name string, val any) error {
	if _, ok := e.vars[name]; ok {
		e.vars[name] = val
		return nil
	} else if e.outer != nil {
		return e.outer.set(name, val)
	} else {
		return errors.New("not found")
	}
}

func (e *environment) getAt(depth int, name string) (any, error) {
	return e.ancestor(depth).get(name)
}

func (e *environment) setAt(depth int, name string, val any) error {
	return e.ancestor(depth).set(name, val)
}

func (e *environment) ancestor(depth int) *environment {
	env := e
	for depth > 0 {
		env = env.outer
		depth--
	}
	return env
}
