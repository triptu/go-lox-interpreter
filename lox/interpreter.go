package lox

import (
	"fmt"
)

/*
Interpreter also implements the visitor interface for the AST nodes.
*/

type interpreter struct {
	globals *environment // permanent reference to the global environment
	env     *environment // reference to the environment of the current scope/block
}

var _ exprVisitor = (*interpreter)(nil)
var _ stmtVisitor = (*interpreter)(nil)

func newInterpreter() *interpreter {
	globals := newEnvironment()
	globals.define("clock", nativeClock{}) // a native language provided function
	globals.define("print", nativePrint{})
	return &interpreter{
		globals: globals,
		env:     globals,
	}
}

func (i interpreter) interpret(statements []stmt) error {
	for _, st := range statements {
		if err := i.execute(st); err != nil {
			return err
		}
	}
	return nil
}

func (i interpreter) visitExprStmt(s sExpr) error {
	_, err := i.evaluate(s.expression)
	return err
}

func (i interpreter) visitPrintStmt(s sPrint) error {
	val := getJustVal(i.evaluate(s.expression))
	fmt.Println(getLiteralStr(val))
	return nil
}

func (i interpreter) visitBlockStmt(s sBlock) error {
	return i.executeBlock(s.statements, i.env)
}

func (i interpreter) visitIfStmt(s sIf) error {
	val := getJustVal(i.evaluate(s.condition))
	if isTruthy(val) {
		return i.execute(s.thenBranch)
	} else if s.elseBranch != nil {
		return i.execute(s.elseBranch)
	}
	return nil
}

func (i interpreter) visitFunctionStmt(s sFunction) error {
	i.env.define(s.name.lexeme, loxFunction{s})
	return nil
}

func (i interpreter) visitWhileStmt(s sWhile) error {
	for {
		val := getJustVal(i.evaluate(s.condition))
		if !isTruthy(val) {
			break
		}
		if err := i.execute(s.body); err != nil {
			return err
		}
	}
	return nil
}

func (i interpreter) executeBlock(statements []stmt, outerEnv *environment) error {
	defer func() { i.env = outerEnv }() // restore outer environment at the end
	env := newChildEnvironment(outerEnv)
	i.env = env
	for _, st := range statements {
		if err := i.execute(st); err != nil {
			return err
		}
	}
	return nil
}

func (i interpreter) execute(stmt stmt) error {
	return stmt.accept(i)
}

func (i interpreter) evaluate(expr expr) (any, error) {
	return expr.accept(i)
}

/*
var a = 123;
*/
func (i interpreter) visitVarStmt(s sVar) error {
	var val any
	if s.initializer != nil {
		val, _ = i.evaluate(s.initializer)
	}
	i.env.define(s.name.lexeme, val)
	return nil
}

/*
a = 123;
*/
func (i interpreter) visitAssignExpr(e eAssign) (any, error) {
	val := getJustVal(i.evaluate(e.value))
	err := i.env.set(e.name.lexeme, val)
	if err != nil {
		logRuntimeError(e.name.line, "undefined variable '"+e.name.lexeme+"'.")
	}
	return val, nil
}

func (i interpreter) visitBinaryExpr(e eBinary) (any, error) {
	left := getJustVal(i.evaluate(e.left))
	right := getJustVal(i.evaluate(e.right))
	switch e.operator.tokenType {
	case tPlus:
		if isString(left) || isString(right) {
			// if either side is string, convert the other side to string as well
			return getLiteralStr(left) + getLiteralStr(right), nil
		} else if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64), nil
		} else {
			logRuntimeError(e.operator.line, "for plus, operands must be two numbers or two strings.")
		}
	case tMinus:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) - right.(float64), nil
	case tStar:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) * right.(float64), nil
	case tSlash:
		validateNumberOperand2(left, right, e.operator)
		validateNonZeroDenom(right.(float64), e.operator)
		return left.(float64) / right.(float64), nil
	case tGreater:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) > right.(float64), nil
	case tGreaterEqual:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) >= right.(float64), nil
	case tLess:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) < right.(float64), nil
	case tLessEqual:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) <= right.(float64), nil
	case tEqualEqual:
		return left == right, nil
	case tBangEqual:
		return left != right, nil
	}
	return nil, nil // unreachable
}

func (i interpreter) visitCallExpr(e eCall) (any, error) {
	callee := getJustVal(i.evaluate(e.callee))
	var args []any
	for _, arg := range e.arguments {
		args = append(args, getJustVal(i.evaluate(arg)))
	}
	callee2, ok := callee.(callable)
	if !ok {
		logRuntimeError(e.paren.line, "Can only call functions and classes.")
	}
	if len(args) != callee2.arity() {
		logRuntimeError(e.paren.line,
			fmt.Sprintf("Expected %d arguments but got %d.", callee2.arity(), len(args)))
	}
	return callee2.call(i, args)
}

func (i interpreter) visitReturnStmt(s sReturn) error {
	var value any
	if s.value != nil {
		value = getJustVal(i.evaluate(s.value))
	}
	return returnAsError{value}
}

func (i interpreter) visitGetExpr(e eGet) (any, error) {
	panic("implement me")
}

func (i interpreter) visitGroupingExpr(e eGrouping) (any, error) {
	return i.evaluate(e.expression)
}

func (i interpreter) visitLiteralExpr(e eLiteral) (any, error) {
	return e.value, nil
}

func (i interpreter) visitLogicalExpr(e eLogical) (any, error) {
	left := getJustVal(i.evaluate(e.left))
	if (e.operator.tokenType == tOr && isTruthy(left)) ||
		(e.operator.tokenType == tAnd && !isTruthy(left)) {
		return left, nil
	}
	return i.evaluate(e.right)
}

func (i interpreter) visitSetExpr(e eSet) (any, error) {
	panic("implement me")
}

func (i interpreter) visitSuperExpr(e eSuper) (any, error) {
	panic("implement me")
}

func (i interpreter) visitThisExpr(e eThis) (any, error) {
	panic("implement me")
}

func (i interpreter) visitUnaryExpr(e eUnary) (any, error) {
	right := getJustVal(i.evaluate(e.right))
	switch e.operator.tokenType {
	case tMinus:
		validateNumberOperand(right, e.operator)
		return -right.(float64), nil
	case tBang:
		return !isTruthy(right), nil
	default:
		return nil, nil // unreachable
	}
}

func (i interpreter) visitVariableExpr(e eVariable) (any, error) {
	val, err := i.env.get(e.name.lexeme)
	if err != nil {
		logRuntimeError(e.name.line, "undefined variable '"+e.name.lexeme+"'.")
	}
	return val, err
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

/*
shortcut for some places, technically a bad thing to do. we'll just call it
at places where we don't expect error to be returned.

Note that in the interpreter actual errors are not really being bubbled, they lead
to immediate panic. This error stuff is being used more like a control flow for return
statement.

Maybe when we later expand to give stack traces or enhanced runtime error reporting,
this function will be removed.
*/
func getJustVal[T any](val T, _ error) T {
	return val
}
