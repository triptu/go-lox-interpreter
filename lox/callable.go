package lox

import (
	"fmt"
	"strconv"
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
	declaration   sFunction
	closure       *environment
	isInitializer bool // this is a constructor for some class
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

func (f loxFunction) arity() int {
	return len(f.declaration.parameters)
}

func (f loxFunction) bind(instance loxClassInstance) loxFunction {
	env := newChildEnvironment(f.closure)
	env.define("this", instance)
	return loxFunction{declaration: f.declaration, closure: env, isInitializer: f.isInitializer}
}

func (f loxFunction) call(i interpreter, arguments []any) (any, error) {
	env := newChildEnvironment(f.closure)
	for i, param := range f.declaration.parameters {
		env.define(param.lexeme, arguments[i])
	}

	// note that in the parsing stage, we've stored the function's body as
	// a list of statements, and not as a block.
	err := i.executeBlock(f.declaration.body, env)
	if err != nil {
		if returnValue, ok := err.(returnAsError); ok {
			if f.isInitializer { // return from constructor always returns class instance
				return f.closure.getAt(0, "this")
			}
			return returnValue.value, nil
		}
		return nil, err
	}
	if f.isInitializer {
		return f.closure.getAt(0, "this")
	}
	return nil, nil
}

func (f loxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
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
			println(a[0])
			return nil, nil
		},
	})
	globals.define("input", nativeFunction{
		fn: func(i interpreter, a []any) (any, error) {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return nil, err
			}
			return input, nil
		},
	})
	globals.define("parseNumber", nativeFunction{
		arityCnt: 1,
		fn: func(i interpreter, a []any) (any, error) {
			return strconv.ParseFloat(a[0].(string), 64)
		},
	})
}
