package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/interpreteter"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/top_down_parse"
)

const filepath = "example.txt"

func main() {
	log.Println("start calculator")

	scanner := lexer.NewScanner(filepath)

	scanner.PrintTokens()

	parser := top_down_parse.NewParser()

	tree, errTDP := parser.TopDownParse(scanner)
	if errTDP != nil {
		log.Fatalf("err in TopDownParse %+v", errTDP)
	}

	tree.Print()

	fmt.Println("solver:", interpreteter.Solve(tree))

	log.Println("end calculator")
}
