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
	currentLine := 0

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
	currentColumn := 0
	trimGapStep := 0

	// Компилируем регулярные выражения один раз
	reAloneZero := regexp.MustCompile(`(?:^|[^0])0(?:$|[^0\S])`)

	reSequenceUnits := regexp.MustCompile(`1+`)

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

func processToken(line *string, matchIndex []int, currentLine, currentColumn *int, typeToken string) TokenAndError {
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
		// Если нашел нечетное число кавычек, то на этом буквальная строка заканчивается, line обрезаю до конца
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
