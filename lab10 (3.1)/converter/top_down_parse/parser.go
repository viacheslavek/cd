package top_down_parse

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
)

type Parser struct {
	table map[string][]string
}

func NewParser() Parser {
	return Parser{
		table: newTable(),
	}
}

func (p Parser) TopDownParse(scanner *lexer.Scanner) (*TreeNode, error) {
	type stackNode struct {
		itn *InnerTreeNode
		val string
	}
	s := NewStack[stackNode]()

	root := newTreeNode()
	root.addNode(newInnerTreeNode(""))
	s.Push(stackNode{itn: root.Root, val: Declarations})

	t := scanner.NextToken()

	for t.GetType() != lexer.EopTag {
		topNode, err := s.Pop()
		if err != nil {
			return newTreeNode(), fmt.Errorf("failed to get top node: %w", err)
		}

		if isTerminal(topNode.val) {
			topNode.itn.Children = append(topNode.itn.Children, newLeafTreeNode(t))
			t = scanner.NextToken()
		} else if neighbourhoods, ok := p.table[newTableKey(topNode.val, lexer.TagToString[t.GetType()])]; ok {
			in := newInnerTreeNode(topNode.val)
			topNode.itn.Children = append(topNode.itn.Children, in)

			for i := len(neighbourhoods) - 1; i >= 0; i-- {
				s.Push(stackNode{itn: in, val: neighbourhoods[i]})
			}
		} else {
			return newTreeNode(), fmt.Errorf("failed do parse in table with val %s and token %s",
				topNode.val, t.GetValue())
		}
	}

	return root, nil
}
