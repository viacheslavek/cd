package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Scanner struct {
	tokens   []IToken
	compiler *Compiler
	position int
}

func NewScanner(filepath string) *Scanner {
	return &Scanner{
		tokens:   ParseFile(filepath),
		compiler: NewCompiler(),
		position: 0,
	}
}

func (s *Scanner) NextToken() IToken {

	// TODO: и вот тут меняю

	token := s.tokens[s.position]
	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == IdentTag {
		idPos := s.compiler.AddIdentifier(token.(IdentToken))
		newToken := token.(IdentToken)
		newToken.SetAttr(idPos)
		return newToken
	}
	if token.GetType() == TermTag ||
		token.GetType() == OperationTag {
		return token
	}
	if token.GetType() == ErrTag {
		newToken := token.(ErrorToken)
		s.compiler.AddMessage(newToken)
	}

	return token
}

func (s *Scanner) GetCompiler() *Compiler {
	return s.compiler
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

		// TODO: вот тут дополняю новыми

		if runeScanner.IsQuote() {
			tokens = append(tokens, processTerminal(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() || runeScanner.IsStarBracket() {
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

	fmt.Println("tokens:")

	return tokens
}

func processTerminal(rs *RunePosition) IToken {
	currentString := make([]rune, 0)

	start := rs.GetCurrentPosition()

	for !rs.IsQuote() {
		if rs.GetRune() == -1 {
			return NewError("the line didn't end", NewFragment(start, rs.GetCurrentPosition()))
		}
		// INFO: Внутри все что угодно, кроме еще одной кавычки
		currentString = append(currentString, rs.GetRune())
		rs.NextRune()
	}

	curPosition := rs.GetCurrentPosition()
	rs.NextRune()

	return NewTerminal(string(currentString), NewFragment(start, curPosition))
}

func processOperand(rs *RunePosition) IToken {

	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		NewOperation(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		NewOperation(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		NewOperation(string(operand), NewFragment(start, curPosition))
	}

	return NewError("the non real error", NewFragment(start, rs.GetCurrentPosition()))
}

func processIdentifier(rs *RunePosition) IToken {
	currentIdent := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()

	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			// Строим слово
			currentIdent = append(currentIdent, rs.GetRune())
		} else {
			// ошибка в Identifier, не буквы или цифры
			curPosition := rs.GetCurrentPosition()
			for !rs.IsWhiteSpace() {
				curPosition = rs.GetCurrentPosition()
				rs.NextRune()
			}
			return NewError("symbol is not a latin letter or digit", NewFragment(start, curPosition))
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}
	return NewIdent(string(currentIdent), NewFragment(start, curPositionToken))
}
