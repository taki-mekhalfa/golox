package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const EX_USAGE = 64
const EX_DATAERR = 65

var hadError bool

func run(code string) (err error) {
	// Read tokens
	return
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			return scanner.Err()
		}
		run(scanner.Text())
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(EX_USAGE)
	}

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
