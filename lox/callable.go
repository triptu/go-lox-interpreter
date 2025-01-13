package lox

import (
	"fmt"
)

type callable interface {
	arity() int // number of arguments needed
	call(interpreter interpreter, arguments []any) (any, error)
	String() string
}

type nativeFunction struct {
	arityCnt int
	fn       func(interpreter, []any) (any, error)
}

type loxFunction struct {
	declaration sFunction
}

type returnAsError struct {
	value any
}

func (r returnAsError) Error() string {
	return fmt.Sprintf("return statement with value %v", r.value)
}

var _ callable = nativeFunction{} // assert interface adherence
var _ callable = loxFunction{}    // assert interface adherence

func (n nativeFunction) arity() int {
	return n.arityCnt
}

func (n nativeFunction) call(i interpreter, arguments []any) (any, error) {
	return n.fn(i, arguments)
}

func (n nativeFunction) String() string {
	return "<native function>"
}

func (n loxFunction) arity() int {
	return len(n.declaration.parameters)
}

func (n loxFunction) call(i interpreter, arguments []any) (any, error) {
	env := newChildEnvironment(i.globals)
	for i, param := range n.declaration.parameters {
		env.define(param.lexeme, arguments[i])
	}

	err := i.executeBlock(n.declaration.body, env)
	if err != nil {
		if _, ok := err.(returnAsError); ok {
			return err.(returnAsError).value, nil
		}
		return nil, err
	}
	return nil, nil
}

func (n loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", n.declaration.name.lexeme)
}
