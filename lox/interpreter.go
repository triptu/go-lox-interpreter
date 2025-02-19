package lox

import (
	"context"
	"errors"
	"fmt"
	"math"
)

/*
Interpreter also implements the visitor interface for the AST nodes.
*/

type interpreter struct {
	ctx     context.Context
	globals *environment  // permanent reference to the global environment
	env     *environment  // reference to the environment of the current scope/block
	locals  map[token]int // store the scope depth for each variable token usage
}

var _ exprVisitor = (*interpreter)(nil)
var _ stmtVisitor = (*interpreter)(nil)

func newInterpreter() *interpreter {
	globals := newEnvironment()
	defineNativeFunctions(globals)
	return &interpreter{
		globals: globals,
		locals:  make(map[token]int),
		env:     globals,
	}
}

func (i interpreter) interpret(statements []stmt, ctx context.Context) error {
	i.ctx = ctx

	done := ctx.Done()
	for _, st := range statements {
		select {
		case <-done:
			return ctx.Err()
		default:
			if err := i.execute(st); err != nil {
				return err
			}
		}
	}
	return nil
}

// store the depth of the scope where the variable was found
func (i interpreter) storeResolvedDepth(name token, depth int) {
	i.locals[name] = depth
}

func (i interpreter) visitExprStmt(s sExpr) error {
	_, err := i.evaluate(s.expression)
	return err
}

func (i interpreter) visitPrintStmt(s sPrint) error {
	val := getJustVal(i.evaluate(s.expression))
	logger.Print(getLiteralStr(val))
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

/*
function declaration - fun abc() {...}
we just store the function in env here, and use it when the function is
actually called. The current env is also bound to the function as closure.
*/
func (i interpreter) visitFunctionStmt(s sFunction) error {
	// note that we also attach the env active at the time of function declaration
	i.env.define(s.name.lexeme, loxFunction{declaration: s, closure: i.env})
	return nil
}

/*
class declaration - class abc() {}
we do it in two stages - defining and then setting so the class can be referenced
in its own methods
*/
func (i interpreter) visitClassStmt(s sClass) error {
	var superclass *loxClass
	if s.superclass != nil {
		superclassVal := getJustVal(i.evaluate(s.superclass))
		if superclassVal, ok := superclassVal.(loxClass); !ok {
			logRuntimeError(s.superclass.name, "Superclass must be a class.")
		} else {
			superclass = &superclassVal
		}
	}

	className := s.name.lexeme
	i.env.define(className, nil)

	if superclass != nil {
		i.env = newChildEnvironment(i.env)
		i.env.define("super", superclass)
	}

	methods := make(map[string]loxFunction)
	for _, method := range s.methods {
		methods[method.name.lexeme] = loxFunction{declaration: method, closure: i.env, isInitializer: method.name.lexeme == "init"}
	}
	klass := loxClass{name: className, methods: methods, superclass: superclass}

	if superclass != nil {
		i.env = i.env.outer
	}
	i.env.set(className, klass)
	return nil
}

func (i interpreter) visitWhileStmt(s sWhile) error {
	done := i.ctx.Done()
	for {
		select {
		case <-done:
			return i.ctx.Err()
		default:
			val := getJustVal(i.evaluate(s.condition))
			if !isTruthy(val) {
				return nil
			}
			if err := i.execute(s.body); err != nil {
				return err
			}
		}
	}
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
	dist, exists := i.locals[e.name]
	var err error
	if exists {
		err = i.env.setAt(dist, e.name.lexeme, val)
	} else {
		err = i.globals.set(e.name.lexeme, val)
	}
	if err != nil {
		logRuntimeError(e.name, "Undefined variable '"+e.name.lexeme+"'.")
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
		} else if isList(left) && isList(right) {
			return left.(*loxList).concat([]any{right.(*loxList)}), nil
		} else {
			logRuntimeError(e.operator, "Operands must be two numbers or two strings.")
		}
	case tMinus:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) - right.(float64), nil
	case tStar:
		validateNumberOperand2(left, right, e.operator)
		return left.(float64) * right.(float64), nil
	case tMod:
		validateNumberOperand2(left, right, e.operator)
		return math.Mod(left.(float64), right.(float64)), nil
	case tXor:
		validateNumberOperand2(left, right, e.operator)
		return float64(int(left.(float64)) ^ int(right.(float64))), nil
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
		return checkEqua(left, right), nil
	case tBangEqual:
		return !checkEqua(left, right), nil
	}
	return nil, nil // unreachable
}

func checkEqua(left any, right any) bool {
	switch left := left.(type) {
	case loxClass:
		if right, ok := right.(loxClass); ok {
			return left.name == right.name
		}
		return false
	case loxClassInstance:
		if right, ok := right.(loxClassInstance); ok {
			return &left == &right
		}
		return false
	case loxFunction:
		if right, ok := right.(loxFunction); ok {
			return left.declaration.name.lexeme == right.declaration.name.lexeme && left.closure == right.closure
		}
		return false
	default:
		return left == right
	}
}

func (i interpreter) visitCallExpr(e eCall) (any, error) {
	callee := getJustVal(i.evaluate(e.callee))
	var args []any
	for _, arg := range e.arguments {
		args = append(args, getJustVal(i.evaluate(arg)))
	}
	callee2, ok := callee.(callable)
	if !ok {
		logRuntimeError(e.paren, "Can only call functions and classes.")
	}
	if len(args) != callee2.arity() {
		logRuntimeError(e.paren,
			fmt.Sprintf("Expected %d arguments but got %d.", callee2.arity(), len(args)))
	}
	// fmt.Printf("calling %v with %v\n", callee2, args)
	return callee2.call(i, args)
}

// creating a list - [1,2,3]
func (i interpreter) visitListExpr(e eList) (any, error) {
	list := getLoxList(nil)
	for _, arg := range e.elements {
		val, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		list.elements = append(list.elements, val)
	}
	return list, nil
}

// accessing array index
func (i interpreter) visitGetIndexExpr(e eGetIndex) (any, error) {
	obj, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	indexFloat, err := i.evaluate(e.key)
	if err != nil {
		return nil, err
	}
	index := int(indexFloat.(float64))

	switch obj2 := obj.(type) {
	case *loxList:
		return obj2.getAtIndex(index), nil
	case string:
		if index < 0 {
			index = len(obj2) + index
		}
		if index >= len(obj2) {
			logRuntimeError(e.bracket, "Index out of bounds")
			return nil, errors.New("unreachable")
		}
		return string(obj2[index]), nil
	default:
		logRuntimeError(e.bracket, "Only lists and strings can be accessed by index.")
		return nil, errors.New("unreachable")
	}
}

// setting array index
func (i interpreter) visitSetIndexExpr(e eSetIndex) (any, error) {
	obj, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	indexFloat, err := i.evaluate(e.key)
	if err != nil {
		return nil, err
	}
	index := int(indexFloat.(float64))

	value, err := i.evaluate(e.value)
	if err != nil {
		return nil, err
	}

	switch obj2 := obj.(type) {
	case *loxList:
		return obj2.setAtIndex(index, value), nil
	default:
		logRuntimeError(e.bracket, "Only lists can be mutated by index.")
		return nil, errors.New("unreachable")
	}
}

func (i interpreter) visitReturnStmt(s sReturn) error {
	var value any
	if s.value != nil {
		value = getJustVal(i.evaluate(s.value))
	}
	return returnAsError{value}
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

/*
class field access -
paper.write("hello").withStyle("bold").withColor("red")
*/
func (i interpreter) visitGetExpr(e eGet) (any, error) {
	obj, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}

	switch obj2 := obj.(type) {
	case loxClassInstance:
		return obj2.get(e.name), nil
	case dataType:
		return obj2.getMethod(e.name), nil
	default:
		logRuntimeError(e.name, "Only instances have properties.")
		return nil, errors.New("unreachable")
	}
}

func (i interpreter) visitSetExpr(e eSet) (any, error) {
	obj, err := i.evaluate(e.object)
	if err != nil {
		return nil, err
	}
	switch obj2 := obj.(type) {
	case loxClassInstance:
		value := getJustVal(i.evaluate(e.value))
		return obj2.set(e.name, value), nil
	default:
		logRuntimeError(e.name, "Only instances have fields.")
		return nil, errors.New("unreachable")
	}
}

func (i interpreter) visitSuperExpr(e eSuper) (any, error) {
	distance, ok := i.locals[e.keyword]
	if !ok {
		logRuntimeError(e.keyword, "Couldn't find 'super' in current scope.")
		return nil, errors.New("unreachable")
	}
	superclass, err := i.env.getAt(distance, "super")
	if err != nil {
		logRuntimeError(e.keyword, "No parent class to access.")
		return nil, errors.New("unreachable")
	}
	superclass2 := superclass.(*loxClass)
	object, err := i.env.getAt(distance-1, "this")
	if err != nil {
		logRuntimeError(e.keyword, "No 'this' at super class child.")
		return nil, errors.New("unreachable")
	}
	object2 := object.(loxClassInstance)
	method, ok := superclass2.findMethod(e.method.lexeme)
	if !ok {
		logRuntimeError(e.method, "Undefined property '"+e.method.lexeme+"'.")
		return nil, errors.New("unreachable")
	}
	return method.bind(object2), nil
}

func (i interpreter) visitThisExpr(e eThis) (any, error) {
	return i.lookUpVariable(e.keyword)
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
	val, err := i.lookUpVariable(e.name)
	if err != nil {
		logRuntimeError(e.name, "Undefined variable '"+e.name.lexeme+"'.")
	}
	return val, err
}

func (i interpreter) lookUpVariable(name token) (any, error) {
	dist, exists := i.locals[name]
	varName := name.lexeme
	if exists {
		return i.env.getAt(dist, varName)
	} else {
		return i.globals.get(varName)
	}
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

func isList(value any) bool {
	_, ok := value.(*loxList)
	return ok
}

func validateNumberOperand(num any, operator token) {
	if !isNumber(num) {
		logRuntimeError(operator, "Operand must be a number.")
	}
}

func validateNumberOperand2(num1, num2 any, operator token) {
	if !isNumber(num1) || !isNumber(num2) {
		logRuntimeError(operator, "Operands must be numbers.")
	}
}

func validateNonZeroDenom(denom float64, operator token) {
	if denom == 0 {
		logRuntimeError(operator, "Division by zero")
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
