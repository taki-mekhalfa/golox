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
