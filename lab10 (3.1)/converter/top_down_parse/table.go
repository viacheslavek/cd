package top_down_parse

import (
	"fmt"
	"slices"
)

const (
	Declarations  = "DECLARATIONS"
	RewritingRule = "REWRITING_RULE"
	Rewriting     = "REWRITING"
	RewritingOpt  = "REWRITING_OPT"
	Body          = "BODY"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{"DECLARATIONS", "REWRITING_RULE", "REWRITING", "REWRITING_OPT", "BODY"}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}
