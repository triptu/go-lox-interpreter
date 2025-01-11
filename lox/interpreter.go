package lox

/*
Interpreter also implements the visitor interface for the AST nodes.
*/

type interpreter struct {
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
		}
	case tMinus:
		return left.(float64) - right.(float64)
	case tStar:
		return left.(float64) * right.(float64)
	case tSlash:
		return left.(float64) / right.(float64)
	case tGreater:
		return left.(float64) > right.(float64)
	case tGreaterEqual:
		return left.(float64) >= right.(float64)
	case tLess:
		return left.(float64) < right.(float64)
	case tLessEqual:
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
