package main

import (
	"fmt"

	. "github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/token"
	"github.com/taki-mekhalfa/golox/visitor"
)

func main() {
	expr := &Binary{
		Left: &Unary{
			Operator: token.Token{Type: token.MINUS, Lexeme: "-", Line: 1},
			Expr:     &Literal{Value: 123},
		},
		Operator: token.Token{Type: token.STAR, Lexeme: "*", Line: 1},
		Right: &Grouping{
			Expr: &Literal{Value: 45.67},
		},
	}

	printer := visitor.PrettyPrinter{}
	fmt.Println(printer.Print(expr))
}
