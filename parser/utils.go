package parser

import "github.com/taki-mekhalfa/golox/token"

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() token.Token {
	return p.src[p.current]
}

func (p *Parser) next() token.Token {
	p.current++
	return p.src[p.current-1]
}

func (p *Parser) match(typ token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	if p.peek().Type == typ {
		p.next()
		return true
	}

	return false
}

func (p *Parser) reportError(line int, errMessage string) {
	p.ErrorCount++
	p.Error(line, errMessage)
}
