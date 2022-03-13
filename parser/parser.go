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
}

func (p *Parser) Init(src []token.Token) {
	p.src = src
}

func (p *Parser) Parse() ast.Expr {
	expr, _ := p.expression()
	if !p.isAtEnd() {
		p.reportError(p.peek().Line, fmt.Sprintf("unexpected %s, expecting EOF", p.peek().Lexeme))
		return nil
	}
	return expr
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
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
