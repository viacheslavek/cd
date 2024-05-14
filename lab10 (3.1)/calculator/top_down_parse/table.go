package top_down_parse

import (
	"fmt"
	"slices"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{"E", "E'", "T", "T'", "F"}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}
