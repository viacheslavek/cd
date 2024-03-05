package main

import (
	"fmt"
	"lab4/lab_lexer"
)

const filepath = "test_files/ident_correct.txt"

func main() {

	//scanner := lab_lexer.NewScanner(filepath)

	tokens := lab_lexer.ParseFile(filepath)

	fmt.Println("tokens", tokens)

	//token := scanner.NextToken()
	// TODO: сделать итерацию по токенам
	//for token.IsToken() || token.IsError() && token.CurrentType() != lab_lexer.EOF {
	//	fmt.Println(token)
	//	token = lexer.NextToken()
	//}

	fmt.Println("finish")
}
