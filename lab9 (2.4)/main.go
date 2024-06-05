package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab9/recdesc"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

const filepath = "test_files/develop.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	scanner.PrintTokens()
	fmt.Println()

	rp := recdesc.NewParser(scanner)
	program := rp.Parse()
	program.Print()

	fmt.Println("\n\n\nfinish")
}
