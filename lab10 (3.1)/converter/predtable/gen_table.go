package predtable

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

func makeKey(a, b string) string {
	return fmt.Sprintf("%s - %s", a, b)
}

// Функция для получения FIRST множества для правой части правила
func getFirstSet(
	rightSymbols []string, first map[string]map[string]bool, terminals map[string]struct{},
) (map[string]bool, bool) {
	firstSet := make(map[string]bool)

	// Если правая часть правила пустая, добавляем ε (пустую строку) в FIRST множество

	if len(rightSymbols) == 0 {
		log.Fatalf("right symbols is not zero len")
	}

	isEpsilon := false

	if len(rightSymbols[0]) == 0 {
		isEpsilon = true
		return firstSet, isEpsilon
	}

	// Перебираем символы правой части правила
	for _, symbol := range rightSymbols {
		// Если символ - терминал, добавляем его в FIRST множество
		if isTerminal(symbol, terminals) {
			firstSet[symbol] = true
			break
		}

		// Если символ - нетерминал, добавляем все его FIRST множество
		for term := range first[symbol] {
			if term != "" {
				firstSet[term] = true
			} else {
				isEpsilon = true
			}
		}

		// Если в FIRST множестве символа нет ε, прекращаем обработку
		if !first[symbol][""] {
			break
		}
	}

	return firstSet, isEpsilon
}

// Функция для построения таблицы предсказывающего анализа

func getTable(
	rules []semantic.Rule, first, follow map[string]map[string]bool, terminals map[string]struct{},
) (map[string][]string, error) {
	table := make(map[string][]string)

	// Цикл по каждому правилу грамматики
	for _, r := range rules {

		//fmt.Println("current rule", r)

		// Получение множества FOLLOW для левого символа правила
		rFollow := follow[r.LeftSymbol]

		//fmt.Println("rFollow", r.LeftSymbol, ":", rFollow)

		// Получение множества FIRST для правой части правила
		rFirst, isEpsilon := getFirstSet(r.RightSymbols, first, terminals)

		//fmt.Println("rFirst", r.RightSymbols, ":", rFirst)

		// Добавление множества FIRST в таблицу для каждого терминала

		for term := range rFirst {
			if _, ok := table[makeKey(r.LeftSymbol, term)]; ok {
				return table,
					fmt.Errorf("rules is not LL1: two rules in one cell %s", makeKey(r.LeftSymbol, term))
			}
			table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
		}

		// Если множество FIRST содержит ε, добавление множества FOLLOW для соответствующего нетерминала
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
		//tmpFirst := findFirst(n, rules.Rule, rules.Terminal)
		//fmt.Printf("first(%s): %+v\n", n, tmpFirst)
	}

	followSets = make(map[string]map[string]bool)
	for n := range rules.NonTerminal {
		_ = findFollow(n, rules.Rule, rules.Axiom, rules.Terminal)
		//tmpFollow := findFollow(n, rules.Rule, rules.Axiom, rules.Terminal)
		//fmt.Printf("follow(%s): %+v\n", n, tmpFollow)
	}

	return getTable(rules.Rule, firstSets, followSets, rules.Terminal)
}

func PrintGenTable(genTable map[string][]string) {
	fmt.Println("Gen Table:")
	for k, v := range genTable {
		fmt.Println(k, "->", v)
		//fmt.Printf("len: %q\n", v)
	}
}
