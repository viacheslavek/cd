package lexer

import (
	"fmt"
	"sort"
)

type Compiler struct {
	tokens   []Token
	messages map[Fragment]Message
}

func NewCompiler() *Compiler {
	return &Compiler{
		messages: make(map[Fragment]Message),
	}
}

func (c *Compiler) AddMessage(ct CommentToken) {
	c.messages[ct.Coordinate] = NewMessage(true, ct.Value)
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
		fmt.Printf("Type: Comment | i: %d | position: %v | text: %s\n",
			i, position, c.messages[position].text)
	}
}
