package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"E - n": {"T", "E'"},
		"E - )": {"T", "E'"},

		"T - n": {"F", "T'"},
		"T - (": {"F", "T'"},

		"T' - *":   {"*", "F", "T'"},
		"T' - +":   {},
		"T' - )":   {},
		"T' - Eof": {},

		"E' - +":   {"+", "T", "E'"},
		"E' - )":   {},
		"E' - Eof": {},

		"F - n": {"n"},
		"F - (": {"(", "E", ")"},
	}
}

func newGenAxiom() string {
	return "E"
}
