package resolver

import (
	"fmt"

	. "github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/interpreter"
)

const (
	init_ = "init"
)

type meta struct {
	defined bool
	used    bool
	line    int
}

type Resolver struct {
	Error      func(line int, errMessage string)
	ErrorCount int

	scopes []map[string]*meta
	Interp *interpreter.Interpreter

	funcCtx     functionCtx
	insideClass bool
}

func (r *Resolver) Resolve(stmts []Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) VisitGet(g *Get) (void interface{}) {
	// we don't resolve the property name as it's dynamically looked up,
	// resolve only the object expression.
	r.resolveExpr(g.Object)
	return
}

func (r *Resolver) VisitSet(s *Set) (void interface{}) {
	// we don't resolve the property name as it's dynamically looked up,
	// resolve only the object expression.
	r.resolveExpr(s.Value)
	r.resolveExpr(s.Object)
	return
}

func (r *Resolver) VisitClass(c *Class) (void interface{}) {
	r.declare(c.Name.Lexeme, c.Name.Line)
	r.define(c.Name.Lexeme)

	r.beginScope()

	enclosedInClass := r.insideClass
	r.insideClass = true

	r.declare("this", c.Name.Line)
	r.define("this")
	// "this" is used by default to avoid errors
	// related to declared but not used variables
	r.use("this")

	for _, method := range c.Methods {
		enclosingFuncCtx := r.funcCtx
		if method.Name.Lexeme == init_ {
			r.funcCtx = initializer
		} else {
			r.funcCtx = function
		}
		r.reslveFunction(method)
		r.funcCtx = enclosingFuncCtx

		// methods are used by default to avoid errors
		// related to declared but not used variables
		r.use(method.Name.Lexeme)
	}

	r.endScope()

	r.insideClass = enclosedInClass
	return
}

func (r *Resolver) VisitBlock(b *Block) (void interface{}) {
	r.beginScope()
	for _, stmt := range b.Content {
		r.resolveStmt(stmt)
	}
	r.endScope()
	return
}

func (r *Resolver) VisitVarStmt(var_ *VarStmt) (void interface{}) {
	r.declare(var_.Name, var_.Token.Line)
	if var_.Initializer != nil {
		r.resolveExpr(var_.Initializer)
	}
	r.define(var_.Name)
	return nil
}

func (r *Resolver) VisitVar(var_ *Var) (void interface{}) {
	if meta, declared := r.currentScope()[var_.Token.Lexeme]; declared && !meta.defined {
		r.reportError(var_.Token.Line, "Can't read local variable in its own initializer.")
	}
	r.use(var_.Token.Lexeme)
	r.resolve(var_, var_.Token.Lexeme)
	return
}

func (r *Resolver) VisitAssign(a *Assign) (void interface{}) {
	r.resolveExpr(a.Value)
	r.resolve(a, a.Identifier.Lexeme)
	return
}

func (r *Resolver) reslveFunction(f *Function) {
	r.declare(f.Name.Lexeme, f.Name.Line)
	r.define(f.Name.Lexeme)

	r.beginScope()
	for _, param := range f.Params {
		r.declare(param.Lexeme, f.Name.Line)
		r.define(param.Lexeme)
	}
	for _, stmt := range f.Body {
		r.resolveStmt(stmt)
	}
	r.endScope()
}

func (r *Resolver) VisitFunction(f *Function) (void interface{}) {
	enclosingFuncCtx := r.funcCtx
	r.funcCtx = function

	r.reslveFunction(f)

	r.funcCtx = enclosingFuncCtx
	return
}

func (r *Resolver) VisitExprStmt(es *ExprStmt) (void interface{}) {
	r.resolveExpr(es.Expr)
	return
}

func (r *Resolver) VisitIf(if_ *If) (void interface{}) {
	r.resolveExpr(if_.Condition)
	r.resolveStmt(if_.Then)
	if if_.Else != nil {
		r.resolveStmt(if_.Else)
	}
	return
}

func (r *Resolver) VisitPrint(p *Print) (void interface{}) {
	r.resolveExpr(p.Expr)
	return
}

func (r *Resolver) VisitReturn(ret_ *Return) (void interface{}) {
	switch r.funcCtx {
	case initializer:
		if ret_.Value != nil {
			r.reportError(ret_.Token.Line, "Can't return a value from class initializer.")
			return
		}
	case function:
		r.reportError(ret_.Token.Line, "Can't return from top-level code.")
		return
	default:
	}

	if ret_.Value != nil {
		r.resolveExpr(ret_.Value)
	}
	return
}

func (r *Resolver) VisitWhile(while *While) (void interface{}) {
	r.resolveExpr(while.Condition)
	r.resolveStmt(while.Body)
	return
}

func (r *Resolver) VisitBinary(b *Binary) (void interface{}) {
	r.resolveExpr(b.Left)
	r.resolveExpr(b.Right)
	return
}

func (r *Resolver) VisitCall(call *Call) (void interface{}) {
	r.resolveExpr(call.Callee)
	for _, expr := range call.Args {
		r.resolveExpr(expr)
	}
	return
}

func (r *Resolver) VisitGrouping(g *Grouping) (void interface{}) {
	r.resolveExpr(g.Expr)
	return
}

func (r *Resolver) VisitLiteral(l *Literal) (void interface{}) {
	return nil
}

func (r *Resolver) VisitLogical(l *Logical) (void interface{}) {
	r.resolveExpr(l.Left)
	r.resolveExpr(l.Right)
	return
}

func (r *Resolver) VisitUnary(u *Unary) (void interface{}) {
	r.resolveExpr(u.Expr)
	return
}

func (r *Resolver) VisitThis(this *This) (void interface{}) {
	if !r.insideClass {
		r.reportError(this.Keyword.Line, "Can't use 'this' outside of a class.")
		return
	}
	r.resolve(this, "this")
	return
}

func (r *Resolver) resolve(expr Expr, name string) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name]; ok {
			r.Interp.Resolve(expr, len(r.scopes)-i-1)
			return
		}
	}
}

func (r *Resolver) use(name string) {
	if r.currentScope() == nil {
		return
	}
	if _, ok := r.currentScope()[name]; ok {
		r.currentScope()[name].used = true
	}
}

func (r *Resolver) declare(name string, line int) {
	if r.currentScope() == nil {
		return
	}
	if _, ok := r.currentScope()[name]; ok {
		r.reportError(line, "Already a variable with this name in this scope.")
		return
	}
	r.currentScope()[name] = &meta{line: line}
}

func (r *Resolver) define(name string) {
	if r.currentScope() == nil {
		return
	}

	r.currentScope()[name].defined = true
}

func (r *Resolver) resolveStmt(stmt Stmt) interface{} {
	return stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) interface{} {
	return expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]*meta{})
}

func (r *Resolver) endScope() {
	for name, meta := range r.currentScope() {
		if !meta.used {
			r.reportError(meta.line, fmt.Sprintf("%s declared but not used.", name))
		}
	}
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) currentScope() map[string]*meta {
	if len(r.scopes) == 0 {
		return nil
	}
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) reportError(line int, errMessage string) {
	r.ErrorCount++
	r.Error(line, errMessage)
}
