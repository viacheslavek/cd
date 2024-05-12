package semantic

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
	"log"
)

type Semantic struct {
	Tree *top_down_parse.TreeNode
}

func NewSemantic(tree *top_down_parse.TreeNode) *Semantic {
	return &Semantic{
		Tree: tree,
	}
}

func (s *Semantic) StartSemanticAnalysis() error {

	log.Println("start semantic analysis")

	nonTerminalSymbol, terminalSymbol, ac := getTerminalAndNonTerminal(*s.Tree)

	fmt.Println("NonTerminal:", nonTerminalSymbol)
	fmt.Println("Terminal:", terminalSymbol)
	fmt.Println("Axioms:", ac)

	rules := convertTreeToRewritingsRules(*s.Tree)

	fmt.Println(rules)

	return nil
}

func getTerminalAndNonTerminal(tree top_down_parse.TreeNode) (
	nonTerminal map[string]struct{},
	terminal map[string]struct{},
	axiomCount int,
) {
	nonTerminal = make(map[string]struct{})
	terminal = make(map[string]struct{})

	traverseTree(tree.Root, &nonTerminal, &terminal, &axiomCount)

	return nonTerminal, terminal, axiomCount
}

func traverseTree(
	node top_down_parse.TreeNodePrinter, nonTerminals, terminals *map[string]struct{}, axiomCount *int,
) {

	switch n := node.(type) {
	case *top_down_parse.InnerTreeNode:
		// Тут я нахожусь во внутреннем узле - это рабочее пространство дерева
		for _, child := range n.Children {
			traverseTree(child, nonTerminals, terminals, axiomCount)
		}
	case *top_down_parse.LeafTreeNode:
		// Либо терминал, либо нетерминал, либо служебные символы
		if n.Token.GetType() == lexer.TermTag {
			(*terminals)[n.Token.GetValue()] = struct{}{}
		} else if n.Token.GetType() == lexer.NonTermTag {
			(*nonTerminals)[n.Token.GetValue()] = struct{}{}
		} else if n.Token.GetType() == lexer.AxiomTag {
			*axiomCount++
		}
	default:
		log.Println("default?", n)
	}

}

type Rule struct {
	LeftSymbol   string
	RightSymbols []string
}

func convertTreeToRewritingsRules(tree top_down_parse.TreeNode) []Rule {
	rules := make([]Rule, 0)

	return rules
}
