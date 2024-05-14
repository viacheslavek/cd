package lexer

import (
	"fmt"
)

type IToken interface {
	GetType() DomainTag
	GetValue() string
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
	return fmt.Sprintf("%s %s: %s", TagToString[t.Type], t.Coordinate, t.Value)
}

func (t Token) GetType() DomainTag {
	return t.Type
}

func (t Token) GetValue() string {
	return t.Value
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type IntegerToken struct {
	Token
}

func NewIntegerToken(value string, fragment Fragment) IntegerToken {
	return IntegerToken{
		Token: NewToken(IntTag, value, fragment),
	}
}

func (st IntegerToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type OpenBracketToken struct {
	Token
}

func NewOpenBracket(value string, fragment Fragment) OpenBracketToken {
	return OpenBracketToken{
		Token: NewToken(OpenBracketTag, value, fragment),
	}
}

func (st OpenBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type CloseBracketToken struct {
	Token
}

func NewCloseBracket(value string, fragment Fragment) CloseBracketToken {
	return CloseBracketToken{
		Token: NewToken(CloseBracketTag, value, fragment),
	}
}

func (st CloseBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type PlusToken struct {
	Token
}

func NewPlus(value string, fragment Fragment) PlusToken {
	return PlusToken{
		Token: NewToken(PlusTag, value, fragment),
	}
}

func (st PlusToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type MultiplyToken struct {
	Token
}

func NewMultiply(value string, fragment Fragment) MultiplyToken {
	return MultiplyToken{
		Token: NewToken(MultiplyTag, value, fragment),
	}
}

func (st MultiplyToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}
