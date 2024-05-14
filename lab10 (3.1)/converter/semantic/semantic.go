package semantic

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
)

type Rules struct {
	Rule        []Rule
	Axiom       string
	Terminal    map[string]struct{}
	NonTerminal map[string]struct{}
}

type Rule struct {
	LeftSymbol   string
	RightSymbols []string
}

type Semantic struct {
	Tree *top_down_parse.TreeNode
}

func NewSemantic(tree *top_down_parse.TreeNode) *Semantic {
	return &Semantic{
		Tree: tree,
	}
}

func (s *Semantic) StartSemanticAnalysis() (Rules, error) {
	log.Println("start semantic analysis")

	allNonTerminalSymbol, terminalSymbol, ac := getTerminalAndNonTerminal(*s.Tree)

	if ac == 0 {
		return Rules{}, fmt.Errorf("zero axiom, need one")
	}
	if ac > 1 {
		return Rules{}, fmt.Errorf("axiom isn`t be better than 1, give: %d", ac)
	}

	rules := Rules{
		Rule:        make([]Rule, 0),
		Axiom:       "",
		NonTerminal: allNonTerminalSymbol,
		Terminal:    terminalSymbol,
	}

	leftNonTerminals := make(map[string]struct{})

	convertTreeToRewritingsRules(*s.Tree, &rules, &leftNonTerminals)

	if !isFirstSetInSecond(allNonTerminalSymbol, leftNonTerminals) {
		return Rules{}, fmt.Errorf("there are unreachable nonterminals %+v, %+v",
			leftNonTerminals, allNonTerminalSymbol)
	}

	convertEmptiness(&rules.Rule)

	return rules, nil
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

func convertTreeToRewritingsRules(tree top_down_parse.TreeNode, rules *Rules, leftNonTerminals *map[string]struct{}) {
	log.Println("start convert")

	root, errRoot := checkDeclarationNode(tree.Root.Children[0])
	if errRoot != nil {
		log.Printf("error in check root %+v", errRoot)
	}

	handleDeclarations(root, rules, leftNonTerminals)
}

func checkInnerNode(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, ok := node.(*top_down_parse.InnerTreeNode)

	if !ok {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node declaration %t, %s", ok, innerNode.NonTerminal)
	}

	return innerNode, nil
}

func checkDeclarationNode(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Declarations || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node declaration %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleDeclarations(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminals *map[string]struct{}) {
	if len(node.Children) == 2 {
		// Первый - RewritingRule, Второй - Declaration
		rewritingRule, errRR := checkRewritingRuleNode(node.Children[0])
		if errRR != nil {
			log.Printf("error in check handle declaration %+v", errRR)
		}
		handleRewritingRule(rewritingRule, rules, leftNonTerminals)
		declaration, errD := checkDeclarationNode(node.Children[1])
		if errD != nil {
			log.Printf("error in check handle declaration %+v", errD)
		}
		handleDeclarations(declaration, rules, leftNonTerminals)
	} else if len(node.Children) == 1 {
		// Первый - RewritingRule
		rewritingRule, errRR := checkRewritingRuleNode(node.Children[0])
		if errRR != nil {
			log.Printf("error in check handle declaration %+v", errRR)
		}
		handleRewritingRule(rewritingRule, rules, leftNonTerminals)
	} else {
		log.Println("Длина не один и не два в handleDeclarations ???")
	}
}

func checkRewritingRuleNode(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.RewritingRule || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting rule %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleRewritingRule(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminals *map[string]struct{}) {
	if len(node.Children) == 3 {
		// Первый - Axiom, Второй - NonTerminal Leaf, Третий - REWRITING
		axiom, errA := getLeafValue(node.Children[0])
		if errA != nil {
			log.Printf("error in check handle rewriting rule axiom %+v", errA)
		}
		if !isAxiom(axiom) {
			log.Printf("error in axiom checker. Axiom doesn't have axiom tag")
		}

		nonTerminal, errNT := getLeafValue(node.Children[1])
		if errNT != nil {
			log.Printf("error in check handle rewriting rule non terminal %+v", errNT)
		}

		putToLeftNonTerminalTable(nonTerminal, leftNonTerminals)

		rules.Axiom = nonTerminal.GetValue()

		rewriting, errR := checkRewriting(node.Children[2])
		if errR != nil {
			log.Printf("error in check handle rewriting rule rewriting %+v", errR)
		}
		handleRewriting(rewriting, rules, nonTerminal.GetValue())
	} else if len(node.Children) == 2 {
		nonTerminal, errNT := getLeafValue(node.Children[0])
		if errNT != nil {
			log.Printf("error in check handle rewriting rule non terminal %+v", errNT)
		}

		putToLeftNonTerminalTable(nonTerminal, leftNonTerminals)

		rewriting, errR := checkRewriting(node.Children[1])
		if errR != nil {
			log.Printf("error in check handle rewriting rule rewriting %+v", errR)
		}
		handleRewriting(rewriting, rules, nonTerminal.GetValue())
	} else {
		log.Println("Длина не два и не три в handleRewritingRule ???")
	}
}

func checkRewriting(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Rewriting || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func checkLeaf(node top_down_parse.TreeNodePrinter) (*top_down_parse.LeafTreeNode, error) {
	leafNode, ok := node.(*top_down_parse.LeafTreeNode)

	if !ok {
		return &top_down_parse.LeafTreeNode{},
			fmt.Errorf("error in cheking leaf %t", ok)
	}

	return leafNode, nil
}

func getLeafValue(node top_down_parse.TreeNodePrinter) (lexer.IToken, error) {
	leaf, err := checkLeaf(node)
	if err != nil {
		return lexer.Token{},
			fmt.Errorf("error in get leaf value %+v", err)
	}

	return leaf.Token, nil
}

func isAxiom(t lexer.IToken) bool {
	return t.GetType() == lexer.AxiomTag
}

func putToLeftNonTerminalTable(t lexer.IToken, NonTerminalTable *map[string]struct{}) {
	if t.GetType() != lexer.NonTermTag {
		log.Println("nonTerminal has nonTerminal tag", t.GetType(), t.GetValue())
	}
	(*NonTerminalTable)[t.GetValue()] = struct{}{}
}

func handleRewriting(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminal string) {
	if len(node.Children) == 4 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket, Четвертый - REWRITING_OPT

		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting close bracket %+v", errCB)
		}

		rewritingOpt, errRO := checkRewritingOpt(node.Children[3])
		if errRO != nil {
			log.Printf("error in check handle rewriting rewriting opt %+v", errRO)
		}
		handleRewritingOpt(rewritingOpt, rules, leftNonTerminal)
	} else {
		log.Println("Длина четыре в handleRewriting ???")
	}
}

func checkRewritingOpt(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.RewritingOpt || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func checkOpenBracketLeaf(node top_down_parse.TreeNodePrinter) error {
	leaf, err := checkLeaf(node)

	if leaf.Token.GetType() != lexer.OpenBracketTag || err != nil {
		return fmt.Errorf("error in cheking leaf open bracket %s, %+v", leaf.Token, err)
	}

	return nil
}

func checkCloseBracketLeaf(node top_down_parse.TreeNodePrinter) error {
	leaf, err := checkLeaf(node)

	if leaf.Token.GetType() != lexer.CloseBracketTag || err != nil {
		return fmt.Errorf("error in cheking leaf close bracket %s, %+v", leaf.Token, err)
	}

	return nil
}

func checkBody(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Body || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node body %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleRewritingOpt(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminal string) {
	if len(node.Children) == 4 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket, Четвертый - REWRITING_OPT

		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting opt body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt close bracket %+v", errCB)
		}

		rewritingOpt, errRO := checkRewritingOpt(node.Children[3])
		if errRO != nil {
			log.Printf("error in check handle rewriting opt rewriting opt %+v", errRO)
		}
		handleRewritingOpt(rewritingOpt, rules, leftNonTerminal)
	} else if len(node.Children) == 3 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket
		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting opt body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt close bracket %+v", errCB)
		}
	} else if len(node.Children) == 0 {
		// Ничего не делаем
	} else {
		log.Println("Длина не четыре и не три, и не ноль в handleRewritingOpt ???")
	}
}

func handleBody(node *top_down_parse.InnerTreeNode, currentBody *[]string) {
	if len(node.Children) == 2 {
		// Первый - Лист, Второй - BODY
		token, errT := getLeafValue(node.Children[0])
		if errT != nil {
			log.Printf("error in check handle body leaf %+v", errT)
		}
		*currentBody = append(*currentBody, token.GetValue())

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle body body %+v", errB)
		}
		handleBody(body, currentBody)

	} else if len(node.Children) == 0 {
		// Ничего не делаем
	} else {
		log.Println("Длина не два и не ноль в handleBody ???")
	}
}

func (r *Rules) putRule(leftNonTerminal string, body []string) {
	r.Rule = append(r.Rule, Rule{LeftSymbol: leftNonTerminal, RightSymbols: body})
}

func (r *Rules) Print() {
	fmt.Println("RULES:")
	fmt.Println("Terminal:", r.Terminal)
	fmt.Println("NonTerminal:", r.NonTerminal)
	fmt.Println("Axiom:", r.Axiom)
	fmt.Println("Rewriting Rules:")
	for _, rule := range r.Rule {
		fmt.Printf("%s -> %q\n", rule.LeftSymbol, rule.RightSymbols)
	}
}

func isFirstSetInSecond(first, second map[string]struct{}) bool {
	for key := range first {
		if _, ok := second[key]; !ok {
			return false
		}
	}
	return true
}

func convertEmptiness(rules *[]Rule) {
	for i := 0; i < len(*rules); i++ {
		if len((*rules)[i].RightSymbols) == 0 {
			(*rules)[i].RightSymbols = append((*rules)[i].RightSymbols, "")
		}
	}
}
