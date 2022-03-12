package scanner

import "github.com/taki-mekhalfa/golox/token"

type Scanner struct {
	src string

	Error      func(line int, errMessage string)
	ErrorCount int

	line       int
	startPos   int
	currentPos int

	tokens []token.Token
}

const (
	EOF   = rune(0)
	NL    = rune('\n')
	SPACE = rune(' ')
	TAB   = rune('\t')
	CR    = rune('\r')
)

func (s *Scanner) Init(src string) {
	s.src = src
	s.line = 1
}

func (s *Scanner) Tokens() []token.Token {
	return s.tokens
}

func (s *Scanner) Scan() {
	// while we did not consume the entire source
	for !s.isAtEnd() {
		s.startPos = s.currentPos

		c := s.next()
		switch c {
		case '(':
			s.appendToken(token.LEFT_PAREN)
		case ')':
			s.appendToken(token.RIGHT_PAREN)
		case '{':
			s.appendToken(token.LEFT_BRACE)
		case '}':
			s.appendToken(token.RIGHT_BRACE)
		case ',':
			s.appendToken(token.COMMA)
		case '.':
			s.appendToken(token.DOT)
		case '-':
			s.appendToken(token.MINUS)
		case '+':
			s.appendToken(token.PLUS)
		case ';':
			s.appendToken(token.SEMICOLON)
		case '*':
			s.appendToken(token.STAR)
		case '!':
			if s.match('=') {
				s.appendToken(token.BANG_EQUAL)
			} else {
				s.appendToken(token.BANG)
			}

		case '=':
			if s.match('=') {
				s.appendToken(token.EQUAL_EQUAL)
			} else {
				s.appendToken(token.EQUAL)
			}
		case '<':
			if s.match('=') {
				s.appendToken(token.LESS_EQUAL)
			} else {
				s.appendToken(token.LESS)
			}
		case '>':
			if s.match('=') {
				s.appendToken(token.GREATER_EQUAL)
			} else {
				s.appendToken(token.GREATER)
			}
		case '/':
			if s.match('/') {
				// Ignore line comments
				for s.peek() != NL && !s.isAtEnd() {
					s.next()
				}
			} else if s.match('*') {
				// Ignore multiline comments
				s.scanMultiLineComments()
			} else {
				s.appendToken(token.SLASH)
			}
		case SPACE, TAB, CR:
			// Ignore white space
		case NL:
			s.line++
		case '"':
			s.scanString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.scanNumber()
		default:
			switch {
			case isAlpha(c):
				s.scanIdentifier()
			default:
				s.ErrorCount++
				s.Error(s.line, "Unexpected character.")
			}
		}
	}

	s.tokens = append(s.tokens, token.Token{Type: token.EOF, Lexeme: "", Line: s.line})
}

func (s *Scanner) appendToken(typ token.Type) {
	s.tokens = append(s.tokens, token.Token{
		Type:   typ,
		Line:   s.line,
		Lexeme: s.src[s.startPos:s.currentPos],
	})
}

func (s *Scanner) scanString() {
	for {
		if s.isAtEnd() {
			s.ErrorCount++
			s.Error(s.line, "Unterminated string.")
			break
		}
		next := s.next()
		if next == NL {
			s.line++
		} else if next == '"' {
			s.appendToken(token.STRING)
			break
		}
	}
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.next()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.next()
		s.next()
	}
	for isDigit(s.peek()) {
		s.next()
	}
	s.appendToken(token.NUMBER)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.next()
	}

	if typ, ok := token.KeyWords[s.src[s.startPos:s.currentPos]]; ok {
		s.appendToken(typ)
		return
	}
	s.appendToken(token.IDENTIFIER)
}

func (s *Scanner) scanMultiLineComments() {
	for {
		if s.isAtEnd() {
			s.ErrorCount++
			s.Error(s.line, "Unterminated multiline comment.")
			break
		}

		next := s.next()
		if next == NL {
			s.line++
		} else if next == '*' && s.peek() == '/' {
			s.next()
			break
		}
	}

}
