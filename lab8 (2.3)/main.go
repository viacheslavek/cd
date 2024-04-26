package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab8/top_down_parse"
	"log"
)

const filepath = "test_files/mixed.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, err := parser.TopDownParse(scanner)
	if err != nil {
		log.Panic("пупу:", err)
	}

	tree.PrintNode()

	scanner.GetCompiler().PrintMessages()

	fmt.Println("finish")
}
