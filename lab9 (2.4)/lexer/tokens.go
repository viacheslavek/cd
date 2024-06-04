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

type IdentifierToken struct {
	Token
}

func NewIdentifier(value string, fragment Fragment) IdentifierToken {
	return IdentifierToken{
		Token: NewToken(IdentifierTag, value, fragment),
	}
}

func (ntt IdentifierToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ntt.Type], ntt.Coordinate, ntt.Value)
}

type IntegerToken struct {
	Token
}

func NewInteger(value string, fragment Fragment) IntegerToken {
	return IntegerToken{
		Token: NewToken(IntTag, value, fragment),
	}
}

func (tt IntegerToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[tt.Type], tt.Coordinate, tt.Value)
}

type SpecSymbolToken struct {
	Token
}

func NewSpecSymbol(value string, fragment Fragment) SpecSymbolToken {
	return SpecSymbolToken{
		Token: NewToken(SpecSymbolTag, value, fragment),
	}
}

func (st SpecSymbolToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type KeywordToken struct {
	Token
}

func NewKeyword(value string, fragment Fragment) KeywordToken {
	return KeywordToken{
		Token: NewToken(KeywordTag, value, fragment),
	}
}

func (st KeywordToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}
