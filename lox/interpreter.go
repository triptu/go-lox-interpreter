package lox

import "fmt"

/*
Interpreter also implements the visitor interface for the AST nodes.
*/

type interpreter struct {
	env *environment
}

var _ exprVisitor = (*interpreter)(nil)
var _ stmtVisitor = (*interpreter)(nil)

func newInterpreter() *interpreter {
	return &interpreter{
		env: newEnvironment(),
	}
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

func (i interpreter) visitBlockStmt(s sBlock) {
	i.executeBlock(s.statements, i.env)
}

func (i interpreter) visitIfStmt(s sIf) {
	val := i.evaluate(s.condition)
	if isTruthy(val) {
		s.thenBranch.accept(i)
	} else if s.elseBranch != nil {
		s.elseBranch.accept(i)
	}
}

func (i interpreter) visitWhileStmt(s sWhile) {
	for {
		val := i.evaluate(s.condition)
		if !isTruthy(val) {
			break
		}
		s.body.accept(i)
	}
}

func (i interpreter) executeBlock(statements []stmt, outerEnv *environment) {
	env := newChildEnvironment(outerEnv)
	i.env = env
	for _, st := range statements {
		st.accept(i)
	}
	i.env = outerEnv
}

func (i interpreter) evaluate(expr expr) any {
	return expr.accept(i)
}

/*
var a = 123;
*/
func (i interpreter) visitVarStmt(s sVar) {
	var val any
	if s.initializer != nil {
		val = i.evaluate(s.initializer)
	}
	i.env.define(s.name.lexeme, val)
}

/*
a = 123;
*/
func (i interpreter) visitAssignExpr(e eAssign) any {
	val := i.evaluate(e.value)
	err := i.env.set(e.name.lexeme, val)
	if err != nil {
		logRuntimeError(e.name.line, "undefined variable '"+e.name.lexeme+"'.")
	}
	return val
}

func (i interpreter) visitBinaryExpr(e eBinary) any {
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

func (i interpreter) visitCallExpr(e eCall) any {
	panic("implement me")
}

func (i interpreter) visitGetExpr(e eGet) any {
	panic("implement me")
}

func (i interpreter) visitGroupingExpr(e eGrouping) any {
	return i.evaluate(e.expression)
}

func (i interpreter) visitLiteralExpr(e eLiteral) any {
	return e.value
}

func (i interpreter) visitLogicalExpr(e eLogical) any {
	left := i.evaluate(e.left)
	if (e.operator.tokenType == tOr && isTruthy(left)) ||
		(e.operator.tokenType == tAnd && !isTruthy(left)) {
		return left
	}
	return i.evaluate(e.right)
}

func (i interpreter) visitSetExpr(e eSet) any {
	panic("implement me")
}

func (i interpreter) visitSuperExpr(e eSuper) any {
	panic("implement me")
}

func (i interpreter) visitThisExpr(e eThis) any {
	panic("implement me")
}

func (i interpreter) visitUnaryExpr(e eUnary) any {
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

func (i interpreter) visitVariableExpr(e eVariable) any {
	val, err := i.env.get(e.name.lexeme)
	if err != nil {
		logRuntimeError(e.name.line, "undefined variable '"+e.name.lexeme+"'.")
	}
	return val
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
