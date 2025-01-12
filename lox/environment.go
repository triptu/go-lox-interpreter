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

func (e *environment) get(name string) (any, error) {
	if val, ok := e.vars[name]; ok {
		return val, nil
	} else if e.outer != nil {
		return e.outer.get(name)
	} else {
		return nil, errors.New("not found")
	}
}

func (e *environment) set(name string, val any) {
	e.vars[name] = val
}
