package lox

/*
Defines the statement nodes of the AST. A program is a list of statements.
*/

type stmt[R any] interface {
	accept(stmtVisitor[R]) R
}

type stmtVisitor[R any] interface {
	visitExprStmt(sExpr[R]) R
	visitPrintStmt(sPrint[R]) R
}

type sExpr[R any] struct {
	expression expr[R]
}

type sPrint[R any] struct {
	expression expr[R]
}

func (e sExpr[R]) accept(v stmtVisitor[R]) R {
	return v.visitExprStmt(e)
}

func (e sPrint[R]) accept(v stmtVisitor[R]) R {
	return v.visitPrintStmt(e)
}
