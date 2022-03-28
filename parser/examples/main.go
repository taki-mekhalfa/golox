package main

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/printer"
	"github.com/taki-mekhalfa/golox/scanner"
)

func main() {
	scanner := scanner.Scanner{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}

	src := `
		var i = 33;
		print i;
		for (var i = 0; i < 3; i = i + 1){
			print i;
		}
	`

	scanner.Init(src)
	scanner.Scan()

	fmt.Println(scanner.Tokens())

	parser := parser.Parser{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}
	parser.Init(scanner.Tokens())

	stmts := parser.Parse()
	if parser.ErrorCount == 0 {
		printer := printer.PrettyPrinter{}
		for _, stmt := range stmts {
			fmt.Println(printer.PrintStmt(stmt))
		}
	}
}
