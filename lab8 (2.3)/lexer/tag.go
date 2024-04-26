package lexer

type DomainTag int

const (
	NonTermTag DomainTag = iota + 1
	TermTag
	OpenBracketTag
	CloseBracketTag
	AxiomTag
	CommentTag
	EopTag
)

var TagToString = map[DomainTag]string{
	NonTermTag:      "NonTerminal",
	TermTag:         "Terminal",
	OpenBracketTag:  "OpenBracket",
	CloseBracketTag: "CloseBracket",
	AxiomTag:        "AxiomSign",
	CommentTag:      "Comment",
	EopTag:          "Eop",
}
