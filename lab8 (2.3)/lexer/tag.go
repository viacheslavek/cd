package lexer

type DomainTag int

// TODO: добавить комментарии
const (
	IdentTag DomainTag = iota + 1
	NonTermTag
	TermTag
	OperationTag
	ErrTag
	EopTag
)

var tagToString = map[DomainTag]string{
	IdentTag:     "IDENTIFIER",
	NonTermTag:   "NON_TERMINAL",
	TermTag:      "TERMINAL",
	OperationTag: "OPERATION",
	ErrTag:       "ERROR",
	EopTag:       "EOP",
}
