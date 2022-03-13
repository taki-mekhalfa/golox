package ast

type Stmt interface {
	Accept(VisitorStmt) interface{}
}

type Print struct {
	Expr Expr
}

func (p *Print) Accept(v VisitorStmt) interface{} {
	return v.VisitPrint(p)
}

type ExprStmt struct {
	Expr Expr
}

func (es *ExprStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitExprStmt(es)
}
