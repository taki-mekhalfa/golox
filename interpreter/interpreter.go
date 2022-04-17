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

func (i *Interpreter) Init() {
	// tracks the global scope
	i.globals = newEnvironment(nil)
	// starts up from the global scope and tracks the
	// current scope when entering/exiting scopes
	i.env = i.globals

	i.globals.define("clock", clockFn)
}

func (i *Interpreter) VisitWhile(while *While) interface{} {
	for truthness(i.evaluate(while.Condition)) {
		while.Body.Accept(i)
	}

	return nil
}

func (i *Interpreter) VisitIf(if_ *If) interface{} {
	if truthness(i.evaluate(if_.Condition)) {
		return if_.Then.Accept(i)
	}

	if if_.Else != nil {
		return if_.Else.Accept(i)
	}
	return nil
}

func (i *Interpreter) VisitBlock(b *Block) interface{} {
	// save current environment to recover back later
	previous := i.env
	defer func() {
		// pop the current env
		// we run this in a defer to get back the previous envs
		// even in case of a runtime error.
		// this is important when we use the prompt.
		i.env = previous
	}()

	// create a new environment inside the current one
	env := newEnvironment(i.env)
	i.env = env
	// interpret what's inside
	for _, stmt := range b.Content {
		stmt.Accept(i)
	}
	return nil
}

func (i *Interpreter) VisitExprStmt(es *ExprStmt) interface{} {
	_ = i.evaluate(es.Expr)

	return nil
}

func (i *Interpreter) VisitPrint(printExpr *Print) interface{} {
	fmt.Println(fmt.Sprint(i.evaluate(printExpr.Expr)))

	return nil
}

func (i *Interpreter) VisitFunction(f *Function) interface{} {
	i.env.define(f.Name.Lexeme, &function{declaration: f})

	return nil
}

func (i *Interpreter) VisitVarStmt(var_ *VarStmt) interface{} {
	if var_.Initializer == nil {
		i.env.define(var_.Name, nil)
	} else {
		i.env.define(var_.Name, var_.Initializer.Accept(i))
	}

	return nil
}

func (i *Interpreter) VisitBinary(b *Binary) interface{} {
	left, right := b.Left.Accept(i), b.Right.Accept(i)

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

func (i *Interpreter) VisitLogical(l *Logical) interface{} {
	switch l.Operator.Type {
	// golang will take care of short-circuiting both operators
	case token.AND:
		return truthness(i.evaluate(l.Left)) && truthness(i.evaluate(l.Right))
	case token.OR:
		return truthness(i.evaluate(l.Left)) || truthness(i.evaluate(l.Right))
	}

	// should not happen
	return nil
}

func (i *Interpreter) VisitGrouping(g *Grouping) interface{} {
	return g.Expr.Accept(i)
}

func (i *Interpreter) VisitAssign(a *Assign) interface{} {
	if _, defined := i.env.get(a.Identifier.Lexeme); !defined {
		panic(runtimeError{
			token: a.Identifier,
			msg:   fmt.Sprintf("Undefined variable '" + a.Identifier.Lexeme + "'."),
		})
	}
	v := i.evaluate(a.Value)
	i.env.assign(a.Identifier.Lexeme, v)
	return v
}

func (i *Interpreter) VisitVar(var_ *Var) interface{} {
	v, defined := i.env.get(var_.Token.Lexeme)
	if !defined {
		panic(runtimeError{
			token: var_.Token,
			msg:   fmt.Sprintf("Undefined variable '" + var_.Token.Lexeme + "'."),
		})
	}

	return v
}

func (i *Interpreter) VisitLiteral(l *Literal) interface{} {
	return l.Value
}

func (i *Interpreter) VisitCall(c *Call) interface{} {
	callee, ok := i.evaluate(c.Callee).(callable)
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
		args = append(args, i.evaluate(arg))
	}

	return callee.call(i, args)
}

func (p *Interpreter) VisitUnary(u *Unary) interface{} {
	v := u.Expr.Accept(p)

func (i *Interpreter) VisitUnary(u *Unary) interface{} {
	v := u.Expr.Accept(i)
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

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) Interpret(stmts []Stmt) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if runtimeErr, ok := err.(runtimeError); ok {
			i.ErrorCount++
			i.Error(runtimeErr.token.Line, runtimeErr.msg)
			return
		}
		panic(err)
	}()

	for _, stmt := range stmts {
		stmt.Accept(i)
	}
}

func (i *Interpreter) ResetErrors() {
	i.ErrorCount = 0
}
