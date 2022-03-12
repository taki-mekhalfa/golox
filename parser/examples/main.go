package main

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/ast"
	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/scanner"
)

func main() {
	scanner := scanner.Scanner{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}

	src := `4 * 6 / (2+1) + (2) `

	scanner.Init(src)
	scanner.Scan()

	fmt.Println(scanner.Tokens())

	parser := parser.Parser{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}
	parser.Init(scanner.Tokens())

	expr, err := parser.Parse()
	if err == nil {
		printer := ast.AstPrinter{}
		fmt.Println(printer.Print(expr))
	}

}
