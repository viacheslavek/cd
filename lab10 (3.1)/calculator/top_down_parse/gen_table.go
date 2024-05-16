package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"E - n":  {"T", "E'"},
		"E - )":  {"T", "E'"},
		"E' - +": {"+", "T", "E'"},
		"E' - )": {},
		"E' - $": {},
		"F - n":  {"n"},
		"F - (":  {"(", "E", ")"},
		"T - n":  {"F", "T'"},
		"T - (":  {"F", "T'"},
		"T' - *": {"*", "F", "T'"},
		"T' - +": {},
		"T' - )": {},
		"T' - $": {},
	}
}

func newGenAxiom() string {
	return "E"
}
