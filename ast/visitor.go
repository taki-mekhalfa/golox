package ast

type VisitorExpr interface {
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitUnary(*Unary) interface{}
	VisitVar(*Var) interface{}
	VisitAssign(*Assign) interface{}
	VisitLogical(*Logical) interface{}
	VisitCall(*Call) interface{}
	VisitGet(*Get) interface{}
	VisitSet(*Set) interface{}
}

type VisitorStmt interface {
	VisitPrint(*Print) interface{}
	VisitExprStmt(*ExprStmt) interface{}
	VisitVarStmt(*VarStmt) interface{}
	VisitBlock(*Block) interface{}
	VisitIf(*If) interface{}
	VisitWhile(*While) interface{}
	VisitFunction(f *Function) interface{}
	VisitReturn(r *Return) interface{}
	VisitClass(c *Class) interface{}
}
