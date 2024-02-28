package lab_lexer

import "fmt"

const (
	EOF = "EOF"
)

type TokenAndError interface {
	IsToken() bool
	IsError() bool

	String() string
	CurrentType() string
}

type Token struct {
	Type     string
	Value    string
	Position Position
}

func NewToken(t, v string, line, column int) Token {
	return Token{
		Type:     t,
		Value:    v,
		Position: Position{line, column},
	}
}

func (t Token) IsToken() bool { return true }
func (t Token) IsError() bool { return false }

func (t Token) String() string {
	return fmt.Sprintf("%s (%d, %d): %s", t.Type, t.Position.Line, t.Position.Column, t.Value)
}

func (t Token) CurrentType() string {
	return t.Type
}

type SyntaxError struct {
	Message  string
	Position Position
}

func NewError(m string, line, column int) SyntaxError {
	return SyntaxError{
		Message:  m,
		Position: Position{line, column},
	}
}

func (e SyntaxError) IsToken() bool { return false }
func (e SyntaxError) IsError() bool { return true }

func (e SyntaxError) String() string {
	return fmt.Sprintf("syntax error (%d, %d): %s", e.Position.Line, e.Position.Column, e.Message)
}

func (e SyntaxError) CurrentType() string {
	return e.Message
}

type Position struct {
	Line   int
	Column int
}
