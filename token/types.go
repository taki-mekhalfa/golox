package token

type Type int

//go:generate stringer -type=Type
const (
	// Single-character tokens.
	LEFT_PAREN Type = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL

	EQUAL
	EQUAL_EQUAL

	GREATER
	GREATER_EQUAL

	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	CLASS
	VAR
	PRINT
	NIL
	FUN
	RETURN
	SUPER
	THIS

	AND
	OR
	IF
	ELSE
	FALSE
	TRUE

	FOR
	WHILE

	EOF
)

var KeyWords = map[string]Type{
	"class":  CLASS,
	"var":    VAR,
	"print":  PRINT,
	"nil":    NIL,
	"fun":    FUN,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"and":    AND,
	"or":     OR,
	"if":     IF,
	"else":   ELSE,
	"false":  FALSE,
	"true":   TRUE,
	"for":    FOR,
	"while":  WHILE,
}
