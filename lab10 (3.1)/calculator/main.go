package main

import (
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/top_down_parse"
)

const filepath = "example.txt"

func main() {
	log.Println("start calculator")

	// TODO: пишу лексер для калькулятора -> +, *, (, ), INT
	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, errTDP := parser.TopDownParse(scanner)
	if errTDP != nil {
		log.Fatalf("err in TopDownParse %+v", errTDP)
	}

	tree.Print()

	// TODO: делаю интерпретатор выражения по дереву

}
