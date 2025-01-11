package lox

/*
This file contains a lot of struct which are data types to represent different
expressions in our language.
*/

type expr[R any] interface {
	accept(exprVisitor[R]) R
}

/*
action classes implement this interface defining how they would process each type of expression
*/
type exprVisitor[R any] interface {
	visitAssign(eAssign[R]) R
	visitBinary(eBinary[R]) R
	visitCall(eCall[R]) R
	visitGet(eGet[R]) R
	visitGrouping(eGrouping[R]) R
	visitLiteral(eLiteral[R]) R
	visitLogical(eLogical[R]) R
	visitSet(eSet[R]) R
	visitSuper(eSuper[R]) R
	visitThis(eThis[R]) R
	visitUnary(eUnary[R]) R
	visitVariable(eVariable[R]) R
}

type eAssign[R any] struct {
	name  token
	value expr[R]
}

type eBinary[R any] struct {
	left     expr[R]
	operator token
	right    expr[R]
}

type eCall[R any] struct {
	callee    expr[R]
	paren     token
	arguments []expr[R]
}

type eGet[R any] struct {
	object expr[R]
	name   token
}

type eGrouping[R any] struct {
	expression expr[R]
}

type eLiteral[R any] struct {
	value interface{}
}

type eLogical[R any] struct {
	left     expr[R]
	operator token
	right    expr[R]
}

type eSet[R any] struct {
	object expr[R]
	name   token
	value  expr[R]
}

type eSuper[R any] struct {
	keyword token
	method  token
}

type eThis[R any] struct {
	keyword token
}

type eUnary[R any] struct {
	operator token
	right    expr[R]
}

// variable access expression
type eVariable[R any] struct {
	name token
}

// define accept methods for each type of expression

func (e eAssign[R]) accept(v exprVisitor[R]) R {
	return v.visitAssign(e)
}

func (e eBinary[R]) accept(v exprVisitor[R]) R {
	return v.visitBinary(e)
}

func (e eCall[R]) accept(v exprVisitor[R]) R {
	return v.visitCall(e)
}

func (e eGet[R]) accept(v exprVisitor[R]) R {
	return v.visitGet(e)
}

func (e eGrouping[R]) accept(v exprVisitor[R]) R {
	return v.visitGrouping(e)
}

func (e eLiteral[R]) accept(v exprVisitor[R]) R {
	return v.visitLiteral(e)
}

func (e eLogical[R]) accept(v exprVisitor[R]) R {
	return v.visitLogical(e)
}

func (e eSet[R]) accept(v exprVisitor[R]) R {
	return v.visitSet(e)
}

func (e eSuper[R]) accept(v exprVisitor[R]) R {
	return v.visitSuper(e)
}

func (e eThis[R]) accept(v exprVisitor[R]) R {
	return v.visitThis(e)
}

func (e eUnary[R]) accept(v exprVisitor[R]) R {
	return v.visitUnary(e)
}

func (e eVariable[R]) accept(v exprVisitor[R]) R {
	return v.visitVariable(e)
}
