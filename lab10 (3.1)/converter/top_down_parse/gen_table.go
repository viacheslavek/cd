package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"BODY - CloseBracket":          {},
		"BODY - NonTerminal":           {"NonTerminal", "BODY"},
		"BODY - Terminal":              {"Terminal", "BODY"},
		"DECLARATIONS - $":             {},
		"DECLARATIONS - AxiomSign":     {"REWRITING_RULE", "DECLARATIONS"},
		"DECLARATIONS - NonTerminal":   {"REWRITING_RULE", "DECLARATIONS"},
		"REWRITING - OpenBracket":      {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_OPT - $":            {},
		"REWRITING_OPT - AxiomSign":    {},
		"REWRITING_OPT - NonTerminal":  {},
		"REWRITING_OPT - OpenBracket":  {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_RULE - AxiomSign":   {"AxiomSign", "NonTerminal", "REWRITING"},
		"REWRITING_RULE - NonTerminal": {"NonTerminal", "REWRITING"},
	}
}

func newGenAxiom() string {
	return "DECLARATIONS"
}
