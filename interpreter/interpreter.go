package interpreter

import (
	"fmt"

	. "github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/token"
)

type runtimeError struct {
	token token.Token
	msg   string
}

type Interpreter struct {
	Error      func(line int, errMessage string)
	ErrorCount int
	env        *environment
	globals    *environment
}

func (p *Interpreter) Init() {
	// tracks the global scope
	p.globals = &environment{
		values: map[string]interface{}{},
	}
	// starts up from the global scope and tracks the
	// current scope when entering/exiting scopes
	p.env = p.globals

	p.globals.define("clock", clockFn)
}

func (p *Interpreter) VisitWhile(while *While) interface{} {
	for truthness(p.evaluate(while.Condition)) {
		while.Body.Accept(p)
	}

	return nil
}

func (p *Interpreter) VisitIf(if_ *If) interface{} {
	if truthness(p.evaluate(if_.Condition)) {
		return if_.Then.Accept(p)
	}

	if if_.Else != nil {
		return if_.Else.Accept(p)
	}
	return nil
}

func (p *Interpreter) VisitBlock(b *Block) interface{} {
	// create a new environment inside the current one
	env := &environment{values: map[string]interface{}{}, parent: p.env}
	p.env = env
	// interpret what's inside
	p.Interpret(b.Content)
	// pop the current env
	p.env = p.env.parent
	return nil
}

func (p *Interpreter) VisitExprStmt(es *ExprStmt) interface{} {
	_ = p.evaluate(es.Expr)

	return nil
}

func (p *Interpreter) VisitPrint(printExpr *Print) interface{} {
	fmt.Println(fmt.Sprint(p.evaluate(printExpr.Expr)))

	return nil
}

func (p *Interpreter) VisitVarStmt(var_ *VarStmt) interface{} {
	if var_.Initializer == nil {
		p.env.define(var_.Name, nil)
	} else {
		p.env.define(var_.Name, var_.Initializer.Accept(p))
	}

	return nil
}

func (p *Interpreter) VisitBinary(b *Binary) interface{} {
	left, right := b.Left.Accept(p), b.Right.Accept(p)

	switch b.Operator.Type {
	case token.STAR:
		checkNumberOperands(b.Operator, left, right)
		return left.(float64) * right.(float64)
	case token.SLASH:
		checkNumberOperands(b.Operator, left, right)
		rightNumber := right.(float64)
		checkIsNotZero(b.Operator, rightNumber)
		return left.(float64) / rightNumber
	case token.MINUS:
		checkNumberOperands(b.Operator, left, right)
		return left.(float64) - right.(float64)
	case token.PLUS:
		checkOperandsSameType(b.Operator, left, right)
		switch left.(type) {
		case float64:
			return left.(float64) + right.(float64)
		case string:

			return left.(string) + right.(string)
		}
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)

	case token.BANG_EQUAL:
		return left != right
	case token.EQUAL_EQUAL:
		return left == right
	}

	// should not happen
	return nil
}

func (p *Interpreter) VisitLogical(l *Logical) interface{} {
	switch l.Operator.Type {
	// golang will take care of short-circuiting both operators
	case token.AND:
		return truthness(p.evaluate(l.Left)) && truthness(p.evaluate(l.Right))
	case token.OR:
		return truthness(p.evaluate(l.Left)) || truthness(p.evaluate(l.Right))
	}

	// should not happen
	return nil
}

func (p *Interpreter) VisitGrouping(g *Grouping) interface{} {
	return g.Expr.Accept(p)
}

func (p *Interpreter) VisitAssign(a *Assign) interface{} {
	if _, defined := p.env.get(a.Identifier.Lexeme); !defined {
		panic(runtimeError{
			token: a.Identifier,
			msg:   fmt.Sprintf("Undefined variable '" + a.Identifier.Lexeme + "'."),
		})
	}
	v := p.evaluate(a.Value)
	p.env.assign(a.Identifier.Lexeme, v)
	return v
}

func (p *Interpreter) VisitVar(var_ *Var) interface{} {
	v, defined := p.env.get(var_.Token.Lexeme)
	if !defined {
		panic(runtimeError{
			token: var_.Token,
			msg:   fmt.Sprintf("Undefined variable '" + var_.Token.Lexeme + "'."),
		})
	}

	return v
}

func (p *Interpreter) VisitLiteral(l *Literal) interface{} {
	return l.Value
}

func (p *Interpreter) VisitCall(c *Call) interface{} {
	callee, ok := p.evaluate(c.Callee).(callable)
	if !ok {
		panic(runtimeError{
			token: c.ClosingParent,
			msg:   "Can only call functions and classes.",
		})
	}
	if len(c.Args) != callee.arity() {
		panic(runtimeError{
			token: c.ClosingParent,
			msg:   fmt.Sprintf("Expected %d arguments, but got %d.", callee.arity(), len(c.Args)),
		})
	}

	args := []interface{}{}
	for _, arg := range c.Args {
		args = append(args, p.evaluate(arg))
	}

	return callee.call(p, args)
}

func (p *Interpreter) VisitUnary(u *Unary) interface{} {
	v := u.Expr.Accept(p)
	switch u.Operator.Type {
	case token.BANG:
		return !truthness(v)
	case token.MINUS:
		checkNumberOperand(u.Operator, v)
		return -v.(float64)
	}

	// should not happen
	return nil
}

func truthness(v interface{}) bool {
	if v == nil {
		return false
	}
	if boolValue, ok := v.(bool); ok {
		return boolValue
	}

	return true
}

func checkNumberOperand(t token.Token, o interface{}) {
	if _, ok := o.(float64); !ok {
		panic(runtimeError{
			token: t,
			msg:   "Operand must be a number.",
		})
	}
}

func checkNumberOperands(t token.Token, left, right interface{}) {
	_, leftIsNumber := left.(float64)
	_, rightIsNumber := right.(float64)
	if leftIsNumber && rightIsNumber {
		return
	}

	panic(runtimeError{
		token: t,
		msg:   "Operands must be both numbers.",
	})
}

func checkOperandsSameType(t token.Token, left, right interface{}) {
	_, leftIsNumber := left.(float64)
	_, rightIsNumber := right.(float64)
	if leftIsNumber == rightIsNumber {
		return
	}

	panic(runtimeError{
		token: t,
		msg:   "Operands must be both numbers or both strings.",
	})
}

func checkIsNotZero(t token.Token, n float64) {
	if n != 0 {
		return
	}

	panic(runtimeError{
		token: t,
		msg:   "Divided by 0.",
	})
}

func (p *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(p)
}

func (p *Interpreter) Interpret(stmts []Stmt) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if runtimeErr, ok := err.(runtimeError); ok {
			p.ErrorCount++
			p.Error(runtimeErr.token.Line, runtimeErr.msg)
			return
		}
		panic(err)
	}()

	for _, stmt := range stmts {
		stmt.Accept(p)
	}
}

func (p *Interpreter) ResetErrors() {
	p.ErrorCount = 0
}
