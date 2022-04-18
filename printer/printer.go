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
		builder.WriteString(p.PrintExpr(expr))
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
	return fmt.Sprintf("[%s]", var_.Token.Lexeme)
}

func (p PrettyPrinter) VisitCall(c *Call) interface{} {
	var builder strings.Builder
	builder.WriteString(p.PrintExpr(c.Callee))
	builder.WriteString("(")
	for _, arg := range c.Args {
		builder.WriteString(p.PrintExpr(arg))
		builder.WriteString(",")
	}
	builder.WriteString(")")
	return builder.String()
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

	return fmt.Sprintf("var %s = %s", var_.Name, p.PrintExpr(var_.Initializer))
}

func (p PrettyPrinter) VisitBlock(b *Block) interface{} {
	var builder strings.Builder
	builder.WriteString("{\n")
	for _, stmt := range b.Content {
		builder.WriteString(p.PrintStmt(stmt))
		builder.WriteString("\n")
	}
	builder.WriteString("}")

	return builder.String()
}

func (p PrettyPrinter) VisitIf(if_ *If) interface{} {
	var builder strings.Builder
	builder.WriteString("if (")
	builder.WriteString(p.PrintExpr(if_.Condition))
	builder.WriteString(" ) then ")
	builder.WriteString(p.PrintStmt(if_.Then))
	if if_.Else != nil {
		builder.WriteString("\nelse ")
		builder.WriteString(p.PrintStmt(if_.Else))
	}
	return builder.String()
}

func (p PrettyPrinter) VisitWhile(while *While) interface{} {
	var builder strings.Builder
	builder.WriteString("while (")
	builder.WriteString(p.PrintExpr(while.Condition))
	builder.WriteString(") ")
	builder.WriteString(p.PrintStmt(while.Body))
	return builder.String()
}

func (p PrettyPrinter) VisitFunction(f *Function) interface{} {
	var builder strings.Builder
	builder.WriteString("fun ")
	builder.WriteString(f.Name.Lexeme)
	builder.WriteString("(")
	for _, param := range f.Params {
		builder.WriteString(param.Lexeme)
		builder.WriteString(",")
	}
	builder.WriteString(") ")
	builder.WriteString(p.PrintStmt(&Block{Content: f.Body}))
	return builder.String()
}

func (p PrettyPrinter) VisitReturn(r *Return) interface{} {
	return "return " + p.PrintExpr(r.Value)
}

func (p PrettyPrinter) PrintExpr(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p PrettyPrinter) PrintStmt(stmt Stmt) string {
	return stmt.Accept(p).(string)
}
