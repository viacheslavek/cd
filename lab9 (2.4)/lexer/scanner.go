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
			tokens = append(tokens, processInt(runeScanner))
		} else if runeScanner.IsSpecialSymbol() {
			tokens = append(tokens, processSpecialSymbol(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processIdentifier(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				fmt.Println("rune:", string(runeScanner.GetRune()))
				log.Fatalf("incorrect rune in parser %+v", runeScanner.GetCurrentPosition())
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processInt(rs *RunePosition) IToken {
	currentInt := make([]rune, 0)
	start := rs.GetCurrentPosition()

	for rs.IsDigit() {
		currentInt = append(currentInt, rs.GetRune())
		rs.NextRune()
	}

	curPosition := rs.GetCurrentPosition()

	return NewInteger(string(currentInt), NewFragment(start, curPosition))
}

func processSpecialSymbol(rs *RunePosition) IToken {
	start := rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()
	curPosition := rs.GetCurrentPosition()

	return NewSpecSymbol(string(operand), NewFragment(start, curPosition))
}

func processIdentifier(rs *RunePosition) IToken {
	currentIdentifier := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()
	currentIdentifier = append(currentIdentifier, rs.GetRune())
	rs.NextRune()

	for !rs.IsWhiteSpace() && !rs.IsSpecialSymbol() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else if rs.IsUnderlining() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else {
			log.Fatalf("error in process identifier")
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}

	if IsKeyword(string(currentIdentifier)) {
		return NewKeyword(string(currentIdentifier), NewFragment(start, curPositionToken))
	}

	return NewIdentifier(string(currentIdentifier), NewFragment(start, curPositionToken))
}
