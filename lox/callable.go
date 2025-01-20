package lox

import (
	"fmt"
	"time"
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
	closure     *environment
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
	return "<native fn>"
}

func (n loxFunction) arity() int {
	return len(n.declaration.parameters)
}

func (n loxFunction) call(i interpreter, arguments []any) (any, error) {
	env := newChildEnvironment(n.closure)
	for i, param := range n.declaration.parameters {
		env.define(param.lexeme, arguments[i])
	}

	// note that in the parsing stage, we've stored the function's body as
	// a list of statements, and not as a block.
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

/*
Theses are built-in functions that will be available natively in lox.
*/
func defineNativeFunctions(globals *environment) {
	globals.define("clock", nativeFunction{
		fn: func(i interpreter, a []any) (any, error) {
			timeInt := time.Now().UnixMilli() / 1000
			return float64(timeInt), nil
		},
	})
	globals.define("sleep", nativeFunction{ // sleep in seconds
		arityCnt: 1,
		fn: func(i interpreter, a []any) (any, error) {
			time.Sleep(time.Duration(a[0].(float64)) * time.Millisecond)
			return nil, nil
		},
	})
	globals.define("print", nativeFunction{
		arityCnt: 1,
		fn: func(i interpreter, a []any) (any, error) {
			fmt.Println(a[0])
			return nil, nil
		},
	})
}
