package visitor

import (
	"fmt"
	"strings"

	. "github.com/taki-mekhalfa/golox/ast"
)

type PrettyPrinter struct{}

func (p PrettyPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}
	builder.WriteString(")")

	return builder.String()
}

func (p PrettyPrinter) VisitBinary(b *Binary) interface{} {
	return p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (p PrettyPrinter) VisitGrouping(g *Grouping) interface{} {
	return p.parenthesize("group", g.Expr)
}

func (p PrettyPrinter) VisitLiteral(l *Literal) interface{} {
	if l.Value == nil {
		return "nil"
	}
	return fmt.Sprint(l.Value)
}

func (p PrettyPrinter) VisitUnary(u *Unary) interface{} {
	return p.parenthesize(u.Operator.Lexeme, u.Expr)
}

func (p PrettyPrinter) VisitPrint(printStmt *Print) interface{} {
	return fmt.Sprint("PRINT ", p.PrintExpr(printStmt.Expr))
}

func (p PrettyPrinter) VisitExprStmt(exprStmt *ExprStmt) interface{} {
	return p.PrintExpr(exprStmt.Expr)
}

func (p PrettyPrinter) PrintExpr(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p PrettyPrinter) PrintStmt(stmt Stmt) string {
	return stmt.Accept(p).(string)
}
