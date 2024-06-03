package lexer

type DomainTag int

const (
	IdentifierTag DomainTag = iota + 1
	IntTag
	OpenBracketTag
	CloseBracketTag
	AxiomTag
	EopTag
)

var TagToString = map[DomainTag]string{
	IdentifierTag:   "Identifier",
	IntTag:          "Integer",
	OpenBracketTag:  "OpenBracket",
	CloseBracketTag: "CloseBracket",
	AxiomTag:        "AxiomSign",
	EopTag:          "Eop",
}
