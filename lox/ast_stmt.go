package lox

/*
Defines the statement nodes of the AST. A program is a list of statements.
*/

type stmt interface {
	accept(stmtVisitor)
}

type stmtVisitor interface {
	// a = 3;
	visitExprStmt(sExpr)
	// print "hello";
	visitPrintStmt(sPrint)
	// var a = 3; (initializer is optional)
	visitVarStmt(sVar)
	// { var a = 3; }
	visitBlockStmt(sBlock)
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

func (e sExpr) accept(v stmtVisitor) {
	v.visitExprStmt(e)
}

func (e sPrint) accept(v stmtVisitor) {
	v.visitPrintStmt(e)
}

func (e sVar) accept(v stmtVisitor) {
	v.visitVarStmt(e)
}

func (e sBlock) accept(v stmtVisitor) {
	v.visitBlockStmt(e)
}
