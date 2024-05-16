package predtable

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

func makeKey(a, b string) string {
	return fmt.Sprintf("%s - %s", a, b)
}

func getFirstSet(
	rightSymbols []string, first map[string]map[string]bool, terminals map[string]struct{},
) (map[string]bool, bool) {
	firstSet := make(map[string]bool)

	if len(rightSymbols) == 0 {
		log.Fatalf("right symbols is not zero len")
	}

	isEpsilon := false

	if len(rightSymbols[0]) == 0 {
		isEpsilon = true
		return firstSet, isEpsilon
	}

	for _, symbol := range rightSymbols {
		if isTerminal(symbol, terminals) {
			firstSet[symbol] = true
			break
		}

		for term := range first[symbol] {
			if term != "" {
				firstSet[term] = true
			} else {
				isEpsilon = true
			}
		}

		if !first[symbol][""] {
			break
		}
	}

	return firstSet, isEpsilon
}

func getTable(
	rules []semantic.Rule, first, follow map[string]map[string]bool, terminals map[string]struct{},
) (map[string][]string, error) {
	table := make(map[string][]string)

	for _, r := range rules {
		rFollow := follow[r.LeftSymbol]
		rFirst, isEpsilon := getFirstSet(r.RightSymbols, first, terminals)

		for term := range rFirst {
			if _, ok := table[makeKey(r.LeftSymbol, term)]; ok {
				return table,
					fmt.Errorf("rules is not LL1: two rules in one cell %s", makeKey(r.LeftSymbol, term))
			}
			table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
		}

		if isEpsilon {
			for term := range rFollow {
				if _, ok := table[makeKey(r.LeftSymbol, term)]; ok {
					return table,
						fmt.Errorf("rules is not LL1: two rules in one cell %s", makeKey(r.LeftSymbol, term))
				}
				table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
			}
		}
	}

	return table, nil
}

func GenTable(rules semantic.Rules) (map[string][]string, error) {
	firstSets = make(map[string]map[string]bool)
	for n := range rules.NonTerminal {
		_ = findFirst(n, rules.Rule, rules.Terminal)
	}

	followSets = make(map[string]map[string]bool)
	for n := range rules.NonTerminal {
		_ = findFollow(n, rules.Rule, rules.Axiom, rules.Terminal)
	}

	return getTable(rules.Rule, firstSets, followSets, rules.Terminal)
}

func PrintGenTable(genTable map[string][]string) {
	fmt.Println("Gen Table:")
	for k, v := range genTable {
		fmt.Println(k, "->", v)
	}
}
