package main

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/scanner"
	"github.com/taki-mekhalfa/golox/visitor"
)

func main() {
	scanner := scanner.Scanner{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}

	src := `
		print 1 + 1;
		print "ok";
		print "ok" + "boki";
		1 + (2/3);
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
		printer := visitor.PrettyPrinter{}
		for _, stmt := range stmts {
			fmt.Println(printer.PrintStmt(stmt))
		}
	}

}
