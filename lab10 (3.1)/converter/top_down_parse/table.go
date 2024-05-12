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
	nonTerminal   = "NonTerminal"
	terminal      = "Terminal"
	axiomSign     = "AxiomSign"
	OpenBracket   = "OpenBracket"
	CloseBracket  = "CloseBracket"
	eof           = "Eof"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{Declarations, RewritingRule, Rewriting, RewritingOpt, Body}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}

func newTable() map[string][]string {
	return map[string][]string{
		newTableKey(Declarations, axiomSign):   {RewritingRule, Declarations},
		newTableKey(Declarations, nonTerminal): {RewritingRule, Declarations},
		newTableKey(Declarations, eof):         {},

		newTableKey(RewritingRule, axiomSign):   {axiomSign, nonTerminal, Rewriting},
		newTableKey(RewritingRule, nonTerminal): {nonTerminal, Rewriting},

		newTableKey(Rewriting, OpenBracket): {OpenBracket, Body, CloseBracket, RewritingOpt},

		newTableKey(RewritingOpt, axiomSign):   {},
		newTableKey(RewritingOpt, nonTerminal): {},
		newTableKey(RewritingOpt, OpenBracket): {OpenBracket, Body, CloseBracket, RewritingOpt},
		newTableKey(RewritingOpt, eof):         {},

		newTableKey(Body, nonTerminal):  {nonTerminal, Body},
		newTableKey(Body, terminal):     {terminal, Body},
		newTableKey(Body, CloseBracket): {},
	}
}
