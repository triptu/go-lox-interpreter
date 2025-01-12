package lox

/*
This file contains a lot of struct which are data types to represent different
expressions in our language. These are basically the nodes of the
AST(abstract syntax tree) of our language.
*/

type expr[R any] interface {
	accept(exprVisitor[R]) R
}

/*
action classes implement this interface defining how they would process each type of expression
*/
type exprVisitor[R any] interface {
	// a = 3
	visitAssignExpr(eAssign[R]) R
	// 1 + 2, 5 * 6, 2 < 3, etc
	visitBinaryExpr(eBinary[R]) R
	// myFunction(1, 2, 3)
	visitCallExpr(eCall[R]) R
	// myObject.myFunction(1, 2, 3)
	visitGetExpr(eGet[R]) R
	// (1, 2, 3)
	visitGroupingExpr(eGrouping[R]) R
	// 123, "hello", true, false, nil
	visitLiteralExpr(eLiteral[R]) R
	visitLogicalExpr(eLogical[R]) R
	visitSetExpr(eSet[R]) R
	// super.method(1, 2, 3)
	visitSuperExpr(eSuper[R]) R
	// this
	visitThisExpr(eThis[R]) R
	// -1, !true
	visitUnaryExpr(eUnary[R]) R
	// myVariable (accessing variable)
	visitVariableExpr(eVariable[R]) R
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
	return v.visitAssignExpr(e)
}

func (e eBinary[R]) accept(v exprVisitor[R]) R {
	return v.visitBinaryExpr(e)
}

func (e eCall[R]) accept(v exprVisitor[R]) R {
	return v.visitCallExpr(e)
}

func (e eGet[R]) accept(v exprVisitor[R]) R {
	return v.visitGetExpr(e)
}

func (e eGrouping[R]) accept(v exprVisitor[R]) R {
	return v.visitGroupingExpr(e)
}

func (e eLiteral[R]) accept(v exprVisitor[R]) R {
	return v.visitLiteralExpr(e)
}

func (e eLogical[R]) accept(v exprVisitor[R]) R {
	return v.visitLogicalExpr(e)
}

func (e eSet[R]) accept(v exprVisitor[R]) R {
	return v.visitSetExpr(e)
}

func (e eSuper[R]) accept(v exprVisitor[R]) R {
	return v.visitSuperExpr(e)
}

func (e eThis[R]) accept(v exprVisitor[R]) R {
	return v.visitThisExpr(e)
}

func (e eUnary[R]) accept(v exprVisitor[R]) R {
	return v.visitUnaryExpr(e)
}

func (e eVariable[R]) accept(v exprVisitor[R]) R {
	return v.visitVariableExpr(e)
}
