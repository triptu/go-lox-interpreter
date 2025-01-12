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
}

type sExpr struct {
	expression expr[any]
}

type sPrint struct {
	expression expr[any]
}

type sVar struct {
	name        token
	initializer expr[any]
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
