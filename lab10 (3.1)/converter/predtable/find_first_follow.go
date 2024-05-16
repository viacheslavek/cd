package predtable

import (
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

func findFirst(symbol string, rules []semantic.Rule, terminals map[string]struct{}) map[string]bool {
	firstSet := make(map[string]bool)

	if _, ok := firstSets[symbol]; ok {
		return firstSets[symbol]
	}

	if isTerminal(symbol, terminals) {
		firstSet[symbol] = true
		return firstSet
	}

	for _, rule := range rules {
		if rule.LeftSymbol == symbol {

			if len(rule.RightSymbols) == 1 && rule.RightSymbols[0] == "" {
				firstSet[""] = true
			} else {
				for _, rightSymbol := range rule.RightSymbols {

					firstOfRight := findFirst(rightSymbol, rules, terminals)

					for k := range firstOfRight {
						if k != "" {
							firstSet[k] = true
						}
					}

					if _, ok := firstOfRight[""]; !ok {
						break
					}
				}
			}
		}
	}

	firstSets[symbol] = firstSet
	return firstSet
}

func isTerminal(symbol string, terminals map[string]struct{}) bool {
	_, ok := terminals[symbol]
	return ok
}

var firstSets map[string]map[string]bool

var followSets map[string]map[string]bool

func findFollow(
	symbol string, rules []semantic.Rule, startSymbol string, terminals map[string]struct{},
) map[string]bool {
	followSet := make(map[string]bool)

	if symbol == startSymbol {
		followSet["$"] = true
	}

	if _, ok := followSets[symbol]; ok {
		return followSets[symbol]
	}

	visited := make(map[string]bool)

	var findFollowRecursive func(string)

	findFollowRecursive = func(s string) {
		visited[s] = true

		for _, rule := range rules {
			for i, rightSymbol := range rule.RightSymbols {
				if rightSymbol == s && i < len(rule.RightSymbols)-1 {
					nextSymbol := rule.RightSymbols[i+1]

					for firstSymbol := range findFirst(nextSymbol, rules, terminals) {
						if firstSymbol != "" {
							followSet[firstSymbol] = true
						}
					}

					if _, ok := findFirst(nextSymbol, rules, terminals)[""]; ok {
						if !visited[rule.LeftSymbol] {
							for followSymbol := range findFollow(
								rule.LeftSymbol, rules, startSymbol, terminals) {
								followSet[followSymbol] = true
							}
						}
					}
				}

				if rightSymbol == s && i == len(rule.RightSymbols)-1 {
					if !visited[rule.LeftSymbol] {
						for followSymbol := range findFollow(
							rule.LeftSymbol, rules, startSymbol, terminals) {
							followSet[followSymbol] = true
						}
					}
				}
			}
		}
	}

	findFollowRecursive(symbol)

	followSets[symbol] = followSet

	return followSet
}
