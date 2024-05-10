% Лабораторная работа № 2.3 «Синтаксический анализатор на основе
предсказывающего анализа»
% 2 мая 2024 г.
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

## Неформальное описание синтаксиса входного языка

На вход подается список правил переписывания.
В начале правила переписывания может быть указано, является оно аксиомой или нет.
Далее указан нетерминал и то, во что он может быть переписан.
То, во что он может быть переписан указано в круглых скобках.
Внутри находится список терминалов, нетерминалов или пустота.

## Лексическая структура
```
AxiomSign ::= *
OpenBracket ::= (
CloseBracket ::= )
Whitespace ::= [ \t\n\r]+
Comment ::= /*([^(*/)]*)*/
NonTerminal ::= [a-zA-Z][a-zA-Z0-9]*(')?
Terminal    ::= "[^"]+"
```

## Грамматика языка
`DECLARATIONS ::= REWRITING_RULE DECLARATIONS | epsilon `

`REWRITING_RULE ::= AxiomSign NonTerminal REWRITING | NonTerminal REWRITING`

`REWRITING ::= OpenBracket BODY CloseBracket REWRITING_OPT`

`REWRITING_OPT ::= OpenBracket BODY CloseBracket REWRITING_OPT | epsilon`

`BODY ::= Terminal BODY | NonTerminal BODY | epsilon`

## Программная реализация

### пакет lexer

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
		token.GetType() == CloseBracketTag || token.GetType() == OpenBracketTag ||
		token.GetType() == AxiomTag {
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
		} else if runeScanner.IsOpenBracket() || runeScanner.IsCloseBracket() || runeScanner.IsStar() {
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

	return NewTerminal(string(currentString), NewFragment(start, curPosition))
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
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
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

### пакет top_down_parse

`parser.go`
```go
package top_down_parse

import (
	"fmt"

	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
)

type Parser struct {
	table map[string][]string
}

func NewParser() Parser {
	return Parser{
		table: newTable(),
	}
}

func (p Parser) TopDownParse(scanner *lexer.Scanner) (*TreeNode, error) {
	type stackNode struct {
		itn *innerTreeNode
		val string
	}
	s := NewStack[stackNode]()

	root := newTreeNode()
	root.addNode(newInnerTreeNode(""))
	s.Push(stackNode{itn: root.root, val: declarations})

	t := scanner.NextToken()

	for t.GetType() != lexer.EopTag {
		topNode, err := s.Pop()
		if err != nil {
			return newTreeNode(), fmt.Errorf("failed to get top node: %w", err)
		}

		if isTerminal(topNode.val) {
			topNode.itn.children = append(topNode.itn.children, newLeafTreeNode(t))
			t = scanner.NextToken()
		} else if neighbourhoods, ok := p.table[newTableKey(topNode.val, lexer.TagToString[t.GetType()])];
		ok {
			in := newInnerTreeNode(topNode.val)
			topNode.itn.children = append(topNode.itn.children, in)

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
	declarations  = "DECLARATIONS"
	rewritingRule = "REWRITING_RULE"
	rewriting     = "REWRITING"
	rewritingOpt  = "REWRITING_OPT"
	body          = "BODY"
	nonTerminal   = "NonTerminal"
	terminal      = "Terminal"
	axiomSign     = "AxiomSign"
	openBracket   = "OpenBracket"
	closeBracket  = "CloseBracket"
	eof           = "Eof"
)

func newTableKey(nonTerminal, terminal string) string {
	return fmt.Sprintf("%s - %s", nonTerminal, terminal)
}

func terminalValue() []string {
	return []string{declarations, rewritingRule, rewriting, rewritingOpt, body}
}

func isTerminal(s string) bool {
	return !slices.Contains(terminalValue(), s)
}

func newTable() map[string][]string {
	return map[string][]string{
		newTableKey(declarations, axiomSign):   {rewritingRule, declarations},
		newTableKey(declarations, nonTerminal): {rewritingRule, declarations},
		newTableKey(declarations, eof):         {},

		newTableKey(rewritingRule, axiomSign):   {axiomSign, nonTerminal, rewriting},
		newTableKey(rewritingRule, nonTerminal): {nonTerminal, rewriting},

		newTableKey(rewriting, openBracket): {openBracket, body, closeBracket, rewritingOpt},

		newTableKey(rewritingOpt, axiomSign):   {},
		newTableKey(rewritingOpt, nonTerminal): {},
		newTableKey(rewritingOpt, openBracket): {openBracket, body, closeBracket, rewritingOpt},
		newTableKey(rewritingOpt, eof):         {},

		newTableKey(body, nonTerminal):  {nonTerminal, body},
		newTableKey(body, terminal):     {terminal, body},
		newTableKey(body, closeBracket): {},
	}
}

```

`tree.go`
```go
package top_down_parse

import (
	"fmt"
	"strings"

	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
)

type TreeNode struct {
	root *innerTreeNode
}

func newTreeNode() *TreeNode {
	return &TreeNode{}
}

func (tn *TreeNode) Print() {
	tn.root.printNode(0)
}

func (tn *TreeNode) addNode(node *innerTreeNode) {
	tn.root = node
}

type treeNodePrinter interface {
	printNode(offset int)
}

type innerTreeNode struct {
	nonTerminal string
	children    []treeNodePrinter
}

func newInnerTreeNode(nonTerminal string) *innerTreeNode {
	return &innerTreeNode{nonTerminal: nonTerminal, children: make([]treeNodePrinter, 0)}
}

func (itn innerTreeNode) printNode(offset int) {
	fmt.Printf(strings.Repeat("\t", offset) + fmt.Sprintf("Inner node: %s\n", itn.nonTerminal))

	for _, child := range itn.children {
		child.printNode(offset + 1)
	}
}

type leafTreeNode struct {
	token lexer.IToken
}

func newLeafTreeNode(token lexer.IToken) *leafTreeNode {
	return &leafTreeNode{token: token}
}

func (ltn *leafTreeNode) printNode(offset int) {
	if ltn.token.GetType() == lexer.TermTag || ltn.token.GetType() == lexer.NonTermTag {
		fmt.Printf(strings.Repeat("\t", offset) +
			fmt.Sprintf("Leaf: %s - %s\n", lexer.TagToString[ltn.token.GetType()], ltn.token.GetValue()))
	} else {
		fmt.Printf(strings.Repeat("\t", offset) +
			fmt.Sprintf("Leaf: %s\n", lexer.TagToString[ltn.token.GetType()]))
	}
}

```

### main файл

`main.go`
```go
package main

import (
	"fmt"
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab8/lexer"
	"github.com/VyacheslavIsWorkingNow/cd/lab8/top_down_parse"
)

const filepath = "test_files/mixed.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	parser := top_down_parse.NewParser()

	tree, err := parser.TopDownParse(scanner)
	if err != nil {
		log.Panic("пупу:", err)
	}

	tree.Print()

	scanner.GetCompiler().PrintMessages()

	fmt.Println("finish")
}

```

# Тестирование

Входные данные

```
/* аксиома помечена звёздочкой */
  F  ("n") ("(" E ")")
  T  (F T')
  T' ("*" F T') ()
* E  (T E')
  E' ("+" T E') ()
```

Вывод на `stdout`

<!-- ENABLE LONG LINES -->

```
Inner node: 
        Inner node: DECLARATIONS
                Inner node: REWRITING_RULE
                        Leaf: NonTerminal - F
                        Inner node: REWRITING
                                Leaf: OpenBracket
                                Inner node: BODY
                                        Leaf: Terminal - "n"
                                        Inner node: BODY
                                Leaf: CloseBracket
                                Inner node: REWRITING_OPT
                                        Leaf: OpenBracket
                                        Inner node: BODY
                                                Leaf: Terminal - "("
                                                Inner node: BODY
                                                        Leaf: NonTerminal - E
                                                        Inner node: BODY
                                                                Leaf: Terminal - ")"
                                                                Inner node: BODY
                                        Leaf: CloseBracket
                                        Inner node: REWRITING_OPT
                Inner node: DECLARATIONS
                        Inner node: REWRITING_RULE
                                Leaf: NonTerminal - T
                                Inner node: REWRITING
                                        Leaf: OpenBracket
                                        Inner node: BODY
                                                Leaf: NonTerminal - F
                                                Inner node: BODY
                                                        Leaf: NonTerminal - T'
                                                        Inner node: BODY
                                        Leaf: CloseBracket
                                        Inner node: REWRITING_OPT
                        Inner node: DECLARATIONS
                                Inner node: REWRITING_RULE
                                        Leaf: NonTerminal - T'
                                        Inner node: REWRITING
                                                Leaf: OpenBracket
                                                Inner node: BODY
                                                        Leaf: Terminal - "*"
                                                        Inner node: BODY
                                                                Leaf: NonTerminal - F
                                                                Inner node: BODY
                                                                        Leaf: NonTerminal - T'
                                                                        Inner node: BODY
                                                Leaf: CloseBracket
                                                Inner node: REWRITING_OPT
                                                        Leaf: OpenBracket
                                                        Inner node: BODY
                                                        Leaf: CloseBracket
                                                        Inner node: REWRITING_OPT
                                Inner node: DECLARATIONS
                                        Inner node: REWRITING_RULE
                                                Leaf: AxiomSign
                                                Leaf: NonTerminal - E
                                                Inner node: REWRITING
                                                        Leaf: OpenBracket
                                                        Inner node: BODY
                                                                Leaf: NonTerminal - T
                                                                Inner node: BODY
                                                                        Leaf: NonTerminal - E'
                                                                        Inner node: BODY
                                                        Leaf: CloseBracket
                                                        Inner node: REWRITING_OPT
                                        Inner node: DECLARATIONS
                                                Inner node: REWRITING_RULE
                                                        Leaf: NonTerminal - E'
                                                        Inner node: REWRITING
                                                                Leaf: OpenBracket
                                                                Inner node: BODY
                                                                        Leaf: Terminal - "+"
                                                                        Inner node: BODY
                                                                                Leaf: NonTerminal - T
                                                                                Inner node: BODY
                                                                                        Leaf: NonTerminal - E'
                                                                                        Inner node: BODY
                                                                Leaf: CloseBracket
                                                                Inner node: REWRITING_OPT
                                                                        Leaf: OpenBracket
                                                                        Inner node: BODY
                                                                        Leaf: CloseBracket
_____MESSAGES_____
Type: Comment | i: 0 | position: (1,1)-(1,34) | text: /* аксиома помечена звёздочкой */
finish

```

# Вывод
В ходе выполнения лабораторной работы я разобрался с тем,
как строить таблицы предсказывающего анализатора.
Так же узнал об алгоритме top-down-parse.
После полученных из лабораторной 2.2 навыков построить синтаксис языка
оказалось не таким сложным заданием.
Как лексический анализатор решил выбрать свой ООП лексер.
Так же я решил воспользоваться возможностями языка Go:
реализовал стек на дженериках и узлы дерева представил как интерфейс
с методом printNode(offset int), для каждого узла дерева вывода реализовал этот интерфейс.
Интересно было реализовывать синтаксис для языка представления правил грамматики,
механизм предсказывающего анализатора мне понравился. 
