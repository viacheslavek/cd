package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
)

const filepath = "test_files/mixed.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, err := parser.TopDownParse(scanner)
	if err != nil {
		log.Panic("пупу:", err)
	}

	tree.Print()

	scanner.GetCompiler().PrintMessages()

	sem := semantic.NewSemantic(tree)

	errSem := sem.StartSemanticAnalysis()
	if errSem != nil {
		fmt.Println("err in sem", errSem)
	}

	fmt.Println("finish")
}
