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
	token := s.tokens[s.position]
	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == TermTag || token.GetType() == OperationTag || token.GetType() == NonTermTag {
		return token
	}
	if token.GetType() == CommentTag {
		newToken := token.(CommentToken)
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

		if runeScanner.IsQuote() {
			tokens = append(tokens, processTerminal(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() || runeScanner.IsStar() {
			tokens = append(tokens, processOperand(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processNonTerminal(runeScanner))
		} else if runeScanner.IsOpenSlash() {
			tokens = append(tokens, processComment(runeScanner))
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

func processTerminal(rs *RunePosition) IToken {
	currentString := make([]rune, 0)
	currentString = append(currentString, rs.GetRune())
	start := rs.GetCurrentPosition()
	rs.NextRune()

	for !rs.IsQuote() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentString = append(currentString, rs.GetRune())
		rs.NextRune()
	}

	currentString = append(currentString, rs.GetRune())
	curPosition := rs.GetCurrentPosition()
	rs.NextRune()

	return NewTerminal(string(currentString), NewFragment(start, curPosition))
}

func processOperand(rs *RunePosition) IToken {
	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		return NewOperation(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		return NewOperation(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		return NewOperation(string(operand), NewFragment(start, curPosition))
	}

	log.Fatalf("the non real error %v", NewFragment(start, rs.GetCurrentPosition()))
	return nil
}

func processNonTerminal(rs *RunePosition) IToken {
	currentNonTerminal := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()
	currentNonTerminal = append(currentNonTerminal, rs.GetRune())
	rs.NextRune()

	for !rs.IsWhiteSpace() && !rs.IsBrackets() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			currentNonTerminal = append(currentNonTerminal, rs.GetRune())
		} else if rs.IsOneQuote() {
			currentNonTerminal = append(currentNonTerminal, rs.GetRune())
		} else {
			log.Fatalf("error in process nonterminal")
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}

	return NewNonTerminal(string(currentNonTerminal), NewFragment(start, curPositionToken))
}

func processComment(rs *RunePosition) IToken {

	currentText := make([]rune, 0)

	start := rs.GetCurrentPosition()
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	if !rs.IsStar() {
		log.Fatalf("missing quote in comment")
	}
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	for !rs.IsStar() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentText = append(currentText, rs.GetRune())
		rs.NextRune()
	}

	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	if !rs.IsOpenSlash() {
		log.Fatalf("missing slash in comment")
	}
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	curPositionToken := rs.GetCurrentPosition()

	return NewComment(string(currentText), NewFragment(start, curPositionToken))
}
