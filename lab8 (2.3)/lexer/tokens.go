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
	return fmt.Sprintf("%s %s: %s", tagToString[ct.Type], ct.Coordinate, ct.Value)
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
	return fmt.Sprintf("%s %s: %s", tagToString[ntt.Type], ntt.Coordinate, ntt.Value)
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
	return fmt.Sprintf("%s %s: %s", tagToString[tt.Type], tt.Coordinate, tt.Value)
}

type OperationToken struct {
	Token
}

func NewOperation(value string, fragment Fragment) OperationToken {
	return OperationToken{
		Token: NewToken(OperationTag, value, fragment),
	}
}

func (st OperationToken) String() string {
	return fmt.Sprintf("%s %s: %s", tagToString[st.Type], st.Coordinate, st.Value)
}
