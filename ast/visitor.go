package ast

import (
	"fmt"
	"strings"
)

type Visitor interface {
	visitBinary(*Binary) interface{}
	visitGrouping(*Grouping) interface{}
	visitLiteral(*Literal) interface{}
	visitUnary(*Unary) interface{}
}

type AstPrinter struct{}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.accept(p).(string))
	}
	builder.WriteString(")")

	return builder.String()
}

func (p AstPrinter) visitBinary(b *Binary) interface{} {
	return p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (p AstPrinter) visitGrouping(g *Grouping) interface{} {
	return p.parenthesize("group", g.Expr)
}

func (p AstPrinter) visitLiteral(l *Literal) interface{} {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprint(l.Value)
}

func (p AstPrinter) visitUnary(u *Unary) interface{} {
	return p.parenthesize(u.Operator.Lexeme, u.Expr)
}

func (p AstPrinter) Print(expr Expr) string {
	return expr.accept(p).(string)
}
