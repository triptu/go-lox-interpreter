package lox

/*
Defines the statement nodes of the AST. A program is a list of statements.
*/

type stmt interface {
	accept(stmtVisitor) error
}

type stmtVisitor interface {
	// a = 3;
	visitExprStmt(sExpr) error
	// print "hello";
	visitPrintStmt(sPrint) error
	// var a = 3; (initializer is optional)
	visitVarStmt(sVar) error
	// { var a = 3; }
	visitBlockStmt(sBlock) error
	// if (a == 3) { print "hello"; } else { print "world"; }
	visitIfStmt(sIf) error
	// while (a == 3) { print "hello"; a = a + 1; }
	visitWhileStmt(sWhile) error
	// fun foo() { print "hello"; }
	visitFunctionStmt(sFunction) error
	// return 7;
	visitReturnStmt(sReturn) error
	// class Foo { fun bar() { print "hello"; } }
	visitClassStmt(sClass) error
}

type sExpr struct {
	expression expr
}

type sPrint struct {
	expression expr
}

type sVar struct {
	name        token
	initializer expr
}

type sBlock struct {
	statements []stmt
}

type sIf struct {
	condition  expr
	thenBranch stmt
	elseBranch stmt
}

type sWhile struct {
	condition expr
	body      stmt
}

type sClass struct {
	name       token
	superclass *eVariable
	methods    []sFunction
}

type sFunction struct {
	name       token
	parameters []token
	body       []stmt
}

type sReturn struct {
	keyword token
	value   expr
}

func (e sExpr) accept(v stmtVisitor) error {
	return v.visitExprStmt(e)
}

func (e sPrint) accept(v stmtVisitor) error {
	return v.visitPrintStmt(e)
}

func (e sVar) accept(v stmtVisitor) error {
	return v.visitVarStmt(e)
}

func (e sBlock) accept(v stmtVisitor) error {
	return v.visitBlockStmt(e)
}

func (e sIf) accept(v stmtVisitor) error {
	return v.visitIfStmt(e)
}

func (e sWhile) accept(v stmtVisitor) error {
	return v.visitWhileStmt(e)
}

func (e sFunction) accept(v stmtVisitor) error {
	return v.visitFunctionStmt(e)
}

func (e sReturn) accept(v stmtVisitor) error {
	return v.visitReturnStmt(e)
}

func (e sClass) accept(v stmtVisitor) error {
	return v.visitClassStmt(e)
}
