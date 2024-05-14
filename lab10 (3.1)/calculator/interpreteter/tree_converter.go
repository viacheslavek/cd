package interpreteter

import (
	"fmt"
	"log"
	"strconv"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/top_down_parse"
)

func Solve(tree *top_down_parse.TreeNode) int {
	log.Println("start convert")

	return evaluateExpression(tree.Root.Children[0])
}

func evaluateExpression(node top_down_parse.TreeNodePrinter) int {
	switch n := node.(type) {
	case *top_down_parse.InnerTreeNode:
		if n.NonTerminal == "E" {
			if len(n.Children) == 2 {
				fmt.Println("в E 2")
				left := evaluateExpression(n.Children[0])
				right := evaluateExpression(n.Children[1])
				return left + right
			} else if len(n.Children) == 1 {
				fmt.Println("в E 1")
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в E равна не 1 и не 2:", n, len(n.Children))
			}
		} else if n.NonTerminal == "E'" {
			if len(n.Children) == 3 {
				fmt.Println("в E' 3")
				left := evaluateExpression(n.Children[1])
				right := evaluateExpression(n.Children[2])
				return left + right
			} else if len(n.Children) == 2 {
				fmt.Println("в E' 2")
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 0 {
				fmt.Println("в E' 0")
				return 0 // epsilon
			} else {
				log.Fatal("Длина в E' равна не 3 и не 0:", n, len(n.Children))
			}
		} else if n.NonTerminal == "T" {
			if len(n.Children) == 2 {
				fmt.Println("в T 2")
				left := evaluateExpression(n.Children[0])
				right := evaluateExpression(n.Children[1])
				return left * right
			} else if len(n.Children) == 1 {
				fmt.Println("в T 1")
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в T равна не 1 и не 2:", n, len(n.Children))
			}
		} else if n.NonTerminal == "T'" {
			if len(n.Children) == 3 {
				fmt.Println("в T' 3")
				left := evaluateExpression(n.Children[1])
				right := evaluateExpression(n.Children[2])
				return left * right
			} else if len(n.Children) == 2 {
				fmt.Println("в T' 2")
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 0 {
				fmt.Println("в T' 0")
				return 1 // epsilon
			} else {
				log.Fatal("Длина в T' равна не 3 и не 2 и не 0:", n, len(n.Children))
			}
		} else if n.NonTerminal == "F" {
			if len(n.Children) == 3 {
				fmt.Println("в F 3")
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 1 {
				fmt.Println("в F 1")
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в F равна не 3 и не 1:", n, len(n.Children))
			}
		} else {
			log.Fatal("Неизвестный нетерминал", n, len(n.Children))
		}
	case *top_down_parse.LeafTreeNode:
		fmt.Println("in leaf", n, *n)
		value, _ := strconv.Atoi(n.Token.GetValue())
		return value
	}
	return 0
}
