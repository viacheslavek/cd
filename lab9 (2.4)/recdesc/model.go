package recdesc

import (
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

type RecursiveParser struct {
	currentToken lexer.IToken
	scanner      *lexer.Scanner
}

func NewParser(scanner *lexer.Scanner) *RecursiveParser {
	return &RecursiveParser{currentToken: scanner.NextToken(), scanner: scanner}
}

func (rp *RecursiveParser) isExpectedToken(tokenValue string, tokenType lexer.DomainTag) {
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	if t.GetValue() != tokenValue {
		log.Fatalf("incorrect token value. Expected: '%s' given: '%s'", tokenValue, t.GetValue())
	}
	if t.GetType() != tokenType {
		log.Fatalf("incorrect token type. Expected: '%s' given: '%s'",
			lexer.TagToString[tokenType], lexer.TagToString[t.GetType()])
	}
}
