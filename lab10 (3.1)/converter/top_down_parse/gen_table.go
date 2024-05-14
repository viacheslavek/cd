package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"DECLARATIONS - AxiomSign":   {"REWRITING_RULE", "DECLARATIONS"},
		"DECLARATIONS - NonTerminal": {"REWRITING_RULE", "DECLARATIONS"},
		"DECLARATIONS - Eof":         {},

		"REWRITING_RULE - AxiomSign":   {"AxiomSign", "NonTerminal", "REWRITING"},
		"REWRITING_RULE - NonTerminal": {"NonTerminal", "REWRITING"},

		"REWRITING - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},

		"REWRITING_OPT - AxiomSign":   {},
		"REWRITING_OPT - NonTerminal": {},
		"REWRITING_OPT - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_OPT - Eof":         {},

		"BODY - NonTerminal":  {"NonTerminal", "BODY"},
		"BODY - Terminal":     {"Terminal", "BODY"},
		"BODY - CloseBracket": {},
	}
}

func newGenAxiom() string {
	return "DECLARATIONS"
}
