package top_down_parse

import (
	"fmt"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
)

type TreeNode struct {
	root *innerTreeNode
}

func newTreeNode() *TreeNode {
	return &TreeNode{}
}

func (tn *TreeNode) Print() {
	tn.root.printNode(0)
}

func (tn *TreeNode) addNode(node *innerTreeNode) {
	tn.root = node
}

type treeNodePrinter interface {
	printNode(offset int)
}

type innerTreeNode struct {
	nonTerminal string
	children    []treeNodePrinter
}

func newInnerTreeNode(nonTerminal string) *innerTreeNode {
	return &innerTreeNode{nonTerminal: nonTerminal, children: make([]treeNodePrinter, 0)}
}

func (itn innerTreeNode) printNode(offset int) {
	fmt.Printf(strings.Repeat("\t", offset) + fmt.Sprintf("Inner node: %s\n", itn.nonTerminal))

	for _, child := range itn.children {
		child.printNode(offset + 1)
	}
}

type leafTreeNode struct {
	token lexer.IToken
}

func newLeafTreeNode(token lexer.IToken) *leafTreeNode {
	return &leafTreeNode{token: token}
}

func (ltn *leafTreeNode) printNode(offset int) {
	if ltn.token.GetType() == lexer.TermTag || ltn.token.GetType() == lexer.NonTermTag {
		fmt.Printf(strings.Repeat("\t", offset) +
			fmt.Sprintf("Leaf: %s - %s\n", lexer.TagToString[ltn.token.GetType()], ltn.token.GetValue()))
	} else {
		fmt.Printf(strings.Repeat("\t", offset) +
			fmt.Sprintf("Leaf: %s\n", lexer.TagToString[ltn.token.GetType()]))
	}
}
