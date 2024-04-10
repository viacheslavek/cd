package lexer

import (
	"bufio"
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
	if token.GetType() == IdentTag {
		idPos := s.compiler.AddIdentifier(token.(IdentToken))
		newToken := token.(IdentToken)
		newToken.SetAttr(idPos)
		return newToken
	}
	if token.GetType() == StrTag {
		newToken := token.(StringToken)
		newToken.SetText(newToken.Value)
		return newToken
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

		if runeScanner.IsApostrophe() || runeScanner.IsHashtag() {
			tokens = append(tokens, processString(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processIdentifier(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				tokens = append(tokens, processStartError(runeScanner))
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processString(rs *RunePosition) IToken {
	currentString := make([]rune, 0)

	start := rs.GetCurrentPosition()

	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return NewError("the line didn't end", NewFragment(start, rs.GetCurrentPosition()))
		}
		if rs.IsApostrophe() {
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
			for !rs.IsApostrophe() {
				if rs.IsLineTranslation() {
					return NewError("the section didn't end", NewFragment(start, rs.GetCurrentPosition()))
				}
				currentString = append(currentString, rs.GetRune())
				rs.NextRune()
			}
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
		} else if rs.IsHashtag() {
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
			for rs.IsDigit() {
				currentString = append(currentString, rs.GetRune())
				rs.NextRune()
			}
			if !rs.IsWhiteSpace() && !rs.IsHashtag() && !rs.IsApostrophe() {
				for !rs.IsWhiteSpace() {
					rs.NextRune()
				}
				rs.NextRune()
				curPosition := rs.GetCurrentPosition()
				curPosition.column--

				return NewError("the character code has bad end", NewFragment(start, rs.GetCurrentPosition()))
			}
		} else {
			for !rs.IsApostrophe() && !rs.IsLineTranslation() {
				rs.NextRune()
			}
			curPosition := rs.GetCurrentPosition()
			rs.NextRune()
			return NewError("symbol is not a section or character code", NewFragment(start, curPosition))
		}
	}

	curPosition := rs.GetCurrentPosition()
	rs.NextRune()

	return NewString(string(currentString),
		NewFragment(start, curPosition))
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

func processStartError(rs *RunePosition) IToken {

	startPosition := rs.GetCurrentPosition()
	curPosition := rs.GetCurrentPosition()
	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		curPosition = rs.GetCurrentPosition()
		rs.NextRune()
	}
	return NewError("start is not a quote or letter", NewFragment(startPosition, curPosition))
}
