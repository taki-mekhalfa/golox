package ast

import "github.com/taki-mekhalfa/golox/token"

type Expr interface {
	accept(Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) accept(v Visitor) interface{} {
	return v.visitBinary(b)
}

type Grouping struct {
	Expr Expr
}

func (g *Grouping) accept(v Visitor) interface{} {
	return v.visitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (l *Literal) accept(v Visitor) interface{} {
	return v.visitLiteral(l)
}

type Unary struct {
	Operator token.Token
	Expr     Expr
}

func (u *Unary) accept(v Visitor) interface{} {
	return v.visitUnary(u)
}
