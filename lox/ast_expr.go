package lox

/*
This file contains a lot of struct which are data types to represent different
expressions in our language. These are basically the nodes of the
AST(abstract syntax tree) of our language.
*/

type expr interface {
	accept(exprVisitor) (any, error)
}

/*
action classes implement this interface defining how they would process each type of expression
*/
type exprVisitor interface {
	// a = 3
	visitAssignExpr(eAssign) (any, error)
	// 1 + 2, 5 * 6, 2 < 3, etc
	visitBinaryExpr(eBinary) (any, error)
	// myFunction(1, 2, 3)
	visitCallExpr(eCall) (any, error)
	// myObject.myFunction(1, 2, 3) or breakfast.milk.sugarLevel
	visitGetExpr(eGet) (any, error)
	// (1, 2, 3)
	visitGroupingExpr(eGrouping) (any, error)
	// 123, "hello", true, false, nil
	visitLiteralExpr(eLiteral) (any, error)
	// true or false, "abcd" and "efgh"
	// these are not coupled with binary as whether the right expression is evaluated
	// depends on the left expression's evaluation
	visitLogicalExpr(eLogical) (any, error)
	// breakfast.milk.sugarLevel = 4, setting fields of class instance
	visitSetExpr(eSet) (any, error)
	// super.method(1, 2, 3)
	visitSuperExpr(eSuper) (any, error)
	// this (used in class for self-refernce in methods)
	visitThisExpr(eThis) (any, error)
	// -1, !true
	visitUnaryExpr(eUnary) (any, error)
	// myVariable (accessing variable)
	// a + b;
	visitVariableExpr(eVariable) (any, error)
	// [1, 2, 3]
	visitListExpr(eList) (any, error)
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
	paren     token // stored only for error reporting
	arguments []expr
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

// object.name is being accessed
type eGet struct {
	object expr
	name   token
}

// object.name = value
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

type eList struct {
	elements []expr
}

// define accept methods for each type of expression

func (e eAssign) accept(v exprVisitor) (any, error) {
	return v.visitAssignExpr(e)
}

func (e eBinary) accept(v exprVisitor) (any, error) {
	return v.visitBinaryExpr(e)
}

func (e eCall) accept(v exprVisitor) (any, error) {
	return v.visitCallExpr(e)
}

func (e eGet) accept(v exprVisitor) (any, error) {
	return v.visitGetExpr(e)
}

func (e eGrouping) accept(v exprVisitor) (any, error) {
	return v.visitGroupingExpr(e)
}

func (e eLiteral) accept(v exprVisitor) (any, error) {
	return v.visitLiteralExpr(e)
}

func (e eLogical) accept(v exprVisitor) (any, error) {
	return v.visitLogicalExpr(e)
}

func (e eSet) accept(v exprVisitor) (any, error) {
	return v.visitSetExpr(e)
}

func (e eSuper) accept(v exprVisitor) (any, error) {
	return v.visitSuperExpr(e)
}

func (e eThis) accept(v exprVisitor) (any, error) {
	return v.visitThisExpr(e)
}

func (e eUnary) accept(v exprVisitor) (any, error) {
	return v.visitUnaryExpr(e)
}

func (e eVariable) accept(v exprVisitor) (any, error) {
	return v.visitVariableExpr(e)
}

func (e eList) accept(v exprVisitor) (any, error) {
	return v.visitListExpr(e)
}
