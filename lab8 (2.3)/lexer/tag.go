package lexer

type DomainTag int

const (
	IdentTag DomainTag = iota + 1
	StrTag
	ErrTag
	EopTag
)

var tagToString = map[DomainTag]string{
	IdentTag: "IDENTIFIER",
	StrTag:   "STRING",
	ErrTag:   "ERROR",
	EopTag:   "EOP",
}
