package lexer

import (
	"fmt"
)

type IToken interface {
	GetType() DomainTag
}

type Token struct {
	Type       DomainTag
	Value      string
	Coordinate Fragment
}

type Fragment struct {
	start  TokenPosition
	finish TokenPosition
}

func NewFragment(start, finish TokenPosition) Fragment {
	return Fragment{
		start:  start,
		finish: finish,
	}
}

func (f Fragment) String() string {
	return fmt.Sprintf("%s-%s", f.start.String(), f.finish.String())
}

type TokenPosition struct {
	line   int
	column int
}

func NewTokenPosition(line, column int) TokenPosition {
	return TokenPosition{
		line:   line,
		column: column,
	}
}

func (p *TokenPosition) String() string {
	return fmt.Sprintf("(%d,%d)", p.line, p.column)
}

func NewToken(tag DomainTag, value string, coordinate Fragment) Token {
	return Token{
		Type:       tag,
		Value:      value,
		Coordinate: coordinate,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s: %s", tagToString[t.Type], t.Coordinate, t.Value)
}

func (t Token) GetType() DomainTag {
	return t.Type
}

type Message struct {
	isError bool
	text    string
}

func NewMessage(isError bool, text string) Message {
	return Message{
		isError: isError,
		text:    text,
	}
}

type ErrorToken struct {
	Token
}

func NewError(text string, fragment Fragment) ErrorToken {
	return ErrorToken{
		Token: NewToken(ErrTag, text, fragment),
	}
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type IdentToken struct {
	Token
	attr int
}

func NewIdent(value string, fragment Fragment) IdentToken {
	return IdentToken{
		Token: NewToken(IdentTag, value, fragment),
	}
}

func (it *IdentToken) SetAttr(attr int) {
	it.attr = attr
}

func (it *IdentToken) String() string {
	return fmt.Sprintf("%s %s: %d", tagToString[it.Type], it.Coordinate, it.attr)
}

type TerminalToken struct {
	Token
	attr string
}

func NewTerminal(value string, fragment Fragment) TerminalToken {
	return TerminalToken{
		Token: NewToken(TermTag, value, fragment),
	}
}

func (st *TerminalToken) String() string {
	return fmt.Sprintf("%s %s: %s", tagToString[st.Type], st.Coordinate, st.attr)
}

type OperationToken struct {
	Token
	operation int
}

func NewOperation(value string, fragment Fragment) IdentToken {
	return IdentToken{
		Token: NewToken(OperationTag, value, fragment),
	}
}
