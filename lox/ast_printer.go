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
	fmt.Println(e.accept(p))
}

func (p astPrinter) visitAssignExpr(e eAssign) any {
	return p.parenthesize("= "+e.name.lexeme, e.value)
}

/*
for e.g. "1 + 2" => "(+ 1 2)"
*/
func (p astPrinter) visitBinaryExpr(e eBinary) any {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitCallExpr(e eCall) any {
	return p.parenthesize("call", append([]expr{e.callee}, e.arguments...)...)
}

func (p astPrinter) visitGetExpr(e eGet) any {
	return p.parenthesize(".", e.object, eLiteral{value: e.name.lexeme})
}

func (p astPrinter) visitGroupingExpr(e eGrouping) any {
	return p.parenthesize("group", e.expression)
}

func (p astPrinter) visitLiteralExpr(e eLiteral) any {
	return getTokenLiteralStr(e.value)
}

func (p astPrinter) visitLogicalExpr(e eLogical) any {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitSetExpr(e eSet) any {
	return p.parenthesize("=", e.object, eLiteral{value: e.name.lexeme}, e.value)
}

func (p astPrinter) visitSuperExpr(e eSuper) any {
	return p.parenthesize("super " + e.method.lexeme + " " + e.keyword.lexeme)
}

func (p astPrinter) visitThisExpr(e eThis) any {
	return p.parenthesize(e.keyword.lexeme)
}

func (p astPrinter) visitUnaryExpr(e eUnary) any {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p astPrinter) visitVariableExpr(e eVariable) any {
	return e.name.lexeme
}

func (p astPrinter) parenthesize(name string, exprs ...expr) any {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.accept(p).(string))
	}
	sb.WriteString(")")
	return sb.String()
}
