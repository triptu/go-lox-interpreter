package lox

/*
This file contains a lot of struct which are data types to represent different
expressions in our language. These are basically the nodes of the
AST(abstract syntax tree) of our language.
*/

type expr interface {
	accept(exprVisitor) any
}

/*
action classes implement this interface defining how they would process each type of expression
*/
type exprVisitor interface {
	// a = 3
	visitAssignExpr(eAssign) any
	// 1 + 2, 5 * 6, 2 < 3, etc
	visitBinaryExpr(eBinary) any
	// myFunction(1, 2, 3)
	visitCallExpr(eCall) any
	// myObject.myFunction(1, 2, 3)
	visitGetExpr(eGet) any
	// (1, 2, 3)
	visitGroupingExpr(eGrouping) any
	// 123, "hello", true, false, nil
	visitLiteralExpr(eLiteral) any
	// true or false, "abcd" and "efgh"
	// these are not coupled with binary as whether the right expression is evaluated
	// depends on the left expression's evaluation
	visitLogicalExpr(eLogical) any
	visitSetExpr(eSet) any
	// super.method(1, 2, 3)
	visitSuperExpr(eSuper) any
	// this
	visitThisExpr(eThis) any
	// -1, !true
	visitUnaryExpr(eUnary) any
	// myVariable (accessing variable)
	visitVariableExpr(eVariable) any
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

// define accept methods for each type of expression

func (e eAssign) accept(v exprVisitor) any {
	return v.visitAssignExpr(e)
}

func (e eBinary) accept(v exprVisitor) any {
	return v.visitBinaryExpr(e)
}

func (e eCall) accept(v exprVisitor) any {
	return v.visitCallExpr(e)
}

func (e eGet) accept(v exprVisitor) any {
	return v.visitGetExpr(e)
}

func (e eGrouping) accept(v exprVisitor) any {
	return v.visitGroupingExpr(e)
}

func (e eLiteral) accept(v exprVisitor) any {
	return v.visitLiteralExpr(e)
}

func (e eLogical) accept(v exprVisitor) any {
	return v.visitLogicalExpr(e)
}

func (e eSet) accept(v exprVisitor) any {
	return v.visitSetExpr(e)
}

func (e eSuper) accept(v exprVisitor) any {
	return v.visitSuperExpr(e)
}

func (e eThis) accept(v exprVisitor) any {
	return v.visitThisExpr(e)
}

func (e eUnary) accept(v exprVisitor) any {
	return v.visitUnaryExpr(e)
}

func (e eVariable) accept(v exprVisitor) any {
	return v.visitVariableExpr(e)
}
