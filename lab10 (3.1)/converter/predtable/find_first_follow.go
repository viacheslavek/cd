package predtable

import (
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

// Функция для поиска множества FIRST для символа
func findFirst(symbol string, rules []semantic.Rule, terminals map[string]struct{}) map[string]bool {
	firstSet := make(map[string]bool)

	// Если символ уже был обработан, вернуть его FIRST множество
	if _, ok := firstSets[symbol]; ok {
		return firstSets[symbol]
	}

	// Если символ терминальный, добавить его в FIRST множество и вернуть
	if isTerminal(symbol, terminals) {
		firstSet[symbol] = true
		return firstSet
	}

	// Для каждого правила, где символ является левым нетерминалом
	for _, rule := range rules {
		if rule.LeftSymbol == symbol {
			// Рассмотрим правила вида A -> ε
			if len(rule.RightSymbols) == 1 && rule.RightSymbols[0] == "" {
				firstSet[""] = true
			} else {
				for _, rightSymbol := range rule.RightSymbols {
					// Для каждого символа X1, X2, ... , Xn
					firstOfRight := findFirst(rightSymbol, rules, terminals)
					// Добавим FIRST(X1) к FIRST(A), кроме ε
					for k := range firstOfRight {
						if k != "" {
							firstSet[k] = true
						}
					}
					// Если FIRST(X1) не содержит ε, завершим цикл
					if _, ok := firstOfRight[""]; !ok {
						break
					}
				}
			}
		}
	}

	// Запомним FIRST множество для будущих вызовов
	firstSets[symbol] = firstSet
	return firstSet
}

func isTerminal(symbol string, terminals map[string]struct{}) bool {
	_, ok := terminals[symbol]
	return ok
}

// Глобальная переменная для хранения FIRST множеств для каждого символа
var firstSets map[string]map[string]bool

// Глобальная переменная для хранения FOLLOW множеств для каждого символа
var followSets map[string]map[string]bool

// Функция для нахождения множества FOLLOW для символа
func findFollow(
	symbol string, rules []semantic.Rule, startSymbol string, terminals map[string]struct{},
) map[string]bool {
	followSet := make(map[string]bool)

	// Добавляем символ конца строки в FOLLOW множество стартового символа
	if symbol == startSymbol {
		followSet["$"] = true
	}

	// Если множество FOLLOW для данного символа уже было найдено, возвращаем его
	if _, ok := followSets[symbol]; ok {
		return followSets[symbol]
	}

	// Флаг для отслеживания рекурсивных вызовов функции
	visited := make(map[string]bool)

	// Функция для рекурсивного поиска множества FOLLOW для символа
	var findFollowRecursive func(string)

	findFollowRecursive = func(s string) {
		visited[s] = true

		// Для каждого правила грамматики
		for _, rule := range rules {
			for i, rightSymbol := range rule.RightSymbols {
				// Если символ является нетерминалом и находится перед искомым символом
				if rightSymbol == s && i < len(rule.RightSymbols)-1 {
					nextSymbol := rule.RightSymbols[i+1]

					// Добавляем FIRST множество следующего символа в FOLLOW множество текущего символа
					for firstSymbol := range findFirst(nextSymbol, rules, terminals) {
						if firstSymbol != "" {
							followSet[firstSymbol] = true
						}
					}

					// Если FIRST множество следующего символа содержит ε, добавляем FOLLOW множество левого символа правила
					if _, ok := findFirst(nextSymbol, rules, terminals)[""]; ok {
						if !visited[rule.LeftSymbol] {
							// Добавляем FOLLOW множество левого символа правила
							for followSymbol := range findFollow(rule.LeftSymbol, rules, startSymbol, terminals) {
								followSet[followSymbol] = true
							}
						}
					}
				}

				// Если символ является последним в правиле и находится перед искомым символом
				if rightSymbol == s && i == len(rule.RightSymbols)-1 {
					// Добавляем FOLLOW множество левого символа правила
					if !visited[rule.LeftSymbol] {
						for followSymbol := range findFollow(rule.LeftSymbol, rules, startSymbol, terminals) {
							followSet[followSymbol] = true
						}
					}
				}
			}
		}
	}

	findFollowRecursive(symbol)

	// Запоминаем множество FOLLOW для данного символа
	followSets[symbol] = followSet

	return followSet
}
