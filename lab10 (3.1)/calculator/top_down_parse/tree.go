package top_down_parse

import (
	"fmt"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
)

type TreeNode struct {
	Root *InnerTreeNode
}

func newTreeNode() *TreeNode {
	return &TreeNode{}
}

func (tn *TreeNode) Print() {
	tn.Root.printNode(0)
}

func (tn *TreeNode) addNode(node *InnerTreeNode) {
	tn.Root = node
}

type TreeNodePrinter interface {
	printNode(offset int)
}

type InnerTreeNode struct {
	NonTerminal string
	Children    []TreeNodePrinter
}

func newInnerTreeNode(nonTerminal string) *InnerTreeNode {
	return &InnerTreeNode{NonTerminal: nonTerminal, Children: make([]TreeNodePrinter, 0)}
}

func (itn InnerTreeNode) printNode(offset int) {
	fmt.Printf(strings.Repeat("..", offset) + fmt.Sprintf("Inner node: %s\n", itn.NonTerminal))

	for _, child := range itn.Children {
		child.printNode(offset + 1)
	}
}

type LeafTreeNode struct {
	Token lexer.IToken
}

func newLeafTreeNode(token lexer.IToken) *LeafTreeNode {
	return &LeafTreeNode{Token: token}
}

func (ltn LeafTreeNode) printNode(offset int) {
	if ltn.Token.GetType() == lexer.IntTag {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s - %s\n", lexer.TagToString[ltn.Token.GetType()], ltn.Token.GetValue()))
	} else {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s\n", lexer.TagToString[ltn.Token.GetType()]))
	}
}
