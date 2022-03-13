package main

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/scanner"
	"github.com/taki-mekhalfa/golox/visitor"
)

func main() {
	errFunc := func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}

	runtimeErrFunc := func(line int, errMessage string) {
		fmt.Printf("[line %d] Runtime Error: %s\n", line, errMessage)
	}

	scanner := scanner.Scanner{Error: errFunc}

	src := `2 * (3/-"muffin")`

	scanner.Init(src)
	scanner.Scan()

	fmt.Println(scanner.Tokens())

	parser := parser.Parser{Error: errFunc}
	parser.Init(scanner.Tokens())

	expr := parser.Parse()
	if parser.ErrorCount == 0 {
		printer := visitor.PrettyPrinter{}
		fmt.Println(printer.Print(expr))
		interpreter := visitor.Interpreter{Error: runtimeErrFunc}
		result := interpreter.Interpret(expr)
		if interpreter.ErrorCount == 0 {
			fmt.Println(result)
		}
	}
}
