% Лабораторная работа № 1.3 «Объектно-ориентированный лексический анализатор»
% 6 марта 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является приобретение навыка реализации лексического анализатора на
объектно-ориентированном языке без применения каких-либо средств автоматизации решения задачи лексического анализа.

# Индивидуальный вариант
Идентификаторы: последовательности латинских букв и цифр, начинающиеся с буквы.

Строковые константы — последовательности строковых секций, записанных слитно. 
Строковые секции: либо последовательность символов, ограниченных апострофами,
апостроф внутри строки описывается как два апострофа подряд, не пересекают границы строк текста,
либо знак «#», за которым следует десятичная константа (код символа).

Пример строковой константы: «'hello'#10#13'world'» (эта строковая константа состоит из 4 строковых секций,
однако является единым токеном).


# Реализация

Файл `compilers.go`
```go
package lab_lexer

import (
	"fmt"
	"sort"
)

type Compiler struct {
	tokens           []Token
	messages         map[Fragment]Message
	identifiersTable map[string]int
	identifiers      []IdentToken
}

func NewCompiler() *Compiler {
	return &Compiler{
		messages:         make(map[Fragment]Message),
		identifiersTable: make(map[string]int),
		identifiers:      make([]IdentToken, 0),
	}
}

func (c *Compiler) AddMessage(et ErrorToken) {
	c.messages[et.Coordinate] = NewMessage(true, et.Value)
}

func (c *Compiler) PrintMessages() {
	sortedMessagesFragments := make([]Fragment, len(c.messages))
	index := 0
	for key := range c.messages {
		sortedMessagesFragments[index] = key
		index++
	}

	sort.Slice(sortedMessagesFragments, func(i, j int) bool {
		return sortedMessagesFragments[i].start.line < sortedMessagesFragments[j].start.line &&
			sortedMessagesFragments[i].start.column < sortedMessagesFragments[j].start.column
	})
	fmt.Println("_____MESSAGES_____")
	for i, position := range sortedMessagesFragments {
		fmt.Printf("Type: Error | i: %d | position: %v | text: %s\n",
			i, position, c.messages[position].text)
	}
}

func (c *Compiler) GetIdentifier(identifier string) IdentToken {
	return c.identifiers[c.identifiersTable[identifier]]
}

func (c *Compiler) AddIdentifier(identifier IdentToken) int {
	val, ok := c.identifiersTable[identifier.Value]
	if !ok {
		iPosition := len(c.identifiers)
		c.identifiers = append(c.identifiers, identifier)
		c.identifiersTable[identifier.Value] = iPosition
		return iPosition
	}
	return val
}

func (c *Compiler) PrintIdentifiers() {
	fmt.Println("____Identifiers____")
	for i, id := range c.identifiers {
		fmt.Println(tagToString[id.Type], id.Coordinate, i, "--", id.Value)
	}
}
```

Файл `position.go`
```go
package lab_lexer

import (
	"bufio"
	"log"
	"unicode"
)

type RunePosition struct {
	line        int
	column      int
	currentLine []rune
	scanner     *bufio.Scanner
	EOF         bool
}

func NewRunePosition(scanner *bufio.Scanner) *RunePosition {
	scanner.Scan()
	return &RunePosition{
		line:        1,
		column:      1,
		currentLine: []rune(scanner.Text() + "\n\a"),
		scanner:     scanner,
		EOF:         false,
	}
}

func (p *RunePosition) NextRune() {
	if p.EOF {
		log.Println("EOF")
		return
	}
	p.column++
	if p.currentLine[p.column-1] == '\a' {
		p.column = 1
		p.line++
		if !p.scanner.Scan() {
			p.EOF = true
		} else {
			p.currentLine = []rune(p.scanner.Text() + "\n\a")
		}
	}
}

func (p *RunePosition) GetRune() rune {
	if p.EOF {
		return -1
	}
	return p.currentLine[p.column-1]
}

func (p *RunePosition) GetCurrentPosition() TokenPosition {
	return NewTokenPosition(p.line, p.column)
}

func (p *RunePosition) IsWhiteSpace() bool {
	return unicode.IsSpace(p.GetRune())
}

func (p *RunePosition) IsLetter() bool {
	return unicode.IsLetter(p.GetRune())
}

func (p *RunePosition) IsLatinLetter() bool {
	r := unicode.ToLower(p.GetRune())
	return r >= 'a' && r <= 'z'
}

func (p *RunePosition) IsDigit() bool {
	return unicode.IsDigit(p.GetRune())
}

func (p *RunePosition) IsQuote() bool {
	return p.GetRune() == '"'
}

func (p *RunePosition) IsApostrophe() bool {
	return p.GetRune() == '\''
}

func (p *RunePosition) IsHashtag() bool {
	return p.GetRune() == '#'
}

func (p *RunePosition) IsLineTranslation() bool {
	return p.GetRune() == '\n'
}

```

Файл `scanner.go`
```go
package lab_lexer

import (
	"bufio"
	"log"
	"os"
)

type Scanner struct {
	tokens   []IToken
	compiler *Compiler
	position int
}

func NewScanner(filepath string) *Scanner {
	return &Scanner{
		tokens:   ParseFile(filepath),
		compiler: NewCompiler(),
		position: 0,
	}
}

func (s *Scanner) NextToken() IToken {

	token := s.tokens[s.position]
	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == IdentTag {
		idPos := s.compiler.AddIdentifier(token.(IdentToken))
		newToken := token.(IdentToken)
		newToken.SetAttr(idPos)
		return newToken
	}
	if token.GetType() == StrTag {
		newToken := token.(StringToken)
		newToken.SetText(newToken.Value)
		return newToken
	}
	if token.GetType() == ErrTag {
		newToken := token.(ErrorToken)
		s.compiler.AddMessage(newToken)
	}

	return token
}

func (s *Scanner) GetCompiler() *Compiler {
	return s.compiler
}

func ParseFile(filepath string) []IToken {
	file, errO := os.Open(filepath)
	if errO != nil {
		log.Fatalf("can't open file in parser %e", errO)
	}
	defer func() {
		_ = file.Close()
	}()

	tokens := make([]IToken, 0)

	scanner := bufio.NewScanner(file)
	runeScanner := NewRunePosition(scanner)

	for runeScanner.GetRune() != -1 {
		for runeScanner.IsWhiteSpace() {
			runeScanner.NextRune()
		}

		if runeScanner.IsApostrophe() || runeScanner.IsHashtag() {
			tokens = append(tokens, processString(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processIdentifier(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				tokens = append(tokens, processStartError(runeScanner))
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processString(rs *RunePosition) IToken {
	currentString := make([]rune, 0)

	start := rs.GetCurrentPosition()

	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return NewError("the line didn't end", NewFragment(start, rs.GetCurrentPosition()))
		}
		if rs.IsApostrophe() {
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
			for !rs.IsApostrophe() {
				if rs.IsLineTranslation() {
					return NewError("the section didn't end", NewFragment(start, rs.GetCurrentPosition()))
				}
				currentString = append(currentString, rs.GetRune())
				rs.NextRune()
			}
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
		} else if rs.IsHashtag() {
			currentString = append(currentString, rs.GetRune())
			rs.NextRune()
			for rs.IsDigit() {
				currentString = append(currentString, rs.GetRune())
				rs.NextRune()
			}
			if !rs.IsWhiteSpace() && !rs.IsHashtag() && !rs.IsApostrophe() {
				for !rs.IsWhiteSpace() {
					rs.NextRune()
				}
				rs.NextRune()
				curPosition := rs.GetCurrentPosition()
				curPosition.column--

				return NewError("the character code has bad end", NewFragment(start, rs.GetCurrentPosition()))
			}
		} else {
			for !rs.IsApostrophe() && !rs.IsLineTranslation() {
				rs.NextRune()
			}
			curPosition := rs.GetCurrentPosition()
			rs.NextRune()
			return NewError("symbol is not a section or character code", NewFragment(start, curPosition))
		}
	}

	curPosition := rs.GetCurrentPosition()
	rs.NextRune()

	return NewString(string(currentString),
		NewFragment(start, curPosition))
}

func processIdentifier(rs *RunePosition) IToken {
	currentIdent := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()

	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			// Строим слово
			currentIdent = append(currentIdent, rs.GetRune())
		} else {
			// ошибка в Identifier, не буквы или цифры
			curPosition := rs.GetCurrentPosition()
			for !rs.IsWhiteSpace() {
				curPosition = rs.GetCurrentPosition()
				rs.NextRune()
			}
			return NewError("symbol is not a latin letter or digit", NewFragment(start, curPosition))
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}
	return NewIdent(string(currentIdent), NewFragment(start, curPositionToken))
}

func processStartError(rs *RunePosition) IToken {

	startPosition := rs.GetCurrentPosition()
	curPosition := rs.GetCurrentPosition()
	for !rs.IsWhiteSpace() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		curPosition = rs.GetCurrentPosition()
		rs.NextRune()
	}
	return NewError("start is not a quote or letter", NewFragment(startPosition, curPosition))
}

```

Файл `tag.go`
```go
package lab_lexer

type DomainTag int

const (
	IdentTag DomainTag = iota + 1
	StrTag
	ErrTag
	EopTag
)

var tagToString = map[DomainTag]string{
	IdentTag: "IDENTIFIER",
	StrTag:   "STRING",
	ErrTag:   "ERROR",
	EopTag:   "EOP",
}

```

Файл `tokens.go`
```go
package lab_lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type IToken interface {
	GetType() DomainTag
}

type Token struct {
	Type       DomainTag
	Value      string
	Coordinate Fragment
}

type Fragment struct {
	start  TokenPosition
	finish TokenPosition
}

func NewFragment(start, finish TokenPosition) Fragment {
	return Fragment{
		start:  start,
		finish: finish,
	}
}

func (f Fragment) String() string {
	return fmt.Sprintf("%s-%s", f.start.String(), f.finish.String())
}

type TokenPosition struct {
	line   int
	column int
}

func NewTokenPosition(line, column int) TokenPosition {
	return TokenPosition{
		line:   line,
		column: column,
	}
}

func (p *TokenPosition) String() string {
	return fmt.Sprintf("(%d,%d)", p.line, p.column)
}

func NewToken(tag DomainTag, value string, coordinate Fragment) Token {
	return Token{
		Type:       tag,
		Value:      value,
		Coordinate: coordinate,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s: %s", tagToString[t.Type], t.Coordinate, t.Value)
}

func (t Token) GetType() DomainTag {
	return t.Type
}

type Message struct {
	isError bool
	text    string
}

func NewMessage(isError bool, text string) Message {
	return Message{
		isError: isError,
		text:    text,
	}
}

type ErrorToken struct {
	Token
}

func NewError(text string, fragment Fragment) ErrorToken {
	return ErrorToken{
		Token: NewToken(ErrTag, text, fragment),
	}
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type IdentToken struct {
	Token
	attr int
}

func NewIdent(value string, fragment Fragment) IdentToken {
	return IdentToken{
		Token: NewToken(IdentTag, value, fragment),
	}
}

func (it *IdentToken) SetAttr(attr int) {
	it.attr = attr
}

func (it IdentToken) String() string {
	return fmt.Sprintf("%s %s: %d", tagToString[it.Type], it.Coordinate, it.attr)
}

type StringToken struct {
	Token
	attr string
}

func NewString(value string, fragment Fragment) StringToken {
	return StringToken{
		Token: NewToken(StrTag, value, fragment),
	}
}

func (st *StringToken) SetText(text string) {
	newText := processText(text)
	st.attr = newText
}

func processText(text string) string {
	newTexts := strings.Split(text, "#")
	ans := ""
	for _, nt := range newTexts {
		if len(nt) == 0 {
		} else if nt[0] == '\'' {
			nt = strings.Replace(nt, "''", "'", -1)
			ans += nt[1 : len(nt)-1]
		} else {
			num, err := strconv.Atoi(nt)
			if err != nil {
				fmt.Println("num problem")
			}
			ans += string(rune(num))
		}
	}
	return ans
}

func (st StringToken) String() string {
	return fmt.Sprintf("%s %s: %s", tagToString[st.Type], st.Coordinate, st.attr)
}

```



# Тестирование

Для тестирования было написано несколько тестирующих файлов,
покрывающих большую часть тестовых случаев

Входные данные

Файл `ident_correct.txt`
```
abcdf abcdf1      a12345678
     a1 b23
aBC12345qwer a
a1a1a1
```

Файл `ident_error.txt`
```
abcdf 1abcdf1 1a12345678 ab
a1 3b23
abcфыъ12345qwer
a1a1a1)) asd
{000}
pppu
```

Файл `mixed.txt`
```
'asdf''hjk'#33 aa 'aaa aaa'
abcdf abcdf1 3333 b#16#17#18
#10
```

Файл `strings_correct.txt`
```
'asdf' #10
'asdf''hjk'#15 #16#17#18
'' 'asd fgh 20 klj: qwerty'
'asd '' asd' 'asd''asd' ''''
'hello'#10#13'world'
```

Файл `strings_error.txt`
```
'asdf' #10
'asdf''hjk'#15aa b#16#17#18
'' 'asd fgh 20 klj: qwerty'
'asd ' asd' 'asd''asd'
'hello' #10 #13'world'
```

Вывод на `stdout` для *mixed.txt*

```
STRING (1,1)-(1,15): asdf'hjk!
IDENTIFIER (1,16)-(1,17): 0
STRING (1,19)-(1,28): aaa aaa
IDENTIFIER (2,1)-(2,5): 1
IDENTIFIER (2,7)-(2,12): 2
STRING (3,1)-(3,4): 

_____MESSAGES_____
Type: Error | i: 0 | position: (2,14)-(2,17) | text: start is not a quote or letter
Type: Error | i: 1 | position: (2,19)-(2,28) | text: symbol is not a latin letter or digit
____Identifiers____
IDENTIFIER (1,16)-(1,17) 0 -- aa
IDENTIFIER (2,1)-(2,5) 1 -- abcdf
IDENTIFIER (2,7)-(2,12) 2 -- abcdf1
finish
```

# Вывод
В ходе данной лабораторной работы был получен навык реализации лексического анализатора на
объектно-ориентированном языке без применения каких-либо средств автоматизации
решения задачи лексического анализа.
Опробовал объектный подход в Golang, для этого потребовалось знание интерфейсов и встраивания структур.
При этом во время написания этого лексического анализатора мне показалось удобнее использовать 
самописные средства для решения задачи лексического анализа, чем регулярные выражения, но
при этом объем кода и сложность возросла. 

