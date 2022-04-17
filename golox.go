package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/taki-mekhalfa/golox/interpreter"
	"github.com/taki-mekhalfa/golox/parser"
	"github.com/taki-mekhalfa/golox/scanner"
)

const EX_USAGE = 64
const EX_DATAERR = 65

var syntaxErrFunc = func(line int, errMessage string) {
	fmt.Printf("[line %d] Syntax Error: %s\n", line, errMessage)
}

var runtimeErrFunc = func(line int, errMessage string) {
	fmt.Printf("[line %d] Runtime Error: %s\n", line, errMessage)
}

var interpreter_ = interpreter.Interpreter{Error: runtimeErrFunc}

func run(code string) error {
	scanner := scanner.Scanner{Error: syntaxErrFunc}
	scanner.Init(code)
	scanner.Scan()
	if scanner.ErrorCount != 0 {
		return fmt.Errorf("encountred %d scanner errors", scanner.ErrorCount)
	}

	parser := parser.Parser{Error: syntaxErrFunc}
	parser.Init(scanner.Tokens())

	stmts := parser.Parse()
	if parser.ErrorCount != 0 {
		return fmt.Errorf("encountred %d parser errors", parser.ErrorCount)
	}

	interpreter_.Interpret(stmts)
	if interpreter_.ErrorCount != 0 {
		return fmt.Errorf("encountred %d interpreter errors", parser.ErrorCount)
	}

	return nil
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
		_ = run(scanner.Text())
		interpreter_.ResetErrors()
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(EX_USAGE)
	}

	interpreter_.Init()

	if len(os.Args) == 2 {
		b, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Could not read the source file: %+v", err)
			os.Exit(1)
		}
		if err := run(string(b)); err != nil {
			os.Exit(EX_DATAERR)
		}
	} else {
		runPrompt()
	}
}
