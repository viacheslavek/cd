package lab_lexer

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
	// TODO: тут встраиваю все для компилятора
	if s.position < len(s.tokens) {
		token := s.tokens[s.position]
		s.position++
		return token
	}
	// TODO: возвращаю специальный токен EOF
	return Token{}

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

		switch runeScanner.GetRune() {
		case '"':
			runeScanner.NextRune()
			tokens = append(tokens, processString(runeScanner))
		default:
			if runeScanner.IsLetter() {
				tokens = append(tokens, processIdentifier(runeScanner))
			} else {
				if runeScanner.GetRune() == -1 {
					tokens = append(tokens, NewEOP())
				} else {
					fmt.Println("это не планировалось")
					tokens = append(tokens, processStartError(runeScanner))
					return tokens
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	return tokens
}

func processString(rs *RunePosition) IToken {
	currentString := make([]rune, 0)
	currentString = append(currentString, '"')

	start := rs.GetCurrentPosition()

	for !rs.IsQuote() {
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
			if !rs.IsApostrophe() && !rs.IsQuote() && !rs.IsHashtag() {
				for !rs.IsQuote() && !rs.IsLineTranslation() {
					rs.NextRune()
				}
				rs.NextRune()
				curPosition := rs.GetCurrentPosition()
				curPosition.column--

				return NewError("the character code has bad end", NewFragment(start, rs.GetCurrentPosition()))
			}
		} else {
			for !rs.IsQuote() && !rs.IsLineTranslation() {
				rs.NextRune()
			}
			rs.NextRune()
			curPosition := rs.GetCurrentPosition()
			curPosition.column--
			return NewError("symbol is not a section or character code", NewFragment(start, curPosition))
		}
	}

	currentString = append(currentString, rs.GetRune())
	rs.NextRune()

	curPosition := rs.GetCurrentPosition()
	curPosition.column--
	return NewString(string(currentString),
		NewFragment(start, curPosition), string(currentString))
}

func processIdentifier(rs *RunePosition) IToken {
	currentIdent := make([]rune, 0)

	start := rs.GetCurrentPosition()

	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if rs.IsLetter() || rs.IsDigit() {
			// Строим слово
			currentIdent = append(currentIdent, rs.GetRune())
		} else {
			// ошибка в Identifier, не буквы или цифры
			for !rs.IsWhiteSpace() {
				rs.NextRune()
			}
			curPosition := rs.GetCurrentPosition()
			curPosition.column--
			return NewError("symbol is not a letter or digit", NewFragment(start, curPosition))
		}
		rs.NextRune()
	}
	curPosition := rs.GetCurrentPosition()
	curPosition.column--
	return NewIdent(string(currentIdent), NewFragment(start, curPosition), string(currentIdent))
}

// TODO: делаю обработку стартовой ошибки - если не с буквы и не с кавычки началось
func processStartError(rs *RunePosition) IToken {
	fmt.Println("Разбираем ошибку, что у меня нет таких токенов, с которых могло бы начинаться")
	return Token{}
}
