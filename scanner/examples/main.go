package main

import (
	"fmt"

	"github.com/taki-mekhalfa/golox/scanner"
)

func main() {
	scanner := scanner.Scanner{Error: func(line int, errMessage string) {
		fmt.Printf("[line %d] Error: %s\n", line, errMessage)
	}}

	src := `
		// THIS IS A COMMENT
		var numTrue_or33_nil = 45345;
		var pi = 3.14; // delicious Pi
		/*
			This is a multi line comment haha
			you should stop here 
		*/
		fun hello(a,b) {
			return a + (b * pi);
		}
	`

	scanner.Init(src)
	scanner.Scan()
	fmt.Println(scanner.Tokens())
}
