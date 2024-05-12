package main

import (
	"fmt"
	"log"
)

type Rule struct {
	LeftSymbol   string
	RightSymbols []string
}

// Функция для поиска множества FIRST для символа
func findFirst(symbol string, rules []Rule) map[string]bool {
	firstSet := make(map[string]bool)

	// Если символ уже был обработан, вернуть его FIRST множество
	if _, ok := firstSets[symbol]; ok {
		return firstSets[symbol]
	}

	// Если символ терминальный, добавить его в FIRST множество и вернуть
	if isTerminal(symbol) {
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
					firstOfRight := findFirst(rightSymbol, rules)
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

var terminal = map[string]struct{}{
	"+": {},
	"*": {},
	"n": {},
	"(": {},
	")": {},
}

func isTerminal(symbol string) bool {
	_, ok := terminal[symbol]
	return ok
}

var firstSets map[string]map[string]bool

var followSets map[string]map[string]bool

// Функция для нахождения множества FOLLOW для символа
func findFollow(symbol string, rules []Rule, startSymbol string) map[string]bool {
	followSet := make(map[string]bool)

	// Добавляем символ конца строки в FOLLOW множество стартового символа
	if symbol == startSymbol {
		followSet["$"] = true
	}

	// Если множество FOLLOW для данного символа уже было найдено, возвращаем его
	if _, ok := followSets[symbol]; ok {
		return followSets[symbol]
	}

	visited := make(map[string]bool)

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
					for firstSymbol := range findFirst(nextSymbol, rules) {
						if firstSymbol != "" {
							followSet[firstSymbol] = true
						}
					}

					// Если FIRST множество следующего символа содержит ε,
					// добавляем FOLLOW множество левого символа правила
					if _, ok := findFirst(nextSymbol, rules)[""]; ok {
						if !visited[rule.LeftSymbol] {
							// Добавляем FOLLOW множество левого символа правила
							for followSymbol := range findFollow(rule.LeftSymbol, rules, startSymbol) {
								followSet[followSymbol] = true
							}
						}
					}
				}

				// Если символ является последним в правиле и находится перед искомым символом
				if rightSymbol == s && i == len(rule.RightSymbols)-1 {
					// Добавляем FOLLOW множество левого символа правила
					if !visited[rule.LeftSymbol] {
						for followSymbol := range findFollow(rule.LeftSymbol, rules, startSymbol) {
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

// Таблица

func makeKey(a, b string) string {
	return fmt.Sprintf("%s - %s", a, b)
}

// Функция для получения FIRST множества для правой части правила
func getFirstSet(rightSymbols []string, first map[string]map[string]bool) (map[string]bool, bool) {
	firstSet := make(map[string]bool)

	// Если правая часть правила пустая, добавляем ε (пустую строку) в FIRST множество

	if len(rightSymbols) == 0 {
		log.Fatalf("right symbols is not zero len")
	}

	isEpsilon := false

	if len(rightSymbols[0]) == 0 {
		fmt.Println("правая часть пустая")
		isEpsilon = true
		return firstSet, isEpsilon
	}

	// Перебираем символы правой части правила
	for _, symbol := range rightSymbols {
		// Если символ - терминал, добавляем его в FIRST множество
		if isTerminal(symbol) {
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

func GetTable(rules []Rule, first, follow map[string]map[string]bool) map[string][]string {
	table := make(map[string][]string)

	// Цикл по каждому правилу грамматики
	for _, r := range rules {

		fmt.Println("current rule", r)

		// Получение множества FOLLOW для левого символа правила
		rFollow := follow[r.LeftSymbol]

		fmt.Println("rFollow", r.LeftSymbol, ":", rFollow)

		// Получение множества FIRST для правой части правила
		rFirst, isEpsilon := getFirstSet(r.RightSymbols, first)

		fmt.Println("rFirst", r.RightSymbols, ":", rFirst)

		// Добавление множества FIRST в таблицу для каждого терминала
		for term := range rFirst {
			table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
		}

		// Если множество FIRST содержит ε, добавление множества FOLLOW для соответствующего нетерминала
		if isEpsilon {
			for term := range rFollow {
				table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
			}
		}

		fmt.Println("_____________")
	}

	return table
}

//func main() {
// Пример правил переписывания изи - работает
//rules := []Rule{
//	{"E", []string{"T", "F"}},
//	{"F", []string{"+", "T", "F"}},
//	{"F", []string{""}}, // Пустое правило
//	{"T", []string{"n"}},
//	{"T", []string{"(", "E", ")"}},
//}
//
//// Пример вызова функции для символа A
//firstSets = make(map[string]map[string]bool)
//firstE := findFirst("E", rules)
//fmt.Println("FIRST(E):", firstE)
//firstF := findFirst("F", rules)
//fmt.Println("FIRST(F):", firstF)
//firstT := findFirst("T", rules)
//fmt.Println("FIRST(T):", firstT)
//
//followSets = make(map[string]map[string]bool)
//
//followE := findFollow("E", rules, "E")
//fmt.Println("FOLLOW(E):", followE)
//followF := findFollow("F", rules, "E")
//fmt.Println("FOLLOW(F):", followF)
//followT := findFollow("T", rules, "E")
//fmt.Println("FOLLOW(T):", followT)
//
//table := GetTable(rules, firstSets, followSets)
//
//fmt.Println("TABLE")
//
//for k, v := range table {
//	fmt.Println(k, "->", v)
//}

// Пример правил переписывания с лекций - проверяю
//	rules := []Rule{
//		{"E1", []string{"T1", "E2"}},
//		{"E2", []string{"+", "T1", "E2"}},
//		{"E2", []string{""}}, // Пустое правило
//		{"T1", []string{"F", "T2"}},
//		{"T2", []string{"*", "F", "T2"}},
//		{"T2", []string{""}},
//		{"F", []string{"n"}},
//		{"F", []string{"(", "E1", ")"}},
//	}
//	// Пример вызова функции для символа A
//	firstSets = make(map[string]map[string]bool)
//	firstE1 := findFirst("E1", rules)
//	fmt.Println("FIRST(E1):", firstE1)
//	firstE2 := findFirst("E2", rules)
//	fmt.Println("FIRST(E2):", firstE2)
//	firstT1 := findFirst("T1", rules)
//	fmt.Println("FIRST(T1):", firstT1)
//	firstT2 := findFirst("T2", rules)
//	fmt.Println("FIRST(T2):", firstT2)
//	firstF := findFirst("F", rules)
//	fmt.Println("FIRST(F):", firstF)
//
//	fmt.Println(firstSets)
//
//	followSets = make(map[string]map[string]bool)
//
//	followE1 := findFollow("E1", rules, "E1")
//	fmt.Println("FOLLOW(E1):", followE1)
//	followE2 := findFollow("E2", rules, "E1")
//	fmt.Println("FOLLOW(E2):", followE2)
//	followT1 := findFollow("T1", rules, "E1")
//	fmt.Println("FOLLOW(T1):", followT1)
//	followT2 := findFollow("T2", rules, "E1")
//	fmt.Println("FOLLOW(T2):", followT2)
//	followF := findFollow("F", rules, "E1")
//	fmt.Println("FOLLOW(F):", followF)
//
//	fmt.Println(followSets)
//
//	fmt.Println("TABLE")
//
//	table := GetTable(rules, firstSets, followSets)
//	for k, v := range table {
//		fmt.Println(k, "->", v)
//		fmt.Printf("len: %q\n", v)
//	}
//
//}
//
