package printer

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

func (p PrettyPrinter) VisitLogical(l *Logical) interface{} {
	return p.parenthesize(l.Operator.Lexeme, l.Left, l.Right)
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

func (p PrettyPrinter) VisitVar(var_ *Var) interface{} {
	return fmt.Sprintf("var[%s]", var_.Token.Lexeme)
}

func (p PrettyPrinter) VisitUnary(u *Unary) interface{} {
	return p.parenthesize(u.Operator.Lexeme, u.Expr)
}

func (p PrettyPrinter) VisitPrint(printStmt *Print) interface{} {
	return fmt.Sprint("PRINT ", p.PrintExpr(printStmt.Expr))
}

func (p PrettyPrinter) VisitAssign(a *Assign) interface{} {
	return fmt.Sprintf("var[%s]=%s", a.Identifier.Lexeme, p.PrintExpr(a.Value))
}

func (p PrettyPrinter) VisitExprStmt(exprStmt *ExprStmt) interface{} {
	return p.PrintExpr(exprStmt.Expr)
}

func (p PrettyPrinter) VisitVarStmt(var_ *VarStmt) interface{} {
	if var_.Initializer == nil {
		return fmt.Sprintf("var %s", var_.Name)
	}

	return fmt.Sprintf("var %s = %s", var_.Name, var_.Initializer.Accept(p))
}

func (p PrettyPrinter) VisitBlock(b *Block) interface{} {
	var builder strings.Builder
	builder.WriteString("{\n")
	for _, stmt := range b.Content {
		builder.WriteString(stmt.Accept(p).(string))
		builder.WriteString("\n")
	}
	builder.WriteString("}")

	return builder.String()
}

func (p PrettyPrinter) VisitIf(if_ *If) interface{} {
	var builder strings.Builder
	builder.WriteString("if (")
	builder.WriteString(if_.Condition.Accept(p).(string))
	builder.WriteString(" ) then ")
	builder.WriteString(if_.Then.Accept(p).(string))
	if if_.Else != nil {
		builder.WriteString("\nelse ")
		builder.WriteString(if_.Else.Accept(p).(string))
	}
	return builder.String()
}

func (p PrettyPrinter) VisitWhile(while *While) interface{} {
	var builder strings.Builder
	builder.WriteString("while (")
	builder.WriteString(while.Condition.Accept(p).(string))
	builder.WriteString(") ")
	builder.WriteString(while.Body.Accept(p).(string))
	return builder.String()
}

func (p PrettyPrinter) PrintExpr(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p PrettyPrinter) PrintStmt(stmt Stmt) string {
	return stmt.Accept(p).(string)
}
