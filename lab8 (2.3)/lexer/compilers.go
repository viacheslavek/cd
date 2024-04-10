package lexer

import (
	"fmt"
	"sort"
)

type Compiler struct {
	tokens           []Token
	messages         map[Fragment]Message
	identifiersTable map[string]int
	identifiers      []IdentToken
}

func NewCompiler() *Compiler {
	return &Compiler{
		messages:         make(map[Fragment]Message),
		identifiersTable: make(map[string]int),
		identifiers:      make([]IdentToken, 0),
	}
}

func (c *Compiler) AddMessage(et ErrorToken) {
	c.messages[et.Coordinate] = NewMessage(true, et.Value)
}

func (c *Compiler) PrintMessages() {
	sortedMessagesFragments := make([]Fragment, len(c.messages))
	index := 0
	for key := range c.messages {
		sortedMessagesFragments[index] = key
		index++
	}

	sort.Slice(sortedMessagesFragments, func(i, j int) bool {
		return sortedMessagesFragments[i].start.line < sortedMessagesFragments[j].start.line &&
			sortedMessagesFragments[i].start.column < sortedMessagesFragments[j].start.column
	})
	fmt.Println("_____MESSAGES_____")
	for i, position := range sortedMessagesFragments {
		fmt.Printf("Type: Error | i: %d | position: %v | text: %s\n",
			i, position, c.messages[position].text)
	}
}

func (c *Compiler) GetIdentifier(identifier string) IdentToken {
	return c.identifiers[c.identifiersTable[identifier]]
}

func (c *Compiler) AddIdentifier(identifier IdentToken) int {
	val, ok := c.identifiersTable[identifier.Value]
	if !ok {
		iPosition := len(c.identifiers)
		c.identifiers = append(c.identifiers, identifier)
		c.identifiersTable[identifier.Value] = iPosition
		return iPosition
	}
	return val
}

func (c *Compiler) PrintIdentifiers() {
	fmt.Println("____Identifiers____")
	for i, id := range c.identifiers {
		fmt.Println(tagToString[id.Type], id.Coordinate, i, "--", id.Value)
	}
}
