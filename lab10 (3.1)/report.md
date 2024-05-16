% Лабораторная работа № 3.1 «Самоприменимый генератор компиляторов
на основе предсказывающего анализа»
% 15 мая 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является изучение алгоритма построения таблиц предсказывающего анализатора.

# Индивидуальный вариант
```
/* аксиома помечена звёздочкой */
F  ("n") ("(" E ")")
T  (F T')
T' ("*" F T') ()
* E  (T E')
  E' ("+" T E') ()
```

# Реализация

## Модуль калькулятора

### Интерпретатор дерева вывода

`tree_converter.go`
```go
package interpreteter

import (
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
				left := evaluateExpression(n.Children[0])
				right := evaluateExpression(n.Children[1])
				return left + right
			} else if len(n.Children) == 1 {
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в E равна не 1 и не 2:", n, len(n.Children))
			}
		} else if n.NonTerminal == "E'" {
			if len(n.Children) == 3 {
				left := evaluateExpression(n.Children[1])
				right := evaluateExpression(n.Children[2])
				return left + right
			} else if len(n.Children) == 2 {
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 0 {
				return 0 // epsilon
			} else {
				log.Fatal("Длина в E' равна не 3 и не 0:", n, len(n.Children))
			}
		} else if n.NonTerminal == "T" {
			if len(n.Children) == 2 {
				left := evaluateExpression(n.Children[0])
				right := evaluateExpression(n.Children[1])
				return left * right
			} else if len(n.Children) == 1 {
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в T равна не 1 и не 2:", n, len(n.Children))
			}
		} else if n.NonTerminal == "T'" {
			if len(n.Children) == 3 {
				left := evaluateExpression(n.Children[1])
				right := evaluateExpression(n.Children[2])
				return left * right
			} else if len(n.Children) == 2 {
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 0 {
				return 1 // epsilon
			} else {
				log.Fatal("Длина в T' равна не 3 и не 2 и не 0:", n, len(n.Children))
			}
		} else if n.NonTerminal == "F" {
			if len(n.Children) == 3 {
				return evaluateExpression(n.Children[1])
			} else if len(n.Children) == 1 {
				return evaluateExpression(n.Children[0])
			} else {
				log.Fatal("Длина в F равна не 3 и не 1:", n, len(n.Children))
			}
		} else {
			log.Fatal("Неизвестный нетерминал", n, len(n.Children))
		}
	case *top_down_parse.LeafTreeNode:
		value, _ := strconv.Atoi(n.Token.GetValue())
		return value
	}
	return 0
}

```

### Лексер

`position.go`
```go
package lexer

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

func (p *RunePosition) IsDigit() bool {
	return unicode.IsDigit(p.GetRune())
}

func (p *RunePosition) IsOpenBracket() bool {
	return p.GetRune() == '('
}

func (p *RunePosition) IsCloseBracket() bool {
	return p.GetRune() == ')'
}

func (p *RunePosition) IsPlus() bool {
	return p.GetRune() == '+'
}

func (p *RunePosition) IsMultiply() bool {
	return p.GetRune() == '*'
}

func (p *RunePosition) IsBrackets() bool {
	return p.IsOpenBracket() || p.IsCloseBracket()
}

```

`scanner.go`
```go
package lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Scanner struct {
	tokens   []IToken
	position int
}

func NewScanner(filepath string) *Scanner {
	return &Scanner{
		tokens:   ParseFile(filepath),
		position: 0,
	}
}

func (s *Scanner) PrintTokens() {
	fmt.Println("TOKENS:")
	for _, p := range s.tokens {
		fmt.Println(p)
	}
}

func (s *Scanner) NextToken() IToken {
	token := s.tokens[s.position]

	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == IntTag ||
		token.GetType() == CloseBracketTag || token.GetType() == OpenBracketTag ||
		token.GetType() == PlusTag || token.GetType() == MultiplyTag {
		return token
	}

	return token
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

		if runeScanner.IsDigit() {
			tokens = append(tokens, processInteger(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() ||
			runeScanner.IsPlus() || runeScanner.IsMultiply() {
			tokens = append(tokens, processOperand(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				fmt.Println("rune:", string(runeScanner.GetRune()))
				log.Fatalf("incorrect rune in parser")
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processInteger(rs *RunePosition) IToken {
	currentInt := make([]rune, 0)
	start := rs.GetCurrentPosition()

	for rs.IsDigit() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentInt = append(currentInt, rs.GetRune())
		rs.NextRune()
	}

	curPosition := rs.GetCurrentPosition()

	return NewIntegerToken(string(currentInt), NewFragment(start, curPosition))
}

func processOperand(rs *RunePosition) IToken {
	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		return NewOpenBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		return NewCloseBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		return NewMultiply(string(operand), NewFragment(start, curPosition))
	} else if operand == '+' {
		return NewPlus(string(operand), NewFragment(start, curPosition))
	}

	log.Fatalf("the non real error %v", NewFragment(start, rs.GetCurrentPosition()))
	return nil
}

```

`tag.go`
```go
package lexer

type DomainTag int

const (
	IntTag DomainTag = iota + 1
	OpenBracketTag
	CloseBracketTag
	PlusTag
	MultiplyTag
	EopTag
)

var TagToString = map[DomainTag]string{
	IntTag:          "n",
	OpenBracketTag:  "(",
	CloseBracketTag: ")",
	PlusTag:         "+",
	MultiplyTag:     "*",
	EopTag:          "Eop",
}

```

`tokens.go`
```go
package lexer

import (
	"fmt"
)

type IToken interface {
	GetType() DomainTag
	GetValue() string
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
	return fmt.Sprintf("%s %s: %s", TagToString[t.Type], t.Coordinate, t.Value)
}

func (t Token) GetType() DomainTag {
	return t.Type
}

func (t Token) GetValue() string {
	return t.Value
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type IntegerToken struct {
	Token
}

func NewIntegerToken(value string, fragment Fragment) IntegerToken {
	return IntegerToken{
		Token: NewToken(IntTag, value, fragment),
	}
}

func (st IntegerToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type OpenBracketToken struct {
	Token
}

func NewOpenBracket(value string, fragment Fragment) OpenBracketToken {
	return OpenBracketToken{
		Token: NewToken(OpenBracketTag, value, fragment),
	}
}

func (st OpenBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type CloseBracketToken struct {
	Token
}

func NewCloseBracket(value string, fragment Fragment) CloseBracketToken {
	return CloseBracketToken{
		Token: NewToken(CloseBracketTag, value, fragment),
	}
}

func (st CloseBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type PlusToken struct {
	Token
}

func NewPlus(value string, fragment Fragment) PlusToken {
	return PlusToken{
		Token: NewToken(PlusTag, value, fragment),
	}
}

func (st PlusToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type MultiplyToken struct {
	Token
}

func NewMultiply(value string, fragment Fragment) MultiplyToken {
	return MultiplyToken{
		Token: NewToken(MultiplyTag, value, fragment),
	}
}

func (st MultiplyToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

```

### top-down-parse

`gen_table.go`
```go
package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"E - n":  {"T", "E'"},
		"E - )":  {"T", "E'"},
		"E' - +": {"+", "T", "E'"},
		"E' - )": {},
		"E' - $": {},
		"F - n":  {"n"},
		"F - (":  {"(", "E", ")"},
		"T - n":  {"F", "T'"},
		"T - (":  {"F", "T'"},
		"T' - *": {"*", "F", "T'"},
		"T' - +": {},
		"T' - )": {},
		"T' - $": {},
	}
}

func newGenAxiom() string {
	return "E"
}

```

`parser.go`
```go
package top_down_parse

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
)

type Parser struct {
	table map[string][]string
	axiom string
}

func NewParser() Parser {
	return Parser{
		table: newGenTable(),
		axiom: newGenAxiom(),
	}
}

func (p Parser) TopDownParse(scanner *lexer.Scanner) (*TreeNode, error) {
	type stackNode struct {
		itn *InnerTreeNode
		val string
	}
	s := NewStack[stackNode]()

	root := newTreeNode()
	root.addNode(newInnerTreeNode(""))

	s.Push(stackNode{itn: root.Root, val: p.axiom})

	t := scanner.NextToken()

	for t.GetType() != lexer.EopTag {
		topNode, err := s.Pop()
		if err != nil {
			return newTreeNode(), fmt.Errorf("failed to get top node: %w", err)
		}

		if isTerminal(topNode.val) {
			topNode.itn.Children = append(topNode.itn.Children, newLeafTreeNode(t))
			t = scanner.NextToken()
		} else if neighbourhoods, ok := p.table[newTableKey(topNode.val,
			lexer.TagToString[t.GetType()])]; ok {
			in := newInnerTreeNode(topNode.val)
			topNode.itn.Children = append(topNode.itn.Children, in)

			for i := len(neighbourhoods) - 1; i >= 0; i-- {
				s.Push(stackNode{itn: in, val: neighbourhoods[i]})
			}
		} else {
			return newTreeNode(), fmt.Errorf("failed do parse in table with val %s and token %s",
				topNode.val, t.GetValue())
		}
	}

	return root, nil
}

```

`stack.go`
```go
package top_down_parse

import "errors"

type Stack[T any] struct {
	buffer []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		buffer: make([]T, 0),
	}
}

func (s *Stack[T]) Push(elem T) {
	s.buffer = append(s.buffer, elem)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.buffer) > 0 {
		elem := s.buffer[len(s.buffer)-1]
		s.buffer = s.buffer[:len(s.buffer)-1]
		return elem, nil
	}
	var tmp T
	return tmp, errors.New("empty buffer")
}

func (s *Stack[T]) GetElems() []T {
	return s.buffer
}

```

`table.go`
```go
package top_down_parse

import (
	"fmt"
	"slices"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{"E", "E'", "T", "T'", "F"}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}

```

`tree.go`
```go
package top_down_parse

import (
	"fmt"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
)

type TreeNode struct {
	Root *InnerTreeNode
}

func newTreeNode() *TreeNode {
	return &TreeNode{}
}

func (tn *TreeNode) Print() {
	tn.Root.printNode(0)
}

func (tn *TreeNode) addNode(node *InnerTreeNode) {
	tn.Root = node
}

type TreeNodePrinter interface {
	printNode(offset int)
}

type InnerTreeNode struct {
	NonTerminal string
	Children    []TreeNodePrinter
}

func newInnerTreeNode(nonTerminal string) *InnerTreeNode {
	return &InnerTreeNode{NonTerminal: nonTerminal, Children: make([]TreeNodePrinter, 0)}
}

func (itn InnerTreeNode) printNode(offset int) {
	fmt.Printf(strings.Repeat("..", offset) + fmt.Sprintf("Inner node: %s\n", itn.NonTerminal))

	for _, child := range itn.Children {
		child.printNode(offset + 1)
	}
}

type LeafTreeNode struct {
	Token lexer.IToken
}

func newLeafTreeNode(token lexer.IToken) *LeafTreeNode {
	return &LeafTreeNode{Token: token}
}

func (ltn LeafTreeNode) printNode(offset int) {
	if ltn.Token.GetType() == lexer.IntTag {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s - %s\n",
				lexer.TagToString[ltn.Token.GetType()], ltn.Token.GetValue()))
	} else {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s\n", lexer.TagToString[ltn.Token.GetType()]))
	}
}

```

### main

`main.go`
```go
package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/interpreteter"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/calculator/top_down_parse"
)

const filepath = "example.txt"

func main() {
	log.Println("start calculator")

	scanner := lexer.NewScanner(filepath)

	scanner.PrintTokens()

	parser := top_down_parse.NewParser()

	tree, errTDP := parser.TopDownParse(scanner)
	if errTDP != nil {
		log.Fatalf("err in TopDownParse %+v", errTDP)
	}

	tree.Print()

	fmt.Println("solver:", interpreteter.Solve(tree))

	log.Println("end calculator")
}

```

`go.mod`
```go
module github.com/VyacheslavIsWorkingNow/cd/lab10/calculator

go 1.22

```


## Модуль конвертора


### Лексер

`compilers.go`
```go
package lexer

import (
	"fmt"
	"sort"
)

type Compiler struct {
	tokens   []Token
	messages map[Fragment]Message
}

func NewCompiler() *Compiler {
	return &Compiler{
		messages: make(map[Fragment]Message),
	}
}

func (c *Compiler) AddMessage(ct CommentToken) {
	c.messages[ct.Coordinate] = NewMessage(true, ct.Value)
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
		fmt.Printf("Type: Comment | i: %d | position: %v | text: %s\n",
			i, position, c.messages[position].text)
	}
}

```

`position.go`
```go
package lexer

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
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func (p *RunePosition) IsDigit() bool {
	return unicode.IsDigit(p.GetRune())
}

func (p *RunePosition) IsLowLine() bool {
	return p.GetRune() == '_'
}

func (p *RunePosition) IsQuote() bool {
	return p.GetRune() == '"'
}
func (p *RunePosition) IsOneQuote() bool {
	return p.GetRune() == '\''
}

func (p *RunePosition) IsOpenBracket() bool {
	return p.GetRune() == '('
}

func (p *RunePosition) IsCloseBracket() bool {
	return p.GetRune() == ')'
}

func (p *RunePosition) IsStar() bool {
	return p.GetRune() == '*'
}

func (p *RunePosition) IsBrackets() bool {
	return p.IsOpenBracket() || p.IsCloseBracket()
}

func (p *RunePosition) IsOpenSlash() bool {
	return p.GetRune() == '/'
}

func (p *RunePosition) IsLineTranslation() bool {
	return p.GetRune() == '\n'
}

```

`scanner.go`
```go
package lexer

import (
	"bufio"
	"fmt"
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

func (s *Scanner) PrintTokens() {
	fmt.Println("TOKENS:")
	for _, p := range s.tokens {
		fmt.Println(p)
	}
}

func (s *Scanner) NextToken() IToken {
	token := s.tokens[s.position]

	if token.GetType() != EopTag {
		s.position++
	}
	if token.GetType() == TermTag || token.GetType() == NonTermTag ||
		token.GetType() == CloseBracketTag ||
		token.GetType() == OpenBracketTag || token.GetType() == AxiomTag {
		return token
	}
	if token.GetType() == CommentTag {
		newToken := token.(CommentToken)
		s.compiler.AddMessage(newToken)
		return s.NextToken()
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

		if runeScanner.IsQuote() {
			tokens = append(tokens, processTerminal(runeScanner))
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() ||
			runeScanner.IsStar() {
			tokens = append(tokens, processOperand(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processNonTerminal(runeScanner))
		} else if runeScanner.IsOpenSlash() {
			tokens = append(tokens, processComment(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				fmt.Println("rune:", string(runeScanner.GetRune()))
				log.Fatalf("incorrect rune in parser")
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processTerminal(rs *RunePosition) IToken {
	currentString := make([]rune, 0)
	currentString = append(currentString, rs.GetRune())
	start := rs.GetCurrentPosition()
	rs.NextRune()

	for !rs.IsQuote() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentString = append(currentString, rs.GetRune())
		rs.NextRune()
	}

	currentString = append(currentString, rs.GetRune())
	curPosition := rs.GetCurrentPosition()
	rs.NextRune()

	term := string(currentString[1 : len(currentString)-1])

	return NewTerminal(term, NewFragment(start, curPosition))
}

func processOperand(rs *RunePosition) IToken {
	start, curPosition := rs.GetCurrentPosition(), rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()

	if operand == '(' {
		return NewOpenBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == ')' {
		return NewCloseBracket(string(operand), NewFragment(start, curPosition))
	} else if operand == '*' {
		return NewAxiom(string(operand), NewFragment(start, curPosition))
	}

	log.Fatalf("the non real error %v", NewFragment(start, rs.GetCurrentPosition()))
	return nil
}

func processNonTerminal(rs *RunePosition) IToken {
	currentNonTerminal := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()
	currentNonTerminal = append(currentNonTerminal, rs.GetRune())
	rs.NextRune()

	for !rs.IsWhiteSpace() && !rs.IsBrackets() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() || rs.IsLowLine() {
			currentNonTerminal = append(currentNonTerminal, rs.GetRune())
		} else if rs.IsOneQuote() {
			currentNonTerminal = append(currentNonTerminal, rs.GetRune())
		} else {
			log.Fatalf("error in process nonterminal")
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}

	return NewNonTerminal(string(currentNonTerminal), NewFragment(start, curPositionToken))
}

func processComment(rs *RunePosition) IToken {

	currentText := make([]rune, 0)

	start := rs.GetCurrentPosition()
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	if !rs.IsStar() {
		log.Fatalf("missing quote in comment")
	}
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	for !rs.IsStar() {
		if rs.GetRune() == -1 {
			log.Fatalf("the line didn't end %v", NewFragment(start, rs.GetCurrentPosition()))
		}
		currentText = append(currentText, rs.GetRune())
		rs.NextRune()
	}

	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	if !rs.IsOpenSlash() {
		log.Fatalf("missing slash in comment")
	}
	currentText = append(currentText, rs.GetRune())
	rs.NextRune()

	curPositionToken := rs.GetCurrentPosition()

	return NewComment(string(currentText), NewFragment(start, curPositionToken))
}

```

`tag.go`
```go
package lexer

type DomainTag int

const (
	NonTermTag DomainTag = iota + 1
	TermTag
	OpenBracketTag
	CloseBracketTag
	AxiomTag
	CommentTag
	EopTag
)

var TagToString = map[DomainTag]string{
	NonTermTag:      "NonTerminal",
	TermTag:         "Terminal",
	OpenBracketTag:  "OpenBracket",
	CloseBracketTag: "CloseBracket",
	AxiomTag:        "AxiomSign",
	CommentTag:      "Comment",
	EopTag:          "Eop",
}

```

`tokens.go`
```go
package lexer

import (
	"fmt"
)

type IToken interface {
	GetType() DomainTag
	GetValue() string
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
	return fmt.Sprintf("%s %s: %s", TagToString[t.Type], t.Coordinate, t.Value)
}

func (t Token) GetType() DomainTag {
	return t.Type
}

func (t Token) GetValue() string {
	return t.Value
}

type Message struct {
	isComment bool
	text      string
}

func NewMessage(isComment bool, text string) Message {
	return Message{
		isComment: isComment,
		text:      text,
	}
}

type CommentToken struct {
	Token
}

func NewComment(text string, fragment Fragment) CommentToken {
	return CommentToken{
		Token: NewToken(CommentTag, text, fragment),
	}
}

func (ct CommentToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ct.Type], ct.Coordinate, ct.Value)
}

type EOPToken struct {
	Token
}

func NewEOP() EOPToken {
	return EOPToken{
		Token: NewToken(EopTag, "end of file", Fragment{}),
	}
}

type NonTermToken struct {
	Token
}

func NewNonTerminal(value string, fragment Fragment) NonTermToken {
	return NonTermToken{
		Token: NewToken(NonTermTag, value, fragment),
	}
}

func (ntt NonTermToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ntt.Type], ntt.Coordinate, ntt.Value)
}

type TerminalToken struct {
	Token
}

func NewTerminal(value string, fragment Fragment) TerminalToken {
	return TerminalToken{
		Token: NewToken(TermTag, value, fragment),
	}
}

func (tt TerminalToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[tt.Type], tt.Coordinate, tt.Value)
}

type OpenBracketToken struct {
	Token
}

func NewOpenBracket(value string, fragment Fragment) OpenBracketToken {
	return OpenBracketToken{
		Token: NewToken(OpenBracketTag, value, fragment),
	}
}

func (st OpenBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type CloseBracketToken struct {
	Token
}

func NewCloseBracket(value string, fragment Fragment) CloseBracketToken {
	return CloseBracketToken{
		Token: NewToken(CloseBracketTag, value, fragment),
	}
}

func (st CloseBracketToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type AxiomToken struct {
	Token
}

func NewAxiom(value string, fragment Fragment) AxiomToken {
	return AxiomToken{
		Token: NewToken(AxiomTag, value, fragment),
	}
}

func (st AxiomToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

```

### Генерация таблицы

`downoload.go`
```go
package predtable

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
)

var tableView = `
package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
%s
	}
}

`

var axiomView = `
func newGenAxiom() string {
	return "%s"
}

`

func UploadTableToFile(filepath string, genTable map[string][]string, rules semantic.Rules) error {
	newFileView := makeTableView(genTable) + makeAxiomView(rules.Axiom)
	return uploadNewFile(filepath, newFileView)
}

func makeTableView(genTable map[string][]string) string {
	keys := make([]string, 0, len(genTable))
	for key := range genTable {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var lines []string
	for _, key := range keys {
		values := genTable[key]
		var quotedValues []string
		for _, value := range values {
			if value != "" {
				quotedValues = append(quotedValues, fmt.Sprintf(`"%s"`, value))
			}
		}
		line := fmt.Sprintf("\t\t\"%s\": {%s},\n", key, strings.Join(quotedValues, ", "))
		lines = append(lines, line)
	}

	return fmt.Sprintf(tableView, strings.Join(lines, ""))
}

func makeAxiomView(axiom string) string {
	return fmt.Sprintf(axiomView, axiom)
}

func uploadNewFile(filepath, view string) error {
	err := os.WriteFile(filepath, []byte(view), 0644)
	if err != nil {
		return fmt.Errorf("failed uploud view: %w", err)
	}
	return nil
}

```

`find_first_follow.go`
```go
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

```

`gen_table.go`
```go
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
					fmt.Errorf("rules is not LL1: two rules in one cell %s",
						makeKey(r.LeftSymbol, term))
			}
			table[makeKey(r.LeftSymbol, term)] = r.RightSymbols
		}
		
		if isEpsilon {
			for term := range rFollow {
				if _, ok := table[makeKey(r.LeftSymbol, term)]; ok {
					return table,
						fmt.Errorf("rules is not LL1: two rules in one cell %s",
							makeKey(r.LeftSymbol, term))
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

```


### Семантический анализ

`semantic.go`
```go
package semantic

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
)

type Rules struct {
	Rule        []Rule
	Axiom       string
	Terminal    map[string]struct{}
	NonTerminal map[string]struct{}
}

type Rule struct {
	LeftSymbol   string
	RightSymbols []string
}

type Semantic struct {
	Tree *top_down_parse.TreeNode
}

func NewSemantic(tree *top_down_parse.TreeNode) *Semantic {
	return &Semantic{
		Tree: tree,
	}
}

func (s *Semantic) StartSemanticAnalysis() (Rules, error) {
	log.Println("start semantic analysis")

	allNonTerminalSymbol, terminalSymbol, ac := getTerminalAndNonTerminal(*s.Tree)

	if ac == 0 {
		return Rules{}, fmt.Errorf("zero axiom, need one")
	}
	if ac > 1 {
		return Rules{}, fmt.Errorf("axiom isn`t be better than 1, give: %d", ac)
	}

	rules := Rules{
		Rule:        make([]Rule, 0),
		Axiom:       "",
		NonTerminal: allNonTerminalSymbol,
		Terminal:    terminalSymbol,
	}

	leftNonTerminals := make(map[string]struct{})

	convertTreeToRewritingsRules(*s.Tree, &rules, &leftNonTerminals)

	if !isFirstSetInSecond(allNonTerminalSymbol, leftNonTerminals) {
		return Rules{}, fmt.Errorf("there are unreachable nonterminals %+v, %+v",
			leftNonTerminals, allNonTerminalSymbol)
	}

	convertEmptiness(&rules.Rule)

	return rules, nil
}

func getTerminalAndNonTerminal(tree top_down_parse.TreeNode) (
	nonTerminal map[string]struct{},
	terminal map[string]struct{},
	axiomCount int,
) {
	nonTerminal = make(map[string]struct{})
	terminal = make(map[string]struct{})

	traverseTree(tree.Root, &nonTerminal, &terminal, &axiomCount)

	return nonTerminal, terminal, axiomCount
}

func traverseTree(
	node top_down_parse.TreeNodePrinter, nonTerminals, terminals *map[string]struct{}, axiomCount *int,
) {

	switch n := node.(type) {
	case *top_down_parse.InnerTreeNode:
		// Тут я нахожусь во внутреннем узле - это рабочее пространство дерева
		for _, child := range n.Children {
			traverseTree(child, nonTerminals, terminals, axiomCount)
		}
	case *top_down_parse.LeafTreeNode:
		// Либо терминал, либо нетерминал, либо служебные символы
		if n.Token.GetType() == lexer.TermTag {
			(*terminals)[n.Token.GetValue()] = struct{}{}
		} else if n.Token.GetType() == lexer.NonTermTag {
			(*nonTerminals)[n.Token.GetValue()] = struct{}{}
		} else if n.Token.GetType() == lexer.AxiomTag {
			*axiomCount++
		}
	default:
		log.Println("default?", n)
	}

}

func convertTreeToRewritingsRules(
	tree top_down_parse.TreeNode, rules *Rules, leftNonTerminals *map[string]struct{},
) {
	log.Println("start convert")

	root, errRoot := checkDeclarationNode(tree.Root.Children[0])
	if errRoot != nil {
		log.Printf("error in check root %+v", errRoot)
	}

	handleDeclarations(root, rules, leftNonTerminals)
}

func checkInnerNode(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, ok := node.(*top_down_parse.InnerTreeNode)

	if !ok {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node declaration %t, %s", ok, innerNode.NonTerminal)
	}

	return innerNode, nil
}

func checkDeclarationNode(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Declarations || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node declaration %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleDeclarations(
	node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminals *map[string]struct{},
) {
	if len(node.Children) == 2 {
		// Первый - RewritingRule, Второй - Declaration
		rewritingRule, errRR := checkRewritingRuleNode(node.Children[0])
		if errRR != nil {
			log.Printf("error in check handle declaration %+v", errRR)
		}
		handleRewritingRule(rewritingRule, rules, leftNonTerminals)
		declaration, errD := checkDeclarationNode(node.Children[1])
		if errD != nil {
			log.Printf("error in check handle declaration %+v", errD)
		}
		handleDeclarations(declaration, rules, leftNonTerminals)
	} else if len(node.Children) == 1 {
		// Первый - RewritingRule
		rewritingRule, errRR := checkRewritingRuleNode(node.Children[0])
		if errRR != nil {
			log.Printf("error in check handle declaration %+v", errRR)
		}
		handleRewritingRule(rewritingRule, rules, leftNonTerminals)
	} else {
		log.Println("Длина не один и не два в handleDeclarations ???")
	}
}

func checkRewritingRuleNode(
	node top_down_parse.TreeNodePrinter,
	) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.RewritingRule || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting rule %s, %+v",
				innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleRewritingRule(
	node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminals *map[string]struct{},
) {
	if len(node.Children) == 3 {
		// Первый - Axiom, Второй - NonTerminal Leaf, Третий - REWRITING
		axiom, errA := getLeafValue(node.Children[0])
		if errA != nil {
			log.Printf("error in check handle rewriting rule axiom %+v", errA)
		}
		if !isAxiom(axiom) {
			log.Printf("error in axiom checker. Axiom doesn't have axiom tag")
		}

		nonTerminal, errNT := getLeafValue(node.Children[1])
		if errNT != nil {
			log.Printf("error in check handle rewriting rule non terminal %+v", errNT)
		}

		putToLeftNonTerminalTable(nonTerminal, leftNonTerminals)

		rules.Axiom = nonTerminal.GetValue()

		rewriting, errR := checkRewriting(node.Children[2])
		if errR != nil {
			log.Printf("error in check handle rewriting rule rewriting %+v", errR)
		}
		handleRewriting(rewriting, rules, nonTerminal.GetValue())
	} else if len(node.Children) == 2 {
		nonTerminal, errNT := getLeafValue(node.Children[0])
		if errNT != nil {
			log.Printf("error in check handle rewriting rule non terminal %+v", errNT)
		}

		putToLeftNonTerminalTable(nonTerminal, leftNonTerminals)

		rewriting, errR := checkRewriting(node.Children[1])
		if errR != nil {
			log.Printf("error in check handle rewriting rule rewriting %+v", errR)
		}
		handleRewriting(rewriting, rules, nonTerminal.GetValue())
	} else {
		log.Println("Длина не два и не три в handleRewritingRule ???")
	}
}

func checkRewriting(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Rewriting || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func checkLeaf(node top_down_parse.TreeNodePrinter) (*top_down_parse.LeafTreeNode, error) {
	leafNode, ok := node.(*top_down_parse.LeafTreeNode)

	if !ok {
		return &top_down_parse.LeafTreeNode{},
			fmt.Errorf("error in cheking leaf %t", ok)
	}

	return leafNode, nil
}

func getLeafValue(node top_down_parse.TreeNodePrinter) (lexer.IToken, error) {
	leaf, err := checkLeaf(node)
	if err != nil {
		return lexer.Token{},
			fmt.Errorf("error in get leaf value %+v", err)
	}

	return leaf.Token, nil
}

func isAxiom(t lexer.IToken) bool {
	return t.GetType() == lexer.AxiomTag
}

func putToLeftNonTerminalTable(t lexer.IToken, NonTerminalTable *map[string]struct{}) {
	if t.GetType() != lexer.NonTermTag {
		log.Println("nonTerminal has nonTerminal tag", t.GetType(), t.GetValue())
	}
	(*NonTerminalTable)[t.GetValue()] = struct{}{}
}

func handleRewriting(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminal string) {
	if len(node.Children) == 4 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket, Четвертый - REWRITING_OPT

		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting close bracket %+v", errCB)
		}

		rewritingOpt, errRO := checkRewritingOpt(node.Children[3])
		if errRO != nil {
			log.Printf("error in check handle rewriting rewriting opt %+v", errRO)
		}
		handleRewritingOpt(rewritingOpt, rules, leftNonTerminal)
	} else {
		log.Println("Длина четыре в handleRewriting ???")
	}
}

func checkRewritingOpt(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.RewritingOpt || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node rewriting %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func checkOpenBracketLeaf(node top_down_parse.TreeNodePrinter) error {
	leaf, err := checkLeaf(node)

	if leaf.Token.GetType() != lexer.OpenBracketTag || err != nil {
		return fmt.Errorf("error in cheking leaf open bracket %s, %+v", leaf.Token, err)
	}

	return nil
}

func checkCloseBracketLeaf(node top_down_parse.TreeNodePrinter) error {
	leaf, err := checkLeaf(node)

	if leaf.Token.GetType() != lexer.CloseBracketTag || err != nil {
		return fmt.Errorf("error in cheking leaf close bracket %s, %+v", leaf.Token, err)
	}

	return nil
}

func checkBody(node top_down_parse.TreeNodePrinter) (*top_down_parse.InnerTreeNode, error) {
	innerNode, err := checkInnerNode(node)

	if innerNode.NonTerminal != top_down_parse.Body || err != nil {
		return &top_down_parse.InnerTreeNode{},
			fmt.Errorf("error in cheking inner node body %s, %+v", innerNode.NonTerminal, err)
	}

	return innerNode, nil
}

func handleRewritingOpt(node *top_down_parse.InnerTreeNode, rules *Rules, leftNonTerminal string) {
	if len(node.Children) == 4 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket, Четвертый - REWRITING_OPT

		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting opt body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt close bracket %+v", errCB)
		}

		rewritingOpt, errRO := checkRewritingOpt(node.Children[3])
		if errRO != nil {
			log.Printf("error in check handle rewriting opt rewriting opt %+v", errRO)
		}
		handleRewritingOpt(rewritingOpt, rules, leftNonTerminal)
	} else if len(node.Children) == 3 {
		// Первый - OpenBracket, Второй - BODY, Третий - CloseBracket
		errOB := checkOpenBracketLeaf(node.Children[0])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt open bracket %+v", errOB)
		}

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle rewriting opt body %+v", errB)
		}
		currentBody := make([]string, 0)
		handleBody(body, &currentBody)

		rules.putRule(leftNonTerminal, currentBody)

		errCB := checkCloseBracketLeaf(node.Children[2])
		if errOB != nil {
			log.Printf("error in check handle rewriting opt close bracket %+v", errCB)
		}
	} else if len(node.Children) == 0 {
		// Ничего не делаем
	} else {
		log.Println("Длина не четыре и не три, и не ноль в handleRewritingOpt ???")
	}
}

func handleBody(node *top_down_parse.InnerTreeNode, currentBody *[]string) {
	if len(node.Children) == 2 {
		// Первый - Лист, Второй - BODY
		token, errT := getLeafValue(node.Children[0])
		if errT != nil {
			log.Printf("error in check handle body leaf %+v", errT)
		}
		*currentBody = append(*currentBody, token.GetValue())

		body, errB := checkBody(node.Children[1])
		if errB != nil {
			log.Printf("error in check handle body body %+v", errB)
		}
		handleBody(body, currentBody)

	} else if len(node.Children) == 0 {
		// Ничего не делаем
	} else {
		log.Println("Длина не два и не ноль в handleBody ???")
	}
}

func (r *Rules) putRule(leftNonTerminal string, body []string) {
	r.Rule = append(r.Rule, Rule{LeftSymbol: leftNonTerminal, RightSymbols: body})
}

func (r *Rules) Print() {
	fmt.Println("RULES:")
	fmt.Println("Terminal:", r.Terminal)
	fmt.Println("NonTerminal:", r.NonTerminal)
	fmt.Println("Axiom:", r.Axiom)
	fmt.Println("Rewriting Rules:")
	for _, rule := range r.Rule {
		fmt.Printf("%s -> %q\n", rule.LeftSymbol, rule.RightSymbols)
	}
}

func isFirstSetInSecond(first, second map[string]struct{}) bool {
	for key := range first {
		if _, ok := second[key]; !ok {
			return false
		}
	}
	return true
}

func convertEmptiness(rules *[]Rule) {
	for i := 0; i < len(*rules); i++ {
		if len((*rules)[i].RightSymbols) == 0 {
			(*rules)[i].RightSymbols = append((*rules)[i].RightSymbols, "")
		}
	}
}

```

### top-down-parse

`gen_table.go`
```go

package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"BODY - CloseBracket": {},
		"BODY - NonTerminal": {"NonTerminal", "BODY"},
		"BODY - Terminal": {"Terminal", "BODY"},
		"DECLARATIONS - $": {},
		"DECLARATIONS - AxiomSign": {"REWRITING_RULE", "DECLARATIONS"},
		"DECLARATIONS - NonTerminal": {"REWRITING_RULE", "DECLARATIONS"},
		"REWRITING - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_OPT - $": {},
		"REWRITING_OPT - AxiomSign": {},
		"REWRITING_OPT - NonTerminal": {},
		"REWRITING_OPT - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_RULE - AxiomSign": {"AxiomSign", "NonTerminal", "REWRITING"},
		"REWRITING_RULE - NonTerminal": {"NonTerminal", "REWRITING"},

	}
}


func newGenAxiom() string {
	return "DECLARATIONS"
}


```

`parser.go`
```go
package top_down_parse

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
)

type Parser struct {
	table map[string][]string
	axiom string
}

func NewParser() Parser {
	return Parser{
		table: newGenTable(),
		axiom: newGenAxiom(),
	}
}

func (p Parser) TopDownParse(scanner *lexer.Scanner) (*TreeNode, error) {
	type stackNode struct {
		itn *InnerTreeNode
		val string
	}
	s := NewStack[stackNode]()

	root := newTreeNode()
	root.addNode(newInnerTreeNode(""))

	s.Push(stackNode{itn: root.Root, val: p.axiom})

	t := scanner.NextToken()

	for t.GetType() != lexer.EopTag {
		topNode, err := s.Pop()
		if err != nil {
			return newTreeNode(), fmt.Errorf("failed to get top node: %w", err)
		}

		if isTerminal(topNode.val) {
			topNode.itn.Children = append(topNode.itn.Children, newLeafTreeNode(t))
			t = scanner.NextToken()
		} else if neighbourhoods,
		ok := p.table[newTableKey(topNode.val, lexer.TagToString[t.GetType()])]; ok {
			in := newInnerTreeNode(topNode.val)
			topNode.itn.Children = append(topNode.itn.Children, in)

			for i := len(neighbourhoods) - 1; i >= 0; i-- {
				s.Push(stackNode{itn: in, val: neighbourhoods[i]})
			}
		} else {
			return newTreeNode(), fmt.Errorf("failed do parse in table with val %s and token %s",
				topNode.val, t.GetValue())
		}
	}

	return root, nil
}

```

`stack.go`
```go
package top_down_parse

import "errors"

type Stack[T any] struct {
	buffer []T
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		buffer: make([]T, 0),
	}
}

func (s *Stack[T]) Push(elem T) {
	s.buffer = append(s.buffer, elem)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.buffer) > 0 {
		elem := s.buffer[len(s.buffer)-1]
		s.buffer = s.buffer[:len(s.buffer)-1]
		return elem, nil
	}
	var tmp T
	return tmp, errors.New("empty buffer")
}

func (s *Stack[T]) GetElems() []T {
	return s.buffer
}

```

`table.go`
```go
package top_down_parse

import (
	"fmt"
	"slices"
)

const (
	Declarations  = "DECLARATIONS"
	RewritingRule = "REWRITING_RULE"
	Rewriting     = "REWRITING"
	RewritingOpt  = "REWRITING_OPT"
	Body          = "BODY"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{"DECLARATIONS", "REWRITING_RULE", "REWRITING", "REWRITING_OPT", "BODY"}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}

```

`tree.go`
```go
package top_down_parse

import (
	"fmt"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
)

type TreeNode struct {
	Root *InnerTreeNode
}

func newTreeNode() *TreeNode {
	return &TreeNode{}
}

func (tn *TreeNode) Print() {
	tn.Root.printNode(0)
}

func (tn *TreeNode) addNode(node *InnerTreeNode) {
	tn.Root = node
}

type TreeNodePrinter interface {
	printNode(offset int)
}

type InnerTreeNode struct {
	NonTerminal string
	Children    []TreeNodePrinter
}

func newInnerTreeNode(nonTerminal string) *InnerTreeNode {
	return &InnerTreeNode{NonTerminal: nonTerminal, Children: make([]TreeNodePrinter, 0)}
}

func (itn InnerTreeNode) printNode(offset int) {
	fmt.Printf(strings.Repeat("..", offset) + fmt.Sprintf("Inner node: %s\n", itn.NonTerminal))

	for _, child := range itn.Children {
		child.printNode(offset + 1)
	}
}

type LeafTreeNode struct {
	Token lexer.IToken
}

func newLeafTreeNode(token lexer.IToken) *LeafTreeNode {
	return &LeafTreeNode{Token: token}
}

func (ltn LeafTreeNode) printNode(offset int) {
	if ltn.Token.GetType() == lexer.TermTag || ltn.Token.GetType() == lexer.NonTermTag {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s - %s\n", lexer.TagToString[ltn.Token.GetType()], ltn.Token.GetValue()))
	} else {
		fmt.Printf(strings.Repeat("..", offset) +
			fmt.Sprintf("Leaf: %s\n", lexer.TagToString[ltn.Token.GetType()]))
	}
}

```


### main

`main.go`
```go
package main

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/predtable"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/semantic"
	"github.com/VyacheslavIsWorkingNow/cd/lab10/converter/top_down_parse"
)

const filepath = "test_files/grammarForGrammar.txt"

const (
	tablePath = "top_down_parse/gen_table.go"
)

func main() {

	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, errTDP := parser.TopDownParse(scanner)
	if errTDP != nil {
		log.Fatalf("err in TopDownParse %+v", errTDP)
	}

	tree.Print()

	scanner.GetCompiler().PrintMessages()

	sem := semantic.NewSemantic(tree)

	rules, errSem := sem.StartSemanticAnalysis()
	if errSem != nil {
		log.Fatalf("err in semantic %+v", errSem)
	}

	rules.Print()

	genTable, errT := predtable.GenTable(rules)
	if errT != nil {
		log.Fatalf("err in gen table: %+v", errSem)
	}

	predtable.PrintGenTable(genTable)

	errUF := predtable.UploadTableToFile(tablePath, genTable, rules)
	if errUF != nil {
		log.Fatalf("err in upload table: %+v", errUF)
	}

	fmt.Println("finish")
}

```

`go.mod`
```go
module github.com/VyacheslavIsWorkingNow/cd/lab10/converter

go 1.22

```


# Тестирование

## Модуль конвертора
Входные данные

Два теста:

Для грамматики самой грамматики 

```
* DECLARATIONS (REWRITING_RULE DECLARATIONS) () 

REWRITING_RULE ("AxiomSign" "NonTerminal" REWRITING) ("NonTerminal" REWRITING)

REWRITING ("OpenBracket" BODY "CloseBracket" REWRITING_OPT)

REWRITING_OPT ("OpenBracket" BODY "CloseBracket" REWRITING_OPT) ()

BODY ("Terminal" BODY) ("NonTerminal" BODY) ()
```

На stdout идут только логи, таблица сохраняется в `gen_table.go`
Логи:
```
Inner node: 
..Inner node: DECLARATIONS
....Inner node: REWRITING_RULE
......Leaf: AxiomSign
......Leaf: NonTerminal - DECLARATIONS
......Inner node: REWRITING
........Leaf: OpenBracket
........Inner node: BODY
..........Leaf: NonTerminal - REWRITING_RULE
..........Inner node: BODY
............Leaf: NonTerminal - DECLARATIONS
............Inner node: BODY
........Leaf: CloseBracket
........Inner node: REWRITING_OPT
..........Leaf: OpenBracket
..........Inner node: BODY
..........Leaf: CloseBracket
..........Inner node: REWRITING_OPT
....Inner node: DECLARATIONS
......Inner node: REWRITING_RULE
........Leaf: NonTerminal - REWRITING_RULE
........Inner node: REWRITING
..........Leaf: OpenBracket
..........Inner node: BODY
............Leaf: Terminal - AxiomSign
............Inner node: BODY
..............Leaf: Terminal - NonTerminal
..............Inner node: BODY
................Leaf: NonTerminal - REWRITING
................Inner node: BODY
..........Leaf: CloseBracket
..........Inner node: REWRITING_OPT
............Leaf: OpenBracket
............Inner node: BODY
..............Leaf: Terminal - NonTerminal
..............Inner node: BODY
................Leaf: NonTerminal - REWRITING
................Inner node: BODY
............Leaf: CloseBracket
............Inner node: REWRITING_OPT
......Inner node: DECLARATIONS
........Inner node: REWRITING_RULE
..........Leaf: NonTerminal - REWRITING
..........Inner node: REWRITING
............Leaf: OpenBracket
............Inner node: BODY
..............Leaf: Terminal - OpenBracket
..............Inner node: BODY
................Leaf: NonTerminal - BODY
................Inner node: BODY
..................Leaf: Terminal - CloseBracket
..................Inner node: BODY
....................Leaf: NonTerminal - REWRITING_OPT
....................Inner node: BODY
............Leaf: CloseBracket
............Inner node: REWRITING_OPT
........Inner node: DECLARATIONS
..........Inner node: REWRITING_RULE
............Leaf: NonTerminal - REWRITING_OPT
............Inner node: REWRITING
..............Leaf: OpenBracket
..............Inner node: BODY
................Leaf: Terminal - OpenBracket
................Inner node: BODY
..................Leaf: NonTerminal - BODY
..................Inner node: BODY
....................Leaf: Terminal - CloseBracket
....................Inner node: BODY
......................Leaf: NonTerminal - REWRITING_OPT
......................Inner node: BODY
..............Leaf: CloseBracket
..............Inner node: REWRITING_OPT
................Leaf: OpenBracket
................Inner node: BODY
................Leaf: CloseBracket
................Inner node: REWRITING_OPT
..........Inner node: DECLARATIONS
............Inner node: REWRITING_RULE
..............Leaf: NonTerminal - BODY
..............Inner node: REWRITING
................Leaf: OpenBracket
................Inner node: BODY
..................Leaf: Terminal - Terminal
..................Inner node: BODY
....................Leaf: NonTerminal - BODY
....................Inner node: BODY
................Leaf: CloseBracket
................Inner node: REWRITING_OPT
..................Leaf: OpenBracket
..................Inner node: BODY
....................Leaf: Terminal - NonTerminal
....................Inner node: BODY
......................Leaf: NonTerminal - BODY
......................Inner node: BODY
..................Leaf: CloseBracket
..................Inner node: REWRITING_OPT
....................Leaf: OpenBracket
....................Inner node: BODY
....................Leaf: CloseBracket
_____MESSAGES_____
2024/05/16 16:11:27 start semantic analysis
2024/05/16 16:11:27 start convert
RULES:
Terminal: map[AxiomSign:{} CloseBracket:{} NonTerminal:{} OpenBracket:{} Terminal:{}]
NonTerminal: map[BODY:{} DECLARATIONS:{} REWRITING:{} REWRITING_OPT:{} REWRITING_RULE:{}]
Axiom: DECLARATIONS
Rewriting Rules:
DECLARATIONS -> ["REWRITING_RULE" "DECLARATIONS"]
DECLARATIONS -> [""]
REWRITING_RULE -> ["AxiomSign" "NonTerminal" "REWRITING"]
REWRITING_RULE -> ["NonTerminal" "REWRITING"]
REWRITING -> ["OpenBracket" "BODY" "CloseBracket" "REWRITING_OPT"]
REWRITING_OPT -> ["OpenBracket" "BODY" "CloseBracket" "REWRITING_OPT"]
REWRITING_OPT -> [""]
BODY -> ["Terminal" "BODY"]
BODY -> ["NonTerminal" "BODY"]
BODY -> [""]
Gen Table:
REWRITING_OPT - NonTerminal -> []
REWRITING_OPT - AxiomSign -> []
REWRITING_OPT - $ -> []
BODY - CloseBracket -> []
DECLARATIONS - $ -> []
REWRITING_RULE - AxiomSign -> [AxiomSign NonTerminal REWRITING]
REWRITING_RULE - NonTerminal -> [NonTerminal REWRITING]
REWRITING - OpenBracket -> [OpenBracket BODY CloseBracket REWRITING_OPT]
BODY - NonTerminal -> [NonTerminal BODY]
DECLARATIONS - AxiomSign -> [REWRITING_RULE DECLARATIONS]
DECLARATIONS - NonTerminal -> [REWRITING_RULE DECLARATIONS]
REWRITING_OPT - OpenBracket -> [OpenBracket BODY CloseBracket REWRITING_OPT]
BODY - Terminal -> [Terminal BODY]
finish

```

Таблица:
```go

package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"BODY - CloseBracket": {},
		"BODY - NonTerminal": {"NonTerminal", "BODY"},
		"BODY - Terminal": {"Terminal", "BODY"},
		"DECLARATIONS - $": {},
		"DECLARATIONS - AxiomSign": {"REWRITING_RULE", "DECLARATIONS"},
		"DECLARATIONS - NonTerminal": {"REWRITING_RULE", "DECLARATIONS"},
		"REWRITING - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_OPT - $": {},
		"REWRITING_OPT - AxiomSign": {},
		"REWRITING_OPT - NonTerminal": {},
		"REWRITING_OPT - OpenBracket": {"OpenBracket", "BODY", "CloseBracket", "REWRITING_OPT"},
		"REWRITING_RULE - AxiomSign": {"AxiomSign", "NonTerminal", "REWRITING"},
		"REWRITING_RULE - NonTerminal": {"NonTerminal", "REWRITING"},

	}
}


func newGenAxiom() string {
	return "DECLARATIONS"
}

```

Для грамматики в задании

Логи: 
```
Inner node: 
..Inner node: DECLARATIONS
....Inner node: REWRITING_RULE
......Leaf: NonTerminal - F
......Inner node: REWRITING
........Leaf: OpenBracket
........Inner node: BODY
..........Leaf: Terminal - n
..........Inner node: BODY
........Leaf: CloseBracket
........Inner node: REWRITING_OPT
..........Leaf: OpenBracket
..........Inner node: BODY
............Leaf: Terminal - (
............Inner node: BODY
..............Leaf: NonTerminal - E
..............Inner node: BODY
................Leaf: Terminal - )
................Inner node: BODY
..........Leaf: CloseBracket
..........Inner node: REWRITING_OPT
....Inner node: DECLARATIONS
......Inner node: REWRITING_RULE
........Leaf: NonTerminal - T
........Inner node: REWRITING
..........Leaf: OpenBracket
..........Inner node: BODY
............Leaf: NonTerminal - F
............Inner node: BODY
..............Leaf: NonTerminal - T'
..............Inner node: BODY
..........Leaf: CloseBracket
..........Inner node: REWRITING_OPT
......Inner node: DECLARATIONS
........Inner node: REWRITING_RULE
..........Leaf: NonTerminal - T'
..........Inner node: REWRITING
............Leaf: OpenBracket
............Inner node: BODY
..............Leaf: Terminal - *
..............Inner node: BODY
................Leaf: NonTerminal - F
................Inner node: BODY
..................Leaf: NonTerminal - T'
..................Inner node: BODY
............Leaf: CloseBracket
............Inner node: REWRITING_OPT
..............Leaf: OpenBracket
..............Inner node: BODY
..............Leaf: CloseBracket
..............Inner node: REWRITING_OPT
........Inner node: DECLARATIONS
..........Inner node: REWRITING_RULE
............Leaf: AxiomSign
............Leaf: NonTerminal - E
............Inner node: REWRITING
..............Leaf: OpenBracket
..............Inner node: BODY
................Leaf: NonTerminal - T
................Inner node: BODY
..................Leaf: NonTerminal - E'
..................Inner node: BODY
..............Leaf: CloseBracket
..............Inner node: REWRITING_OPT
..........Inner node: DECLARATIONS
............Inner node: REWRITING_RULE
..............Leaf: NonTerminal - E'
..............Inner node: REWRITING
................Leaf: OpenBracket
................Inner node: BODY
..................Leaf: Terminal - +
..................Inner node: BODY
....................Leaf: NonTerminal - T
....................Inner node: BODY
......................Leaf: NonTerminal - E'
......................Inner node: BODY
................Leaf: CloseBracket
................Inner node: REWRITING_OPT
..................Leaf: OpenBracket
..................Inner node: BODY
..................Leaf: CloseBracket
_____MESSAGES_____
Type: Comment | i: 0 | position: (1,1)-(1,34) | text: /* аксиома помечена звёздочкой */
2024/05/16 16:14:15 start semantic analysis
2024/05/16 16:14:15 start convert
RULES:
Terminal: map[(:{} ):{} *:{} +:{} n:{}]
NonTerminal: map[E:{} E':{} F:{} T:{} T':{}]
Axiom: E
Rewriting Rules:
F -> ["n"]
F -> ["(" "E" ")"]
T -> ["F" "T'"]
T' -> ["*" "F" "T'"]
T' -> [""]
E -> ["T" "E'"]
E' -> ["+" "T" "E'"]
E' -> [""]
Gen Table:
T' - + -> []
E - n -> [T E']
E' - $ -> []
F - ( -> [( E )]
T - n -> [F T']
T - ( -> [F T']
T' - * -> [* F T']
T' - ) -> []
F - n -> [n]
T' - $ -> []
E - ( -> [T E']
E' - + -> [+ T E']
E' - ) -> []
finish

```

Таблица:
```go

package top_down_parse

func newGenTable() map[string][]string {
	return map[string][]string{
		"E - (": {"T", "E'"},
		"E - n": {"T", "E'"},
		"E' - $": {},
		"E' - )": {},
		"E' - +": {"+", "T", "E'"},
		"F - (": {"(", "E", ")"},
		"F - n": {"n"},
		"T - (": {"F", "T'"},
		"T - n": {"F", "T'"},
		"T' - $": {},
		"T' - )": {},
		"T' - *": {"*", "F", "T'"},
		"T' - +": {},

	}
}


func newGenAxiom() string {
	return "E"
}
```

## Модуль калькулятора
Входные данные

`2 * (3 + 4) + 4 + (4 * 2)`

Вывод на `stdout`
```
TOKENS:
n (2,1)-(2,2): 2
* (2,3)-(2,3): *
( (2,5)-(2,5): (
n (2,6)-(2,7): 3
+ (2,8)-(2,8): +
n (2,10)-(2,11): 4
) (2,11)-(2,11): )
+ (2,13)-(2,13): +
n (2,15)-(2,16): 4
+ (2,17)-(2,17): +
( (2,19)-(2,19): (
n (2,20)-(2,21): 4
* (2,22)-(2,22): *
n (2,24)-(2,25): 2
) (2,25)-(2,25): )
Eop (0,0)-(0,0): end of file
Eop (0,0)-(0,0): end of file
Inner node: 
..Inner node: E
....Inner node: T
......Inner node: F
........Leaf: n - 2
......Inner node: T'
........Leaf: *
........Inner node: F
..........Leaf: (
..........Inner node: E
............Inner node: T
..............Inner node: F
................Leaf: n - 3
..............Inner node: T'
............Inner node: E'
..............Leaf: +
..............Inner node: T
................Inner node: F
..................Leaf: n - 4
................Inner node: T'
..............Inner node: E'
..........Leaf: )
........Inner node: T'
....Inner node: E'
......Leaf: +
......Inner node: T
........Inner node: F
..........Leaf: n - 4
........Inner node: T'
......Inner node: E'
........Leaf: +
........Inner node: T
..........Inner node: F
............Leaf: (
............Inner node: E
..............Inner node: T
................Inner node: F
..................Leaf: n - 4
................Inner node: T'
..................Leaf: *
..................Inner node: F
....................Leaf: n - 2
..................Inner node: T'
..............Inner node: E'
............Leaf: )
2024/05/16 16:16:02 start convert
solver: 26
2024/05/16 16:16:02 end calculator

```

# Вывод
В данной лабораторной работе я изучил алгоритм построения таблиц предсказывающего анализатора
и проработал каждую его часть в коде: first, follow множества, построение таблицы.
Так же я обобщил код top-down-parse, ему на вход подается таблица. То есть я получил 
мощный инструмент для генерации любых таблиц, нужно лишь задать грамматику в нужном виде, 
это меня впечатлило в данной работе больше всего. Так же я проработал два подхода к 
обходу дерева вывода (в конверторе и калькуляторе). Так же меня сильно удивило то,
что лексический анализатор для калькулятора был сделан менее, чем за полчаса, 
виден большой прогресс с начала семестра, когда это занимало больше дня. 
В итоге эта лабороторная однозначно одна из самых больших за все время обучения,
и она объеденила в себе написание лексического анализатора, работу с синтаксическим деревом
и семантический анализ. И нужно было поломать голову над алгоритмами, чтобы это все 
корректно написать. Спасибо за очень крутую лабораторную!


