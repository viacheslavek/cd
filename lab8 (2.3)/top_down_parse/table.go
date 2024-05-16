package top_down_parse

import (
	"fmt"
	"slices"
)

const (
	declarations  = "DECLARATIONS"
	rewritingRule = "REWRITING_RULE"
	rewriting     = "REWRITING"
	rewritingOpt  = "REWRITING_OPT"
	body          = "BODY"
	nonTerminal   = "NonTerminal"
	terminal      = "Terminal"
	axiomSign     = "AxiomSign"
	openBracket   = "OpenBracket"
	closeBracket  = "CloseBracket"
	eof           = "Eof"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{declarations, rewritingRule, rewriting, rewritingOpt, body}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}

func newTable() map[string][]string {
	return map[string][]string{
		newTableKey(declarations, axiomSign):   {rewritingRule, declarations},
		newTableKey(declarations, nonTerminal): {rewritingRule, declarations},
		newTableKey(declarations, eof):         {},

		newTableKey(rewritingRule, axiomSign):   {axiomSign, nonTerminal, rewriting},
		newTableKey(rewritingRule, nonTerminal): {nonTerminal, rewriting},

		newTableKey(rewriting, openBracket): {openBracket, body, closeBracket, rewritingOpt},

		newTableKey(rewritingOpt, axiomSign):   {},
		newTableKey(rewritingOpt, nonTerminal): {},
		newTableKey(rewritingOpt, openBracket): {openBracket, body, closeBracket, rewritingOpt},
		newTableKey(rewritingOpt, eof):         {},

		newTableKey(body, nonTerminal):  {nonTerminal, body},
		newTableKey(body, terminal):     {terminal, body},
		newTableKey(body, closeBracket): {},
	}
}
