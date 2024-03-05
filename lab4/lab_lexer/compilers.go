package lab_lexer

type Compiler struct {
	tokens           []Token
	messages         map[TokenPosition]Message
	identifiersTable map[string]int
	identifiers      []Token
}

func NewCompiler() *Compiler {
	return &Compiler{
		messages:         make(map[TokenPosition]Message),
		identifiersTable: make(map[string]int),
		identifiers:      make([]Token, 0),
	}
}

func (c *Compiler) GetMessages() map[TokenPosition]Message {
	return c.messages
}

func (c *Compiler) GetIdentifier(identifier string) Token {
	return c.identifiers[c.identifiersTable[identifier]]
}

func (c *Compiler) AddIdentifier(identifier Token) {
	if _, ok := c.identifiersTable[identifier.Value]; !ok {
		iPosition := len(c.identifiers)
		c.identifiers = append(c.identifiers, identifier)
		c.identifiersTable[identifier.Value] = iPosition
	}
}
