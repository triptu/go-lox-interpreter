package lox

/*
Defines the statement nodes of the AST. A program is a list of statements.
*/

type stmt interface {
	accept(stmtVisitor)
}

type stmtVisitor interface {
	visitExprStmt(sExpr)
	visitPrintStmt(sPrint)
}

type sExpr struct {
	expression expr[any]
}

type sPrint struct {
	expression expr[any]
}

func (e sExpr) accept(v stmtVisitor) {
	v.visitExprStmt(e)
}

func (e sPrint) accept(v stmtVisitor) {
	v.visitPrintStmt(e)
}
