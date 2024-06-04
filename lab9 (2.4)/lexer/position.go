package lexer

import (
	"bufio"
	"log"
	"unicode"
)

type RunePosition struct {
	line        int
	column      int
	currentLine []rune
	scanner     *bufio.Scanner
	EOF         bool
}

func NewRunePosition(scanner *bufio.Scanner) *RunePosition {
	scanner.Scan()
	return &RunePosition{
		line:        1,
		column:      1,
		currentLine: []rune(scanner.Text() + "\n\a"),
		scanner:     scanner,
		EOF:         false,
	}
}

func (p *RunePosition) NextRune() {
	if p.EOF {
		log.Println("EOF")
		return
	}
	p.column++
	if p.currentLine[p.column-1] == '\a' {
		p.column = 1
		p.line++
		if !p.scanner.Scan() {
			p.EOF = true
		} else {
			p.currentLine = []rune(p.scanner.Text() + "\n\a")
		}
	}
}

func (p *RunePosition) GetRune() rune {
	if p.EOF {
		return -1
	}
	return p.currentLine[p.column-1]
}

func (p *RunePosition) GetCurrentPosition() TokenPosition {
	return NewTokenPosition(p.line, p.column)
}

func (p *RunePosition) IsWhiteSpace() bool {
	return unicode.IsSpace(p.GetRune())
}

func (p *RunePosition) IsLineTranslation() bool {
	return p.GetRune() == '\n'
}

func (p *RunePosition) IsLetter() bool {
	return unicode.IsLetter(p.GetRune())
}

func (p *RunePosition) IsLatinLetter() bool {
	r := unicode.ToLower(p.GetRune())
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func (p *RunePosition) IsUnderlining() bool {
	return p.GetRune() == '_'
}

func (p *RunePosition) IsDigit() bool {
	return unicode.IsDigit(p.GetRune())
}

func (p *RunePosition) IsOpenRoundBracket() bool {
	return p.GetRune() == '('
}

func (p *RunePosition) IsCloseRoundBracket() bool {
	return p.GetRune() == ')'
}

func (p *RunePosition) IsOpenSquareBracket() bool {
	return p.GetRune() == '['
}

func (p *RunePosition) IsCloseSquareBracket() bool {
	return p.GetRune() == ']'
}

func (p *RunePosition) IsOpenCurlyBracket() bool {
	return p.GetRune() == '{'
}

func (p *RunePosition) IsCloseCurlyBracket() bool {
	return p.GetRune() == '}'
}

func (p *RunePosition) IsComma() bool {
	return p.GetRune() == ','
}

func (p *RunePosition) IsSemicolon() bool {
	return p.GetRune() == ';'
}

func (p *RunePosition) IsStar() bool {
	return p.GetRune() == '*'
}

func (p *RunePosition) IsPlus() bool {
	return p.GetRune() == '+'
}

func (p *RunePosition) IsMinus() bool {
	return p.GetRune() == '-'
}

func (p *RunePosition) IsSlash() bool {
	return p.GetRune() == '/'
}

func (p *RunePosition) IsEqual() bool {
	return p.GetRune() == '='
}

func (p *RunePosition) IsSpecialSymbol() bool {
	return p.IsOpenRoundBracket() || p.IsCloseRoundBracket() ||
		p.IsOpenSquareBracket() || p.IsCloseSquareBracket() ||
		p.IsOpenCurlyBracket() || p.IsCloseCurlyBracket() ||
		p.IsMinus() || p.IsPlus() || p.IsSlash() || p.IsStar() ||
		p.IsComma() || p.IsSemicolon() || p.IsEqual()
}

var keywords = map[string]struct{}{
	"enum": {}, "struct": {}, "union": {},
	"sizeof": {},
	"char":   {}, "short": {}, "int": {}, "long": {}, "float": {}, "double": {},
}

func IsKeyword(ident string) bool {
	_, ok := keywords[ident]
	return ok
}
