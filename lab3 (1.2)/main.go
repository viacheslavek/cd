package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab3/lab_lexer"
)

const filepath = "test_files/mixed.txt"

func main() {

	lexer := lab_lexer.NewLexer(filepath)

	token := lexer.NextToken()
	for token.IsToken() || token.IsError() && token.CurrentType() != lab_lexer.EOF {
		fmt.Println(token)
		token = lexer.NextToken()
	}

	fmt.Println("finish")
}
