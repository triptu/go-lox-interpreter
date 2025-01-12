package lox

type callable interface {
	arity() int // number of arguments needed
	call(interpreter interpreter, arguments []any) any
}
