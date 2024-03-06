package main

import (
	"fmt"
	"lab4/lab_lexer"
)

const filepath = "test_files/ident_error.txt"

func main() {

	scanner := lab_lexer.NewScanner(filepath)

	token := scanner.NextToken()
	for token.GetType() != lab_lexer.EopTag {
		if token.GetType() != lab_lexer.ErrTag {
			fmt.Println(token)
		}
		token = scanner.NextToken()
	}

	scanner.GetCompiler().PrintMessages()

	scanner.GetCompiler().PrintIdentifiers()

	fmt.Println("finish")
}
