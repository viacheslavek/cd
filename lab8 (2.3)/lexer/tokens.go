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

type Message struct {
	isComment bool
	text      string
}

func NewMessage(isComment bool, text string) Message {
	return Message{
		isComment: isComment,
		text:      text,
	}
}

type CommentToken struct {
	Token
}

func NewComment(text string, fragment Fragment) CommentToken {
	return CommentToken{
		Token: NewToken(CommentTag, text, fragment),
	}
}

func (ct CommentToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ct.Type], ct.Coordinate, ct.Value)
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type NonTermToken struct {
	Token
}

func NewNonTerminal(value string, fragment Fragment) NonTermToken {
	return NonTermToken{
		Token: NewToken(NonTermTag, value, fragment),
	}
}

func (ntt NonTermToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ntt.Type], ntt.Coordinate, ntt.Value)
}

type TerminalToken struct {
	Token
}

func NewTerminal(value string, fragment Fragment) TerminalToken {
	return TerminalToken{
		Token: NewToken(TermTag, value, fragment),
	}
}

func (tt TerminalToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[tt.Type], tt.Coordinate, tt.Value)
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

type AxiomToken struct {
	Token
}

func NewAxiom(value string, fragment Fragment) AxiomToken {
	return AxiomToken{
		Token: NewToken(AxiomTag, value, fragment),
	}
}

func (st AxiomToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}
