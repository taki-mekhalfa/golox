package token

import "fmt"

type Token struct {
	Type   Type
	Lexeme string
	Line   int
}

func (t Token) String() string {
	return fmt.Sprintf("[%s:%s]", t.Type, t.Lexeme)
}
