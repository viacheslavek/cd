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

	// TODO: обрабатывать знаки - specSymbol - {'-', '+', '*', '/', '(', ')', '[', ']' , '{', '}', ',', ';'}

	for runeScanner.GetRune() != -1 {
		for runeScanner.IsWhiteSpace() {
			runeScanner.NextRune()
		}

		if runeScanner.IsDigit() {
			tokens = append(tokens, processInt(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() || runeScanner.IsStar() {
			tokens = append(tokens, processOperand(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processIdentifier(runeScanner))
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

func processOperand(rs *RunePosition) IToken {
	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		return NewOpenBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		return NewCloseBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		return NewAxiom(string(operand), NewFragment(start, curPosition))
	}

	log.Fatalf("the non real error %v", NewFragment(start, rs.GetCurrentPosition()))
	return nil
}

// TODO: в ident смотреть на ключевые слова в ProcessIdent
// keywords = {enum, struct, union, sizeof, char, char, short, int, long, float, double}
// А так же сохранять идентификаторы в сет
func processIdentifier(rs *RunePosition) IToken {
	currentIdentifier := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()
	currentIdentifier = append(currentIdentifier, rs.GetRune())
	rs.NextRune()

	for !rs.IsWhiteSpace() && !rs.IsBrackets() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else if rs.IsOneQuote() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else {
			log.Fatalf("error in process identifier")
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}

	return NewIdentifier(string(currentIdentifier), NewFragment(start, curPositionToken))
}
