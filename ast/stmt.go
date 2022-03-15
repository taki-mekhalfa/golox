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

type VarStmt struct {
	Name        string
	Initializer Expr
}

func (var_ *VarStmt) Accept(v VisitorStmt) interface{} {
	return v.VisitVarStmt(var_)
}

type Block struct {
	Content []Stmt
}

func (b *Block) Accept(v VisitorStmt) interface{} {
	return v.VisitBlock(b)
}
