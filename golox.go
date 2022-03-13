package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/scanner"
	"github.com/taki-mekhalfa/golox/visitor"
)

const EX_USAGE = 64
const EX_DATAERR = 65

var syntaxErrFunc = func(line int, errMessage string) {
	fmt.Printf("[line %d] Syntax Error: %s\n", line, errMessage)
}

var runtimeErrFunc = func(line int, errMessage string) {
	fmt.Printf("[line %d] Runtime Error: %s\n", line, errMessage)
}

var interpreter = visitor.Interpreter{Error: runtimeErrFunc}

func run(code string) {
	scanner := scanner.Scanner{Error: syntaxErrFunc}
	scanner.Init(code)
	scanner.Scan()
	if scanner.ErrorCount != 0 {
		return
	}

	parser := parser.Parser{Error: syntaxErrFunc}
	parser.Init(scanner.Tokens())

	expr := parser.Parse()
	if parser.ErrorCount != 0 {
		return
	}

	result := interpreter.Interpret(expr)
	if interpreter.ErrorCount == 0 {
		fmt.Println(result)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				panic(err)
			}
			break
		}
		run(scanner.Text())
		interpreter.ResetErrors()
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(EX_USAGE)
	}

	if len(os.Args) == 2 {
		// TODO

		// b, err := ioutil.ReadFile(os.Args[1])
		// if err != nil {
		// 	fmt.Printf("Could not read the source file: %+v", err)
		// 	os.Exit(1)
		// }
		// if err := run(string(b)); err != nil {
		// 	os.Exit(EX_DATAERR)
		// }
	} else {
		runPrompt()
	}
}
