package main

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

const filepath = "test_files/develop.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	token := scanner.NextToken()

	for token.GetType() != lexer.EopTag {
		fmt.Println(token)
		token = scanner.NextToken()
	}

	scanner.GetCompiler().PrintMessages()

	fmt.Println("finish")
}
