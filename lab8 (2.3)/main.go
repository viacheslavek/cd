package main

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
)

const filepath = "test_files/mixed.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	token := scanner.NextToken()
	for token.GetType() != lexer.EopTag {
		if token.GetType() != lexer.ErrTag {
			fmt.Println(token)
		}
		token = scanner.NextToken()
	}

	scanner.GetCompiler().PrintMessages()

	scanner.GetCompiler().PrintIdentifiers()

	fmt.Println("finish")
}
