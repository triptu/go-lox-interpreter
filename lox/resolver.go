package lox

// enum to track if we're inside a function
type functionType int

const (
	fNone = iota
	fFunction
	fMethod
	fInitializer // we're in a constructor, user can't return from there
)

// enum to track if we're inside a class
type classType int

const (
	cNone = iota
	cClass
	cSubClass
)

type resolver struct {
	// Binding is split into two phases: declaration and definition.
	// for each variable, in the scope it is part of, when it's declared we add it as
	// "not ready yet" in the declaration phase. And mark it true when the variable after
	// the initializer has run. This is to avoid the initializer to reference the variable
	// which is being declared.
	// there is no global scope, as if variable isn't part of any local scope, it's
	// obviously part of the global scope.
	scopes       []map[string]bool // stack of nested lexical scopes
	interpreter  *interpreter
	currFunction functionType
	currClass    classType
}

var _ exprVisitor = (*resolver)(nil)
var _ stmtVisitor = (*resolver)(nil)

func newResolver(interpreter *interpreter) *resolver {
	return &resolver{
		scopes:       []map[string]bool{},
		interpreter:  interpreter,
		currFunction: fNone,
		currClass:    cNone,
	}
}

func (r *resolver) resolve(stmts []stmt) {
	r.resolveStmts(stmts)
}

func (r *resolver) resolveStmts(stmts []stmt) {
	for _, s := range stmts {
		if err := r.resolveStmt(s); err != nil {
			if pErr, ok := err.(*parseError); ok {
				logParseError(pErr.token, pErr.msg)
			} else {
				logParseError(token{}, err.Error())
			}
		}
	}
}

func (r *resolver) visitBlockStmt(stmt sBlock) error {
	r.beginScope()
	r.resolveStmts(stmt.statements)
	r.endScope()
	return nil
}

func (r *resolver) resolveStmt(stmt stmt) error {
	return stmt.accept(r)
}

func (r *resolver) resolveExpr(expr expr) (any, error) {
	return expr.accept(r)
}

/*
for e.g. var a = 3; a = a + 1;
*/
func (r *resolver) visitVarStmt(stmt sVar) error {
	if err := r.declare(stmt.name); err != nil {
		return err
	}
	if stmt.initializer != nil {
		// the initializer can't reference the variable which is being declared
		// for e.g. var a = a + 1; is invalid
		// putting variable in our scope at time of declaration helps us detect this
		// error. Otherwise, if there is a variable with the same name in parent scope,
		// interpreter will use that variable's value. This will lead to inconsistent
		// behavior.
		if _, err := r.resolveExpr(stmt.initializer); err != nil {
			return err
		}
	}
	r.define(stmt.name.lexeme)
	return nil
}

/*
called when the variables are used. Note that if it's used in its own initializer, it's
an error. We resolve to the scope the variable is declared in, and save it in the interpreter.
for e.g. a + b;
*/
func (r *resolver) visitVariableExpr(expr eVariable) (any, error) {
	varName := expr.name.lexeme
	if len(r.scopes) != 0 {
		if isReady, exists := r.peekScope()[varName]; exists && !isReady {
			return nil, parseErrorAt(expr.name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr.name)
	return nil, nil
}

func (r *resolver) visitAssignExpr(expr eAssign) (any, error) {
	if _, err := r.resolveExpr(expr.value); err != nil {
		return nil, err
	}
	r.resolveLocal(expr.name)
	return nil, nil
}

func (r *resolver) visitBinaryExpr(expr eBinary) (any, error) {
	if _, err := r.resolveExpr(expr.left); err != nil {
		return nil, err
	}
	if _, err := r.resolveExpr(expr.right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitCallExpr(expr eCall) (any, error) {
	if _, err := r.resolveExpr(expr.callee); err != nil {
		return nil, err
	}
	for _, arg := range expr.arguments {
		if _, err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *resolver) visitListExpr(expr eList) (any, error) {
	for _, arg := range expr.elements {
		if _, err := r.resolveExpr(arg); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *resolver) visitGroupingExpr(expr eGrouping) (any, error) {
	return r.resolveExpr(expr.expression)
}

func (r *resolver) visitLiteralExpr(expr eLiteral) (any, error) {
	return nil, nil
}

func (r *resolver) visitLogicalExpr(expr eLogical) (any, error) {
	if _, err := r.resolveExpr(expr.left); err != nil {
		return nil, err
	}
	if _, err := r.resolveExpr(expr.right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitUnaryExpr(expr eUnary) (any, error) {
	if _, err := r.resolveExpr(expr.right); err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *resolver) visitGetExpr(expr eGet) (any, error) {
	return r.resolveExpr(expr.object)
}

func (r *resolver) visitSetExpr(expr eSet) (any, error) {
	if _, err := r.resolveExpr(expr.object); err != nil {
		return nil, err
	}
	_, err := r.resolveExpr(expr.value)
	return nil, err
}

func (r *resolver) visitSuperExpr(expr eSuper) (any, error) {
	if r.currClass == cNone {
		return nil, parseErrorAt(expr.keyword, "Can't use 'super' outside of a class.")
	} else if r.currClass != cSubClass {
		return nil, parseErrorAt(expr.keyword, "Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr.keyword)
	return nil, nil
}

func (r *resolver) visitThisExpr(expr eThis) (any, error) {
	if r.currClass == cNone {
		return nil, parseErrorAt(expr.keyword, "Can't use 'this' outside of a class.")
	}
	r.resolveLocal(expr.keyword)
	return nil, nil
}

func (r *resolver) visitFunctionStmt(stmt sFunction) error {
	if err := r.declare(stmt.name); err != nil {
		return err
	}
	// we define right away, as it's legal for the function to reference itself for recursion
	r.define(stmt.name.lexeme)

	err := r.resolveFunction(stmt, fFunction)
	return err
}

func (r *resolver) visitExprStmt(stmt sExpr) error {
	_, err := r.resolveExpr(stmt.expression)
	return err
}

func (r *resolver) visitPrintStmt(stmt sPrint) error {
	_, err := r.resolveExpr(stmt.expression)
	return err
}

func (r *resolver) visitReturnStmt(stmt sReturn) error {
	if r.currFunction == fNone {
		return parseErrorAt(stmt.keyword, "Can't return from top-level code.")
	}

	if stmt.value != nil {
		if r.currFunction == fInitializer {
			return parseErrorAt(stmt.keyword, "Can't return a value from an initializer.")
		}
		_, err := r.resolveExpr(stmt.value)
		return err
	}
	return nil
}

func (r *resolver) visitIfStmt(stmt sIf) error {
	if _, err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}
	if err := r.resolveStmt(stmt.thenBranch); err != nil {
		return err
	}
	if stmt.elseBranch != nil {
		return r.resolveStmt(stmt.elseBranch)
	}
	return nil
}

func (r *resolver) visitWhileStmt(stmt sWhile) error {
	if _, err := r.resolveExpr(stmt.condition); err != nil {
		return err
	}
	return r.resolveStmt(stmt.body)
}

func (r *resolver) visitClassStmt(stmt sClass) error {
	enclosingClass := r.currClass
	r.currClass = cClass
	defer func() {
		r.currClass = enclosingClass
		r.endScope()
		if stmt.superclass != nil {
			r.endScope()
		}
	}()
	if err := r.declare(stmt.name); err != nil {
		return err
	}
	r.define(stmt.name.lexeme)

	if stmt.superclass != nil {
		r.currClass = cSubClass
		if stmt.name.lexeme == stmt.superclass.name.lexeme {
			return parseErrorAt(stmt.superclass.name, "A class can't inherit from itself.")
		}
		if _, err := r.resolveExpr(*stmt.superclass); err != nil {
			return err
		}
		r.beginScope()
		r.peekScope()["super"] = true
	}

	r.beginScope()
	r.peekScope()["this"] = true

	for _, method := range stmt.methods {
		var declarationType functionType = fMethod
		if method.name.lexeme == "init" {
			declarationType = fInitializer
		}
		if err := r.resolveFunction(method, declarationType); err != nil {
			return err
		}
		r.define(method.name.lexeme)
	}
	return nil
}

func (r *resolver) resolveFunction(function sFunction, funcType functionType) error {
	enclosingFunction := r.currFunction
	r.currFunction = funcType
	defer func() {
		r.currFunction = enclosingFunction
		r.endScope()
	}()

	r.beginScope()
	for _, param := range function.parameters {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param.lexeme)
	}
	r.resolveStmts(function.body)
	return nil
}

func (r *resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *resolver) endScope() {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

/*
the name is declared. We also check if the name is already declared in the current scope,
which is an error. That is the same variable for e.g. can't be declared twice like below:
var a = 3;
var a = 4;
It becomes a parsing error.
*/
func (r *resolver) declare(nameToken token) error {
	name := nameToken.lexeme
	if len(r.scopes) == 0 {
		return nil
	}
	_, exists := r.peekScope()[name]
	if exists {
		return parseErrorAt(nameToken, "Already a variable with this name in this scope.")
	}
	r.peekScope()[name] = false
	return nil
}

/*
the name is defined, and ready to be used
*/
func (r *resolver) define(name string) {
	if len(r.scopes) == 0 {
		return
	}
	r.peekScope()[name] = true
}

func (r *resolver) peekScope() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

// the exprName token here is one of tVariable, tAssign, tThis or tSuper
func (r *resolver) resolveLocal(exprName token) {
	varName := exprName.lexeme // variable name
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, exists := r.scopes[i][varName]; exists {
			// number of scopes between the current innermost scope and the scope where the variable was found
			depth := len(r.scopes) - 1 - i
			r.interpreter.storeResolvedDepth(exprName, depth)
			return
		}
	}
}
