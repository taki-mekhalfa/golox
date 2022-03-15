package ast

type VisitorExpr interface {
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitUnary(*Unary) interface{}
	VisitVar(*Var) interface{}
	VisitAssign(*Assign) interface{}
}

type VisitorStmt interface {
	VisitPrint(*Print) interface{}
	VisitExprStmt(*ExprStmt) interface{}
	VisitVarStmt(*VarStmt) interface{}
	VisitBlock(*Block) interface{}
}
