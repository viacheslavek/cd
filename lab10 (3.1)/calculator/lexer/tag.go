package lexer

type DomainTag int

const (
	IntTag DomainTag = iota + 1
	OpenBracketTag
	CloseBracketTag
	PlusTag
	MultiplyTag
	EopTag
)

var TagToString = map[DomainTag]string{
	IntTag:          "n",
	OpenBracketTag:  "(",
	CloseBracketTag: ")",
	PlusTag:         "+",
	MultiplyTag:     "*",
	EopTag:          "Eop",
}
