package visitor

import (
	"fmt"

	. "github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/token"
)

type runtimeError struct {
	token token.Token
	msg   string
}

type environment struct {
	values map[string]interface{}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) get(name string) (interface{}, bool) {
	v, ok := e.values[name]
	return v, ok

}

type Interpreter struct {
	Error      func(line int, errMessage string)
	ErrorCount int
	env        *environment
}

func (p *Interpreter) Init() {
	p.env = &environment{
		values: map[string]interface{}{},
	}
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
	p.env.define(a.Identifier.Lexeme, v)
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
