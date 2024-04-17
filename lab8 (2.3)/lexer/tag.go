package lexer

type DomainTag int

const (
	NonTermTag DomainTag = iota + 1
	TermTag
	OperationTag
	CommentTag
	EopTag
)

var tagToString = map[DomainTag]string{
	NonTermTag:   "NON_TERMINAL",
	TermTag:      "TERMINAL",
	OperationTag: "OPERATION",
	CommentTag:   "COMMENT",
	EopTag:       "EOP",
}
