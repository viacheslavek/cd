package lab_lexer

type DomainTag int

const (
	IdentTag DomainTag = iota
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
