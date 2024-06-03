package main

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

const filepath = "test_files/lexer.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	//token := scanner.NextToken()
	//
	//for token.GetType() != lexer.EopTag {
	//	fmt.Println(token)
	//	token = scanner.NextToken()
	//}

	scanner.PrintTokens()

	fmt.Println("finish")
}
