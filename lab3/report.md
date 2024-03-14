% Лабораторная работа № 1.2. «Лексический анализатор на основе регулярных выражений»
% 28 февраля 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является приобретение навыка разработки простейших лексических анализаторов,
работающих на основе поиска в тексте по образцу, заданному регулярным выражением.

# Индивидуальный вариант
Числовые литералы: знак «0» либо последовательности знаков «1».
Строковые литералы: регулярные строки — ограничены двойными кавычками,
могут содержать escape-последовательности `\"`, `\t`, `\n`,
не пересекают границы строк текста; буквальные строки — начинаются на «@"»,
заканчиваются на двойную кавычку, пересекают границы строк текста,
для включения двойной кавычки она удваивается.


# Реализация

**lexer.go**
```go
package lab_lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Lexer struct {
	tokens   []TokenAndError
	position int
}

func NewLexer(filepath string) *Lexer {
	return &Lexer{
		ParseFile(filepath),
		0,
	}
}

func (l *Lexer) NextToken() TokenAndError {
	if l.position < len(l.tokens) {
		token := l.tokens[l.position]
		l.position++
		return token
	}
	return SyntaxError{Message: EOF}

}

func ParseFile(filepath string) []TokenAndError {
	fmt.Println("parseFile")

	file, errO := os.Open(filepath)
	if errO != nil {
		log.Fatalf("can't open file in parser %e", errO)
	}
	defer func() {
		_ = file.Close()
	}()

	tae := make([]TokenAndError, 0)
	currentLine := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parseLine(&tae, scanner, line, &currentLine)

		currentLine++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	return tae
}

func parseLine(taes *[]TokenAndError, scanner *bufio.Scanner, line string, currentLine *int) {
	currentColumn := 1
	trimGapStep := 0

	// Компилируем регулярные выражения один раз
	reAloneZero := regexp.MustCompile(`(?:^|[^0])0(?:$|[^0\S])`)

	reSequenceUnits := regexp.MustCompile(`1+(?:$|[^1\S])`)

	reRegularStrings := regexp.MustCompile(`"(\\.|[^"\\])*"`)

	reLiteralStringsWithCloseQuotes := regexp.MustCompile(`@"([^"]*("")?)*"`)

	reLiteralStringsWithoutCloseQuotes := regexp.MustCompile(`@"([^"]*("")?)*`)

	for len(line) > 0 {
		trimGapStep, line = trimLeadingGapChars(line)

		if len(line) == 0 {
			return
		}

		currentColumn += trimGapStep
		var matchIndex []int
		var tore TokenAndError
		if matchIndex = reAloneZero.FindStringIndex(line); len(matchIndex) > 0 && matchIndex[0] == 0 {
			tore = processToken(&line, matchIndex, currentLine, &currentColumn, "ZERO")

		} else if matchIndex = reSequenceUnits.FindStringIndex(line); len(matchIndex) > 0 && matchIndex[0] == 0 {
			tore = processToken(&line, matchIndex, currentLine, &currentColumn, "ONE_SEQ")

		} else if matchIndex = reRegularStrings.FindStringIndex(line); len(matchIndex) > 0 && matchIndex[0] == 0 {
			tore = processToken(&line, matchIndex, currentLine, &currentColumn, "REGULAR_STR")

		} else if matchIndex = reLiteralStringsWithCloseQuotes.FindStringIndex(line); len(matchIndex) > 0 &&
			matchIndex[0] == 0 {
			tore = processToken(&line, matchIndex, currentLine, &currentColumn, "LITERAL_STR")

		} else if matchIndex = reLiteralStringsWithoutCloseQuotes.FindStringIndex(line); len(matchIndex) > 0 &&
			matchIndex[0] == 0 {
			tore = processLiteralStrInManyLine(scanner, &line, currentLine, &currentColumn)

		} else {
			tore = processSyntaxError(&line, currentLine, &currentColumn,
				reAloneZero, reSequenceUnits, reRegularStrings,
				reLiteralStringsWithCloseQuotes, reLiteralStringsWithoutCloseQuotes)
		}
		*taes = append(*taes, tore)
	}
}

func trimLeadingGapChars(str string) (int, string) {
	trimChars := " \t\n\r"
	trimmed := strings.TrimLeft(str, trimChars)
	return len(str) - len(trimmed), trimmed
}

func processToken(line *string,
	matchIndex []int, currentLine, currentColumn *int, typeToken string) TokenAndError {
	tore := NewToken(typeToken, (*line)[:matchIndex[1]], *currentLine, *currentColumn)
	*currentColumn += matchIndex[1] - matchIndex[0]
	*line = (*line)[matchIndex[1]:]

	return tore
}

func processLiteralStrInManyLine(
	scanner *bufio.Scanner, line *string, currentLine, currentColumn *int,
) TokenAndError {
	tokenValue := *line
	prevLine := *currentLine
	prevColumn := *currentColumn

	for scanner.Scan() {
		newLine := scanner.Text()
		// Если нашел нечетное число кавычек, то на этом буквальная строка заканчивается,
		// line обрезаю до конца
		// буквальной строки и задаю новые параметры currentLine, currentColumn
		// если в строке такого не нашлось, то иду дальше
		index := 0
		if ok := findOddNumberDoubleQuotes(newLine, &index); ok {
			tokenValue += "\n" + newLine[:index]
			*currentLine++
			*currentColumn += index
			*line = newLine[index:]
			return NewToken("LITERAL_STR", tokenValue, prevLine, prevColumn)
		} else {
			tokenValue += "\n" + newLine
			*currentLine++
		}
	}

	// буквальная строка так и не закончилась
	*line = ""
	return NewError(": literal string not end", prevLine, prevColumn)
}

func findOddNumberDoubleQuotes(line string, index *int) (ok bool) {
	for left := 0; left < len(line); left++ {
		if line[left] == '"' {
			right := left
			for ; right < len(line) && line[right] == '"'; right++ {
			}
			if (right-left)%2 != 0 {
				*index += right + 1
				return true
			} else if right >= len(line) {
				return false
			} else {
				*index += right + 1
				return findOddNumberDoubleQuotes(line[right+1:], index)
			}
		}
	}

	return false
}

func processSyntaxError(
	line *string, currentLine, currentColumn *int,
	zero, one, regular, literal, literalFull *regexp.Regexp) TokenAndError {

	tore := NewError("", *currentLine, *currentColumn)

	// Нахожу минимум между ними
	zeroIndex := zero.FindStringIndex(*line)

	oneIndex := one.FindStringIndex(*line)
	regularStrIndex := regular.FindStringIndex(*line)
	literalStrIndex := literal.FindStringIndex(*line)
	literalStrFullIndex := literalFull.FindStringIndex(*line)

	numsPred := [][]int{zeroIndex, oneIndex, regularStrIndex, literalStrIndex, literalStrFullIndex}

	nums := make([]int, 0)

	for _, n := range numsPred {
		if len(n) != 0 {
			nums = append(nums, n[0])
		}
	}

	sort.Ints(nums)
	if len(nums) == 0 {
		// в строке не осталось литералов
		*line = ""
	} else {
		// обрезаем до первого возможно валидного литерала
		*line = (*line)[nums[0]:]
		*currentColumn += nums[0]
	}

	return tore
}
```

**token.go**
```go
package lab_lexer

import "fmt"

const (
	EOF = "EOF"
)

type TokenAndError interface {
	IsToken() bool
	IsError() bool

	String() string
	CurrentType() string
}

type Token struct {
	Type     string
	Value    string
	Position Position
}

func NewToken(t, v string, line, column int) Token {
	return Token{
		Type:     t,
		Value:    v,
		Position: Position{line, column},
	}
}

func (t Token) IsToken() bool { return true }
func (t Token) IsError() bool { return false }

func (t Token) String() string {
	return fmt.Sprintf("%s (%d, %d): %s", t.Type, t.Position.Line, t.Position.Column, t.Value)
}

func (t Token) CurrentType() string {
	return t.Type
}

type SyntaxError struct {
	Message  string
	Position Position
}

func NewError(m string, line, column int) SyntaxError {
	return SyntaxError{
		Message:  m,
		Position: Position{line, column},
	}
}

func (e SyntaxError) IsToken() bool { return false }
func (e SyntaxError) IsError() bool { return true }

func (e SyntaxError) String() string {
	return fmt.Sprintf("syntax error (%d, %d)%s", e.Position.Line, e.Position.Column, e.Message)
}

func (e SyntaxError) CurrentType() string {
	return e.Message
}

type Position struct {
	Line   int
	Column int
}

```

**main.go**
```go
package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab3/lab_lexer"
)

const filepath = "test_files/mixed.txt"

func main() {

	lexer := lab_lexer.NewLexer(filepath)

	token := lexer.NextToken()
	for token.IsToken() || token.IsError() && token.CurrentType() != lab_lexer.EOF {
		fmt.Println(token)
		token = lexer.NextToken()
	}

	fmt.Println("finish")
}
```

# Тестирование

Для тестирования было написано несколько тестирующих файлов,
покрывающих большую часть тестовых случаев

Входные данные

**mixed.txt**
```
0 111 1 "aaa"
00 "aa\t\n\"l"
1 "cccc 1
@"aaaaa
aaaa ""  bbbb
ccccc" "ddd" 1 1111

```

**number_correct.txt**
```
0 0 0 11 1111
0 0 1 1111 11 1 11
0 111 1 111111 0
0 1
```

**number_error.txt**
```
0 1112 00 1 1 1 00
1 1 0 1111 00
```

**str_literal_correct.txt**
```
@"aaaaa
aaaa  bbbb
ccccc" @"ddd"
@" aa"" aa"""" a"
@"aaaa
aaa
" @"aaaaaa"
@"aaaa""aaaaaa" @"aa bb cc"
@"aaa   bbb    ccc"
@"aaaa" @"bbbbb"
```

**str_literal_error.txt**
```
@"aaaa" aaaaa"
@"aaaa" @"bbbbb
```

**str_regular_correct.txt**
```
"abcd" "abcd abcd"
"abcd\nabcd" "abcd\tabcd" "abcd\"abcd"
"abcd \nabcd" "abcd \tabcd" "abcd \"abcd"
"abcd \n abcd" "abcd \t abcd" "abcd \" abcd"
"   abcd   " "   abcd  abcd   "
"abcd"
```

**str_regular_error.txt**
```
"abcd" "abcd abcd
"abcd\nabcd" "abcd \\tabcd" "abcd" abcd"
"   abcd "    abcd  abcd   "
abcd"
```

Вывод на `stdout` для *mixed.txt*

```
parseFile
ZERO (1, 1): 0 
ONE_SEQ (1, 3): 111 
ONE_SEQ (1, 7): 1 
REGULAR_STR (1, 9): "aaa"
syntax error (2, 1)
REGULAR_STR (2, 4): "aa\t\n\"l"
ONE_SEQ (3, 1): 1 
syntax error (3, 3)
ONE_SEQ (3, 9): 1
LITERAL_STR (4, 1): @"aaaaa
aaaa ""  bbbb
ccccc" 
REGULAR_STR (6, 8): "ddd"
ONE_SEQ (6, 14): 1 
ONE_SEQ (6, 16): 1111
finish
```

# Вывод
В ходе данной лабораторной работы был получен навык разработки простейшего лексического анализатора, 
работающего на основе поиска в тексте по образцу, заданному регулярным выражением. При этом я лучше 
стал понимать регулярные выражения и их практическое применение. Ведь большая часть работы для меня 
была в корректном составлении этих регулярных выражений
