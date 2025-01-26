package lox

import (
	"fmt"
	"strings"
)

/*
AST Printer implements the exprVisitor, returning a string representation from the visit methods
*/

type astPrinter struct {
}

func (p astPrinter) print(e expr) {
	val, _ := e.accept(p)
	fmt.Println(val)
}

func (p astPrinter) visitAssignExpr(e eAssign) (any, error) {
	return p.parenthesize("= "+e.name.lexeme, e.value)
}

/*
for e.g. "1 + 2" => "(+ 1 2)"
*/
func (p astPrinter) visitBinaryExpr(e eBinary) (any, error) {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitCallExpr(e eCall) (any, error) {
	return p.parenthesize("call", append([]expr{e.callee}, e.arguments...)...)
}

func (p astPrinter) visitGetExpr(e eGet) (any, error) {
	return p.parenthesize(".", e.object, eLiteral{value: e.name.lexeme})
}

func (p astPrinter) visitGroupingExpr(e eGrouping) (any, error) {
	return p.parenthesize("group", e.expression)
}

func (p astPrinter) visitLiteralExpr(e eLiteral) (any, error) {
	return getTokenLiteralStr(e.value), nil
}

func (p astPrinter) visitLogicalExpr(e eLogical) (any, error) {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitSetExpr(e eSet) (any, error) {
	return p.parenthesize("=", e.object, eLiteral{value: e.name.lexeme}, e.value)
}

func (p astPrinter) visitSuperExpr(e eSuper) (any, error) {
	return p.parenthesize("super " + e.method.lexeme + " " + e.keyword.lexeme)
}

func (p astPrinter) visitThisExpr(e eThis) (any, error) {
	return p.parenthesize(e.keyword.lexeme)
}

func (p astPrinter) visitUnaryExpr(e eUnary) (any, error) {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p astPrinter) visitVariableExpr(e eVariable) (any, error) {
	return e.name.lexeme, nil
}

func (p astPrinter) visitListExpr(e eList) (any, error) {
	var sb strings.Builder
	sb.WriteString("[")
	for _, expr := range e.elements {
		sb.WriteString(" ")
		evald, _ := expr.accept(p)
		sb.WriteString(evald.(string))
	}
	sb.WriteString("]")
	return sb.String(), nil
}

func (p astPrinter) parenthesize(name string, exprs ...expr) (any, error) {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		evald, _ := expr.accept(p)
		sb.WriteString(evald.(string))
	}
	sb.WriteString(")")
	return sb.String(), nil
}
