package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Scanner struct {
	tokens   []IToken
	position int
}

func NewScanner(filepath string) *Scanner {
	return &Scanner{
		tokens:   ParseFile(filepath),
		position: 0,
	}
}

func (s *Scanner) PrintTokens() {
	fmt.Println("TOKENS:")
	for _, p := range s.tokens {
		fmt.Println(p)
	}
}

func (s *Scanner) NextToken() IToken {
	token := s.tokens[s.position]

	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == IntTag ||
		token.GetType() == CloseBracketTag || token.GetType() == OpenBracketTag ||
		token.GetType() == PlusTag || token.GetType() == MultiplyTag {
		return token
	}

	return token
}

func ParseFile(filepath string) []IToken {
	file, errO := os.Open(filepath)
	if errO != nil {
		log.Fatalf("can't open file in parser %e", errO)
	}
	defer func() {
		_ = file.Close()
	}()

	tokens := make([]IToken, 0)

	scanner := bufio.NewScanner(file)
	runeScanner := NewRunePosition(scanner)

	for runeScanner.GetRune() != -1 {
		for runeScanner.IsWhiteSpace() {
			runeScanner.NextRune()
		}

		if runeScanner.IsDigit() {
			tokens = append(tokens, processInteger(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() ||
			runeScanner.IsPlus() || runeScanner.IsMultiply() {
			tokens = append(tokens, processOperand(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				fmt.Println("rune:", string(runeScanner.GetRune()))
				log.Fatalf("incorrect rune in parser")
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processInteger(rs *RunePosition) IToken {
	currentInt := make([]rune, 0)
	start := rs.GetCurrentPosition()

	for rs.IsDigit() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentInt = append(currentInt, rs.GetRune())
		rs.NextRune()
	}

	curPosition := rs.GetCurrentPosition()

	return NewIntegerToken(string(currentInt), NewFragment(start, curPosition))
}

func processOperand(rs *RunePosition) IToken {
	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		return NewOpenBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		return NewCloseBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		return NewMultiply(string(operand), NewFragment(start, curPosition))
	} else if operand == '+' {
		return NewPlus(string(operand), NewFragment(start, curPosition))
	}

	log.Fatalf("the non real error %v", NewFragment(start, rs.GetCurrentPosition()))
	return nil
}
