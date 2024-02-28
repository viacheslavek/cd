package lab_lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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
	patternAloneZero := "0"
	reAloneZero := regexp.MustCompile(patternAloneZero)

	patternSequenceUnits := "1+"
	reSequenceUnits := regexp.MustCompile(patternSequenceUnits)

	patternRegularStrings := `"(\\.|[^"\\])*"`
	reRegularStrings := regexp.MustCompile(patternRegularStrings)

	patternLiteralStringsWithCloseQuotes := `@"([^"]*("")?)*"(\s|$)*`
	reLiteralStringsWithCloseQuotes := regexp.MustCompile(patternLiteralStringsWithCloseQuotes)

	patternLiteralStringsWithoutCloseQuotes := `@"([^"]*("")?)*`
	reLiteralStringsWithoutCloseQuotes := regexp.MustCompile(patternLiteralStringsWithoutCloseQuotes)

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
			// TODO: работаю с ошибками
			fmt.Println("error")
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

	// TODO: не нашел закрывающую кавычку до конца - вернуть ошибку
	return nil
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
