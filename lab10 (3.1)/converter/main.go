package main

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/predtable"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
)

const filepath = "test_files/basic.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, errTDP := parser.TopDownParse(scanner)
	if errTDP != nil {
		log.Fatalf("err in TopDownParse %+v", errTDP)
	}

	tree.Print()

	scanner.GetCompiler().PrintMessages()

	sem := semantic.NewSemantic(tree)

	rules, errSem := sem.StartSemanticAnalysis()
	if errSem != nil {
		log.Fatalf("err in semantic %+v", errSem)
	}

	rules.Print()

	genTable, errT := predtable.GenTable(rules)
	if errT != nil {
		log.Fatalf("err in gen table: %+v", errSem)
	}

	predtable.PrintGenTable(genTable)

	// TODO: сохранять эту структуру в gen_table.go

	fmt.Println("finish")
}
