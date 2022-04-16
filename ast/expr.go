package ast

import "github.com/taki-mekhalfa/golox/token"

type Expr interface {
	Accept(VisitorExpr) interface{}
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(v VisitorExpr) interface{} {
	return v.VisitBinary(b)
}

type Grouping struct {
	Expr Expr
}

func (g *Grouping) Accept(v VisitorExpr) interface{} {
	return v.VisitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(v VisitorExpr) interface{} {
	return v.VisitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Expr     Expr
}

func (u *Unary) Accept(v VisitorExpr) interface{} {
	return v.VisitUnary(u)
}

type Var struct {
	Token token.Token
}

func (var_ *Var) Accept(v VisitorExpr) interface{} {
	return v.VisitVar(var_)
}

type Assign struct {
	Identifier token.Token
	Value      Expr
}

func (a *Assign) Accept(v VisitorExpr) interface{} {
	return v.VisitAssign(a)
}

type Logical struct {
	Operator token.Token
	Left     Expr
	Right    Expr
}

func (l *Logical) Accept(v VisitorExpr) interface{} {
	return v.VisitLogical(l)
}

type Call struct {
	Callee        Expr
	ClosingParent token.Token
	Args          []Expr
}

func (c *Call) Accept(v VisitorExpr) interface{} {
	return v.VisitCall(c)
}
