package lox

import "fmt"

/*
Interpreter also implements the visitor interface for the AST nodes.
*/

type interpreter struct {
}

func (i interpreter) interpret(statements []stmt) {
	for _, st := range statements {
		st.accept(i)
	}
}

func (i interpreter) visitExprStmt(s sExpr) {
	i.evaluate(s.expression)
}

func (i interpreter) visitPrintStmt(s sPrint) {
	val := i.evaluate(s.expression)
	fmt.Println(getLiteralStr(val))
}

func (i interpreter) evaluate(expr expr[any]) any {
	return expr.accept(i)
}

func (i interpreter) visitAssign(e eAssign[any]) any {
	panic("implement me")
}

func (i interpreter) visitBinary(e eBinary[any]) any {
	left := i.evaluate(e.left)
	right := i.evaluate(e.right)
	switch e.operator.tokenType {
	case tPlus:
		if isString(left) && isString(right) {
			return left.(string) + right.(string)
		} else if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64)
		} else {
			logRuntimeError(e.operator.line, "for plus, operands must be two numbers or two strings.")
		}
	case tMinus:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) - right.(float64)
	case tStar:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) * right.(float64)
	case tSlash:
		validateNumberOperand2(left, right, e.operator)
		validateNonZeroDenom(right.(float64), e.operator)
		return left.(float64) / right.(float64)
	case tGreater:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) > right.(float64)
	case tGreaterEqual:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) >= right.(float64)
	case tLess:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) < right.(float64)
	case tLessEqual:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) <= right.(float64)
	case tEqualEqual:
		return left == right
	case tBangEqual:
		return left != right
	}
	return nil // unreachable
}

func (i interpreter) visitCall(e eCall[any]) any {
	panic("implement me")
}

func (i interpreter) visitGet(e eGet[any]) any {
	panic("implement me")
}

func (i interpreter) visitGrouping(e eGrouping[any]) any {
	return i.evaluate(e.expression)
}

func (i interpreter) visitLiteral(e eLiteral[any]) any {
	return e.value
}

func (i interpreter) visitLogical(e eLogical[any]) any {
	panic("implement me")
}

func (i interpreter) visitSet(e eSet[any]) any {
	panic("implement me")
}

func (i interpreter) visitSuper(e eSuper[any]) any {
	panic("implement me")
}

func (i interpreter) visitThis(e eThis[any]) any {
	panic("implement me")
}

func (i interpreter) visitUnary(e eUnary[any]) any {
	right := i.evaluate(e.right)
	switch e.operator.tokenType {
	case tMinus:
		validateNumberOperand(right, e.operator)
		return -right.(float64)
	case tBang:
		return !isTruthy(right)
	default:
		return nil // unreachable
	}
}

func (i interpreter) visitVariable(e eVariable[any]) any {
	panic("implement me")
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	switch value := value.(type) {
	case bool:
		return value
	default:
		return true
	}
}

func isString(value any) bool {
	_, ok := value.(string)
	return ok
}

func isNumber(value any) bool {
	_, ok := value.(float64)
	return ok
}

func validateNumberOperand(num any, operator token) {
	if !isNumber(num) {
		logRuntimeError(operator.line, "Operand must be a number for operator: "+operator.lexeme)
	}
}

func validateNumberOperand2(num1, num2 any, operator token) {
	if !isNumber(num1) || !isNumber(num2) {
		logRuntimeError(operator.line, "Operands must be numbers for operator: "+operator.lexeme)
	}
}

func validateNonZeroDenom(denom float64, operator token) {
	if denom == 0 {
		logRuntimeError(operator.line, "Division by zero")
	}
}
