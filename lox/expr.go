package lox

/*
This file contains a lot of struct which are data types to represent different
expressions in our language.
*/

type expr interface{}

/*
action classes implement this interface defining how they would process each type of expression
*/
type exprVisitor[R any] interface {
	visitAssign(eAssign) R
	visitBinary(eBinary) R
	visitCall(eCall) R
	visitGet(eGet) R
	visitGrouping(eGrouping) R
	visitLiteral(eLiteral) R
	visitLogical(eLogical) R
	visitSet(eSet) R
	visitSuper(eSuper) R
	visitThis(eThis) R
	visitUnary(eUnary) R
	visitVariable(eVariable) R
}

type eAssign struct {
	name  token
	value expr
}

type eBinary struct {
	left     expr
	operator token
	right    expr
}

type eCall struct {
	callee    expr
	paren     token
	arguments []expr
}

type eGet struct {
	object expr
	name   token
}

type eGrouping struct {
	expression expr
}

type eLiteral struct {
	value interface{}
}

type eLogical struct {
	left     expr
	operator token
	right    expr
}

type eSet struct {
	object expr
	name   token
	value  expr
}

type eSuper struct {
	keyword token
	method  token
}

type eThis struct {
	keyword token
}

type eUnary struct {
	operator token
	right    expr
}

// variable access expression
type eVariable struct {
	name token
}
