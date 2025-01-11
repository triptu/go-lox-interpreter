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

func (p astPrinter) print(e expr[string]) {
	fmt.Println(e.accept(p))
}

func (p astPrinter) visitAssign(e eAssign[string]) string {
	return p.parenthesize("= "+e.name.lexeme, e.value)
}

/*
for e.g. "1 + 2" => "(+ 1 2)"
*/
func (p astPrinter) visitBinary(e eBinary[string]) string {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitCall(e eCall[string]) string {
	return p.parenthesize("call", append([]expr[string]{e.callee}, e.arguments...)...)
}

func (p astPrinter) visitGet(e eGet[string]) string {
	return p.parenthesize(".", e.object, eLiteral[string]{value: e.name.lexeme})
}

func (p astPrinter) visitGrouping(e eGrouping[string]) string {
	return p.parenthesize("group", e.expression)
}

func (p astPrinter) visitLiteral(e eLiteral[string]) string {
	if e.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", e.value)
}

func (p astPrinter) visitLogical(e eLogical[string]) string {
	return p.parenthesize(e.operator.lexeme, e.left, e.right)
}

func (p astPrinter) visitSet(e eSet[string]) string {
	return p.parenthesize("=", e.object, eLiteral[string]{value: e.name.lexeme}, e.value)
}

func (p astPrinter) visitSuper(e eSuper[string]) string {
	return p.parenthesize("super " + e.method.lexeme + " " + e.keyword.lexeme)
}

func (p astPrinter) visitThis(e eThis[string]) string {
	return p.parenthesize(e.keyword.lexeme)
}

func (p astPrinter) visitUnary(e eUnary[string]) string {
	return p.parenthesize(e.operator.lexeme, e.right)
}

func (p astPrinter) visitVariable(e eVariable[string]) string {
	return e.name.lexeme
}

func (p astPrinter) parenthesize(name string, exprs ...expr[string]) string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(name)
	for _, expr := range exprs {
		sb.WriteString(" ")
		sb.WriteString(expr.accept(p))
	}
	sb.WriteString(")")
	return sb.String()
}

func testPrinter() {
	p := astPrinter{}
	expr := eBinary[string]{
		left: eUnary[string]{
			operator: token{tokenType: tMinus, lexeme: "-"},
			right: eLiteral[string]{
				value: 123,
			},
		},
		operator: token{tokenType: tStar, lexeme: "*"},
		right: eGrouping[string]{
			expression: eLiteral[string]{
				value: 45.67,
			},
		},
	}
	p.print(expr)
}
