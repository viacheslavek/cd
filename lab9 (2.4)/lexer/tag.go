package lexer

type DomainTag int

const (
	IdentifierTag DomainTag = iota + 1
	IntTag
	KeywordTag
	SpecSymbolTag
	EopTag
)

var TagToString = map[DomainTag]string{
	IdentifierTag: "Identifier",
	IntTag:        "Integer",
	KeywordTag:    "Keyword",
	SpecSymbolTag: "SpecialSymbol",
	EopTag:        "Eop",
}
