package parser

import (
	"fmt"
	"strconv"

	"github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/token"
)

type Parser struct {
	src []token.Token

	Error      func(line int, errMessage string)
	ErrorCount int

	current int

	stmts []ast.Stmt
}

func (p *Parser) Init(src []token.Token) {
	p.src = src
}

func (p *Parser) Parse() []ast.Stmt {
	for !p.isAtEnd() {
		stmt, _ := p.declaration()
		p.stmts = append(p.stmts, stmt)
	}
	return p.stmts
}

func (p *Parser) declaration() (ast.Stmt, error) {
	var stmt ast.Stmt
	var err error
	if p.match(token.VAR) {
		stmt, err = p.var_()
	} else {
		stmt, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		switch p.peek().Type {
		case token.SEMICOLON:
			p.next()
			return
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN, token.LEFT_BRACE:
			return
		default:
			p.next()
		}
	}
}

func (p *Parser) var_() (ast.Stmt, error) {
	if p.peek().Type != token.IDENTIFIER {
		p.reportError(p.peek().Line, "Expected identifier after var.")
		return nil, fmt.Errorf("line %d: expected identifier after var", p.peek().Line)
	}

	varName := p.next().Lexeme
	var initializer ast.Expr
	var err error
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(token.SEMICOLON) {
		p.reportError(p.peek().Line, "Expected ; after variable declaration.")
		return nil, fmt.Errorf("line %d: expected ; after variable declaration", p.peek().Line)
	}

	return &ast.VarStmt{Name: varName, Initializer: initializer}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return p.block()
	}
	if p.match(token.IF) {
		return p.if_()
	}

	return p.expressionStmt()
}

func (p *Parser) if_() (ast.Stmt, error) {
	if !p.match(token.LEFT_PAREN) {
		p.reportError(p.peek().Line, "Expected ( if.")
		return nil, fmt.Errorf("line %d: expected ( after if", p.peek().Line)
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.match(token.RIGHT_PAREN) {
		p.reportError(p.peek().Line, "Expected ) if condition.")
		return nil, fmt.Errorf("line %d: expected ) after if condition", p.peek().Line)
	}
	then, err := p.statement()
	if err != nil {
		return nil, err
	}
	var else_ ast.Stmt
	if p.match(token.ELSE) {
		else_, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &ast.If{Condition: condition, Then: then, Else: else_}, nil
}

func (p *Parser) block() (ast.Stmt, error) {
	var content []ast.Stmt
	for !p.isAtEnd() && p.peek().Type != token.RIGHT_BRACE {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		content = append(content, stmt)
	}
	if p.peek().Type != token.RIGHT_BRACE {
		p.reportError(p.peek().Line, "Expected } after block.")
		return nil, fmt.Errorf("line %d: expected ; after block", p.peek().Line)
	}

	// consume the '}'
	p.next()
	return &ast.Block{Content: content}, nil
}

func (p *Parser) printStmt() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != token.SEMICOLON {
		p.reportError(p.peek().Line, "Expected ; after expression.")
		return nil, fmt.Errorf("line %d: expected ; after expression", p.peek().Line)
	}
	p.next()
	return &ast.Print{Expr: expr}, nil
}

func (p *Parser) expressionStmt() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != token.SEMICOLON {
		p.reportError(p.peek().Line, "Expected ; after expression.")
		return nil, fmt.Errorf("line %d: expected ; after expression", p.peek().Line)
	}
	p.next()
	return &ast.ExprStmt{Expr: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.peek().Type == token.EQUAL {
		equal := p.next()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if var_, ok := expr.(*ast.Var); ok {
			return &ast.Assign{Identifier: var_.Token, Value: value}, nil
		} else {
			p.reportError(equal.Line, "Invalid assignment target.")
		}
	}
	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	left, err := p.and()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.OR:
			op := p.next()
			right, err := p.and()
			if err != nil {
				return nil, err
			}
			left = &ast.Logical{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}

	return left, nil
}

func (p *Parser) and() (ast.Expr, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.AND:
			op := p.next()
			right, err := p.equality()
			if err != nil {
				return nil, err
			}
			left = &ast.Logical{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}

	return left, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	left, err := p.comparison()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.EQUAL_EQUAL, token.BANG_EQUAL:
			op := p.next()
			right, err := p.comparison()
			if err != nil {
				return nil, err
			}
			left = &ast.Binary{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}
	return left, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	left, err := p.term()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL:
			op := p.next()
			right, err := p.term()
			if err != nil {
				return nil, err
			}
			left = &ast.Binary{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}
	return left, nil
}

func (p *Parser) term() (ast.Expr, error) {
	left, err := p.factor()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.MINUS, token.PLUS:
			op := p.next()
			right, err := p.factor()
			if err != nil {
				return nil, err
			}
			left = &ast.Binary{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}
	return left, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
LOOP:
	for {
		switch p.peek().Type {
		case token.SLASH, token.STAR:
			op := p.next()
			right, err := p.unary()
			if err != nil {
				return nil, err
			}
			left = &ast.Binary{Left: left, Operator: op, Right: right}
		default:
			break LOOP
		}
	}
	return left, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	switch p.peek().Type {
	case token.BANG, token.MINUS:
		op := p.next()
		unary, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: op, Expr: unary}, nil
	default:
		return p.primary()
	}
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	}
	if p.peek().Type == token.IDENTIFIER {
		return &ast.Var{Token: p.next()}, nil
	}
	if p.match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	}
	if p.peek().Type == token.STRING {
		lexeme := p.next().Lexeme
		return &ast.Literal{Value: lexeme[1 : len(lexeme)-1]}, nil
	}
	if p.peek().Type == token.NUMBER {
		lexeme := p.next().Lexeme
		// ignore error as this is guaranteed to be a valid float after scanning
		number, _ := strconv.ParseFloat(lexeme, 64)
		return &ast.Literal{Value: number}, nil
	}

	// This should be a left paren
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.peek().Type != token.RIGHT_PAREN {
			p.reportError(p.peek().Line, "Expected ) after expression.")
			return nil, fmt.Errorf("line %d: expected ) after expression", p.peek().Line)
		}
		p.next()
		return &ast.Grouping{Expr: expr}, nil
	}

	p.reportError(p.peek().Line, "Expected an expression.")
	return nil, fmt.Errorf("line %d: expected an expression", p.peek().Line)
}
