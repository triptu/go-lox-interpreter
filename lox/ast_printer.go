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

func (p astPrinter) print(e expr[any]) {
	fmt.Println(e.accept(p))
}

func (p astPrinter) visitAssign(e eAssign[any]) any {
	return p.parenthesize("= "+e.name.lexeme, e.value)
}

/*
for e.g. "1 + 2" => "(+ 1 2)"
*/
func (p astPrinter) visitBinary(e eBinary[any]) any {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitCall(e eCall[any]) any {
	return p.parenthesize("call", append([]expr[any]{e.callee}, e.arguments...)...)
}

func (p astPrinter) visitGet(e eGet[any]) any {
	return p.parenthesize(".", e.object, eLiteral[any]{value: e.name.lexeme})
}

func (p astPrinter) visitGrouping(e eGrouping[any]) any {
	return p.parenthesize("group", e.expression)
}

func (p astPrinter) visitLiteral(e eLiteral[any]) any {
	return getTokenLiteralStr(e.value)
}

func (p astPrinter) visitLogical(e eLogical[any]) any {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitSet(e eSet[any]) any {
	return p.parenthesize("=", e.object, eLiteral[any]{value: e.name.lexeme}, e.value)
}

func (p astPrinter) visitSuper(e eSuper[any]) any {
	return p.parenthesize("super " + e.method.lexeme + " " + e.keyword.lexeme)
}

func (p astPrinter) visitThis(e eThis[any]) any {
	return p.parenthesize(e.keyword.lexeme)
}

func (p astPrinter) visitUnary(e eUnary[any]) any {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p astPrinter) visitVariable(e eVariable[any]) any {
	return e.name.lexeme
}

func (p astPrinter) parenthesize(name string, exprs ...expr[any]) any {
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
