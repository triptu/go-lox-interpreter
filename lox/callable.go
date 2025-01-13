package lox

import (
	"fmt"
	"time"
)

type callable interface {
	arity() int // number of arguments needed
	call(interpreter interpreter, arguments []any) (any, error)
}

const nativeFnStr = "<native fn>"

// define some native functions
type nativeClock struct{}
type nativePrint struct{}
type loxFunction struct {
	declaration sFunction
}

type returnAsError struct {
	value any
}

func (r returnAsError) Error() string {
	return fmt.Sprintf("return statement with value %v", r.value)
}

var _ callable = nativeClock{} // assert interface adherence
var _ callable = nativePrint{} // assert interface adherence
var _ callable = loxFunction{} // assert interface adherence

func (n nativeClock) arity() int {
	return 0
}

func (n nativeClock) call(i interpreter, arguments []any) (any, error) {
	timeInt := time.Now().UnixMilli() / 1000
	return float64(timeInt), nil
}

func (n nativeClock) String() string {
	return nativeFnStr
}

func (n nativePrint) arity() int {
	return 1
}

func (n nativePrint) call(i interpreter, arguments []any) (any, error) {
	fmt.Println(arguments[0])
	return nil, nil
}

func (n nativePrint) String() string {
	return nativeFnStr
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
