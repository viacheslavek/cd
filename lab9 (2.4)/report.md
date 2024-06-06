Лабораторная работа № 2.4 «Рекурсивный спуск»
% 5 июня 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является изучение алгоритмов построения
парсеров методом рекурсивного спуска.

# Индивидуальный вариант
Определения структур, объединений и перечислений языка Си.
В инициализаторах перечислений допустимы знаки операций +, -, *, /,
sizeof, операндами могут служить имена перечислимых значений и целые числа.
Числовые константы могут быть только целочисленными и десятичными.

# Реализация

## Лексическая структура
```
IDENTIFIER = [a-zA-z][a-zA-Z0-9_]*
INTEGER = [0-9]+
```

## Грамматика языка
```
Program ::= DeclarationList
DeclarationList ::= Declaration*
Declaration ::= TypeSpecifier AbstractDeclaratorsOpt ';' 


AbstractDeclaratorsOpt ::= AbstractDeclarators?

AbstractDeclarators ::= AbstractDeclarator (',' AbstractDeclarator)* 

AbstractDeclarator ::= AbstractDeclaratorPointer | AbstractDeclaratorArrayList 

AbstractDeclaratorPointer ::= '*' AbstractDeclarator 
AbstractDeclaratorArrayList ::= AbstractDeclaratorArray+

AbstractDeclaratorArray ::= AbstractDeclaratorPrimArray |
    AbstractDeclaratorPrimSimple | AbstractDeclaratorPrimDifficult

AbstractDeclaratorPrimArray ::= '[' Expression ']'
AbstractDeclaratorPrimSimple ::= IDENTIFIER
AbstractDeclaratorPrimDifficult ::= '(' AbstractDeclarator ')'


TypeSpecifier ::= SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier

SimpleTypeSpecifier ::= char | short | int | long | float | double 


EnumTypeSpecifier ::= ENUM EnumStatement 

EnumStatement ::= IdentEnumStatement | BodyEnumStatement
IdentEnumStatement ::= IDENTIFIER BodyEnumStatement?
BodyEnumStatement ::= '{' EnumeratorList '}'

EnumeratorList ::= Enumerator (',' Enumerator)*  
Enumerator ::= IDENTIFIER ('=' Expression)? 


StructOrUnionSpecifier ::= (struct | union) StructOrUnionStatement 

StructOrUnionStatement ::= IdentStructOrUnionStatement | BodyStructOrUnionStatement
IdentStructOrUnionStatement ::= IDENTIFIER BodyStructOrUnionStatement?
BodyStructOrUnionStatement ::= '{' DeclarationList '}' 


Expression ::= Term (('+' | '-') Term)*  

Term ::= Factor (('*' | '/') Factor)*

Factor ::= sizeof '(' (struct | union | enum) IDENTIFIER ')' 
    | IDENTIFIER | INTEGER | '(' Expression ')'

```

## Программная реализация

### Пакет лексического анализатора

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

func (p *RunePosition) IsLineTranslation() bool {
	return p.GetRune() == '\n'
}

func (p *RunePosition) IsLetter() bool {
	return unicode.IsLetter(p.GetRune())
}

func (p *RunePosition) IsLatinLetter() bool {
	r := unicode.ToLower(p.GetRune())
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func (p *RunePosition) IsUnderlining() bool {
	return p.GetRune() == '_'
}

func (p *RunePosition) IsDigit() bool {
	return unicode.IsDigit(p.GetRune())
}

func (p *RunePosition) IsOpenRoundBracket() bool {
	return p.GetRune() == '('
}

func (p *RunePosition) IsCloseRoundBracket() bool {
	return p.GetRune() == ')'
}

func (p *RunePosition) IsOpenSquareBracket() bool {
	return p.GetRune() == '['
}

func (p *RunePosition) IsCloseSquareBracket() bool {
	return p.GetRune() == ']'
}

func (p *RunePosition) IsOpenCurlyBracket() bool {
	return p.GetRune() == '{'
}

func (p *RunePosition) IsCloseCurlyBracket() bool {
	return p.GetRune() == '}'
}

func (p *RunePosition) IsComma() bool {
	return p.GetRune() == ','
}

func (p *RunePosition) IsSemicolon() bool {
	return p.GetRune() == ';'
}

func (p *RunePosition) IsStar() bool {
	return p.GetRune() == '*'
}

func (p *RunePosition) IsPlus() bool {
	return p.GetRune() == '+'
}

func (p *RunePosition) IsMinus() bool {
	return p.GetRune() == '-'
}

func (p *RunePosition) IsSlash() bool {
	return p.GetRune() == '/'
}

func (p *RunePosition) IsEqual() bool {
	return p.GetRune() == '='
}

func (p *RunePosition) IsSpecialSymbol() bool {
	return p.IsOpenRoundBracket() || p.IsCloseRoundBracket() ||
		p.IsOpenSquareBracket() || p.IsCloseSquareBracket() ||
		p.IsOpenCurlyBracket() || p.IsCloseCurlyBracket() ||
		p.IsMinus() || p.IsPlus() || p.IsSlash() || p.IsStar() ||
		p.IsComma() || p.IsSemicolon() || p.IsEqual()
}

var keywords = map[string]struct{}{
	"enum": {}, "struct": {}, "union": {},
	"sizeof": {},
	"char":   {}, "short": {}, "int": {}, "long": {}, "float": {}, "double": {},
}

func IsKeyword(ident string) bool {
	_, ok := keywords[ident]
	return ok
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
			tokens = append(tokens, processInt(runeScanner))
		} else if runeScanner.IsSpecialSymbol() {
			tokens = append(tokens, processSpecialSymbol(runeScanner))
		} else if runeScanner.IsLetter() && runeScanner.IsLatinLetter() {
			tokens = append(tokens, processIdentifier(runeScanner))
		} else {
			if runeScanner.GetRune() == -1 {
				tokens = append(tokens, NewEOP())
			} else {
				fmt.Println("rune:", string(runeScanner.GetRune()))
				log.Fatalf("incorrect rune in parser %+v", runeScanner.GetCurrentPosition())
			}
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	tokens = append(tokens, NewEOP())

	return tokens
}

func processInt(rs *RunePosition) IToken {
	currentInt := make([]rune, 0)
	start := rs.GetCurrentPosition()

	for rs.IsDigit() {
		currentInt = append(currentInt, rs.GetRune())
		rs.NextRune()
	}

	curPosition := rs.GetCurrentPosition()

	return NewInteger(string(currentInt), NewFragment(start, curPosition))
}

func processSpecialSymbol(rs *RunePosition) IToken {
	start := rs.GetCurrentPosition()
	operand := rs.GetRune()
	rs.NextRune()
	curPosition := rs.GetCurrentPosition()

	return NewSpecSymbol(string(operand), NewFragment(start, curPosition))
}

func processIdentifier(rs *RunePosition) IToken {
	currentIdentifier := make([]rune, 0)

	start := rs.GetCurrentPosition()
	curPositionToken := rs.GetCurrentPosition()
	currentIdentifier = append(currentIdentifier, rs.GetRune())
	rs.NextRune()

	for !rs.IsWhiteSpace() && !rs.IsSpecialSymbol() {
		if rs.GetRune() == -1 {
			return Token{}
		}
		if (rs.IsLetter() && rs.IsLatinLetter()) || rs.IsDigit() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else if rs.IsUnderlining() {
			currentIdentifier = append(currentIdentifier, rs.GetRune())
		} else {
			log.Fatalf("error in process identifier")
		}
		curPositionToken = rs.GetCurrentPosition()
		rs.NextRune()
	}

	if IsKeyword(string(currentIdentifier)) {
		return NewKeyword(string(currentIdentifier), NewFragment(start, curPositionToken))
	}

	return NewIdentifier(string(currentIdentifier), NewFragment(start, curPositionToken))
}

```

`tag.go`
```go
package lexer

type DomainTag int

const (
	IdentifierTag DomainTag = iota + 1
	IntTag
	KeywordTag
	SpecSymbolTag
	EopTag
)

var TagToString = map[DomainTag]string{
	IdentifierTag: "Identifier",
	IntTag:        "Integer",
	KeywordTag:    "Keyword",
	SpecSymbolTag: "SpecialSymbol",
	EopTag:        "Eop",
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

type IdentifierToken struct {
	Token
}

func NewIdentifier(value string, fragment Fragment) IdentifierToken {
	return IdentifierToken{
		Token: NewToken(IdentifierTag, value, fragment),
	}
}

func (ntt IdentifierToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[ntt.Type], ntt.Coordinate, ntt.Value)
}

type IntegerToken struct {
	Token
}

func NewInteger(value string, fragment Fragment) IntegerToken {
	return IntegerToken{
		Token: NewToken(IntTag, value, fragment),
	}
}

func (tt IntegerToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[tt.Type], tt.Coordinate, tt.Value)
}

type SpecSymbolToken struct {
	Token
}

func NewSpecSymbol(value string, fragment Fragment) SpecSymbolToken {
	return SpecSymbolToken{
		Token: NewToken(SpecSymbolTag, value, fragment),
	}
}

func (st SpecSymbolToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

type KeywordToken struct {
	Token
}

func NewKeyword(value string, fragment Fragment) KeywordToken {
	return KeywordToken{
		Token: NewToken(KeywordTag, value, fragment),
	}
}

func (st KeywordToken) String() string {
	return fmt.Sprintf("%s %s: %s", TagToString[st.Type], st.Coordinate, st.Value)
}

```

### Пакет синтаксического спуска

`model.go`
```go
package recdesc

import (
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

type RecursiveParser struct {
	currentToken lexer.IToken
	scanner      *lexer.Scanner
}

func NewParser(scanner *lexer.Scanner) *RecursiveParser {
	return &RecursiveParser{currentToken: scanner.NextToken(), scanner: scanner}
}

func (rp *RecursiveParser) isExpectedToken(tokenValue string, tokenType lexer.DomainTag) {
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	if t.GetValue() != tokenValue {
		log.Fatalf("incorrect token value. Expected: '%s' given: '%s'",
			tokenValue, t.GetValue())
	}
	if t.GetType() != tokenType {
		log.Fatalf("incorrect token type. Expected: '%s' given: '%s'",
			lexer.TagToString[tokenType], lexer.TagToString[t.GetType()])
	}
}

```

`printer.go`
```go
package recdesc

import (
	"fmt"
	"strings"
)

const offsetString = ".."

func (p Program) Print() {
	p.printNode(0)
}

func (p Program) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "Program:")
	p.declarationList.printNode(offset + 1)
}

func (dl DeclarationList) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "DeclarationList:")
	for _, d := range dl.declarations {
		d.printNode(offset + 1)
	}
}

func (d Declaration) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "Declaration:")
	d.typeSpecifier.printNode(offset + 1)
	d.abstractDeclaratorsOpt.printNode(offset + 1)
}

func (ado AbstractDeclaratorsOpt) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclaratorsOpt:")
	ado.abstractDeclarators.printNode(offset + 1)
}

func (ads AbstractDeclarators) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclarators:")
	for _, ad := range ads.abstractDeclarators {
		ad.printNode(offset + 1)
	}
}

func (adal AbstractDeclaratorArrayList) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclaratorArrayList:")
	for _, ada := range adal.abstractDeclaratorArray {
		ada.printNode(offset + 1)
	}
}

func (adpa AbstractDeclaratorPrimArray) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclaratorPrimArray:")
	adpa.expression.printNode(offset + 1)

}

func (adps AbstractDeclaratorPrimSimple) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"AbstractDeclaratorPrimSimple:",
		adps.identifier)
}

func (adpd AbstractDeclaratorPrimDifficult) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclaratorPrimDifficult:")
	adpd.abstractDeclarator.printNode(offset + 1)
}

func (adp AbstractDeclaratorPointer) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "AbstractDeclaratorPointer:")
	adp.abstractDeclarator.printNode(offset + 1)
}

func (sts SimpleTypeSpecifier) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"SimpleTypeSpecifier:", sts.specType)
}

func (ets EnumTypeSpecifier) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "EnumTypeSpecifier:")
	ets.enumStatement.printNode(offset + 1)
}

func (ies IdentEnumStatement) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"IdentEnumStatement:", ies.ident)
	ies.bodyEnumStatementOpt.printNode(offset + 1)
}

func (bes BodyEnumStatement) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "BodyEnumStatement:")
	bes.enumeratorList.printNode(offset + 1)
}

func (bes EnumeratorList) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "EnumeratorList:")
	for _, e := range bes.enumerators {
		e.printNode(offset + 1)
	}
}

func (bes Enumerator) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"Enumerator:", bes.ident)
	bes.expressionOpt.printNode(offset + 1)
}

func (sus StructOrUnionSpecifier) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"StructOrUnionSpecifier:", sus.typeSpecifier)
	sus.structOrUnionStatement.printNode(offset + 1)
}

func (isus IdentStructOrUnionStatement) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"IdentStructOrUnionStatement:", isus.ident)
	isus.bodyStructOrUnionStatementOpt.printNode(offset + 1)
}

func (bsus BodyStructOrUnionStatement) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "BodyStructOrUnionStatement:")
	bsus.declarationList.printNode(offset + 1)
}

func (boe BinaryOperatorExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "BinaryOperatorExpression:")
	boe.left.printNode(offset + 1)
	fmt.Println(strings.Repeat(offsetString, offset)+"operator:", boe.operation)
	boe.right.printNode(offset + 1)
}

func (soe SizeOfExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"SizeOfExpression:",
		soe.typeDeclaration, soe.ident)
}

func (ie IntegerExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"IntegerExpression:",
		ie.integer)
}

func (ie IdentExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"IdentExpression:",
		ie.ident)
}

func (ie InsideExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "InsideExpression:")
	ie.expr.printNode(offset + 1)
}

func (ne NilExpression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset) + "NilExpression")
}

```

`transitions.go`
```go
package recdesc

import (
	"log"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

func (rp *RecursiveParser) Parse() Program {
	return rp.program()
}

// Program ::= DeclarationList
func (rp *RecursiveParser) program() Program {
	declarationList := rp.declarationList()
	return NewProgram(declarationList)
}

// DeclarationList ::= Declaration*
func (rp *RecursiveParser) declarationList() DeclarationList {
	declarationList := make([]Declaration, 0)
	for isDeclarationStart(rp.currentToken) {
		declarationList = append(declarationList, rp.declaration())
	}

	return NewDeclarationList(declarationList)
}

func isDeclarationStart(t lexer.IToken) bool {
	return lexer.IsKeyword(t.GetValue()) && t.GetValue() != "sizeof"
}

// Declaration ::= TypeSpecifier AbstractDeclaratorsOpt ';'
func (rp *RecursiveParser) declaration() Declaration {

	typeSpecifier := rp.typeSpecifier()
	abstractDeclaratorsOpt := rp.abstractDeclaratorsOpt()

	rp.isExpectedToken(";", lexer.SpecSymbolTag)

	return NewDeclaration(abstractDeclaratorsOpt, typeSpecifier)
}

// AbstractDeclaratorsOpt ::= AbstractDeclarators?
func (rp *RecursiveParser) abstractDeclaratorsOpt() AbstractDeclaratorsOpt {
	if rp.currentToken.GetValue() == ";" {
		return NewAbstractDeclaratorsOpt(AbstractDeclarators{})
	}

	ads := rp.abstractDeclarators()

	return NewAbstractDeclaratorsOpt(ads)
}

// AbstractDeclarators ::= AbstractDeclarator (',' AbstractDeclarator)*
func (rp *RecursiveParser) abstractDeclarators() AbstractDeclarators {
	abstractDeclarators := make([]AbstractDeclarator, 0)
	abstractDeclarators = append(abstractDeclarators, rp.abstractDeclarator())

	for rp.currentToken.GetValue() == "," {
		rp.currentToken = rp.scanner.NextToken()
		abstractDeclarators = append(abstractDeclarators, rp.abstractDeclarator())
	}

	return NewAbstractDeclarators(abstractDeclarators)
}

// AbstractDeclarator ::= AbstractDeclaratorPointer | AbstractDeclaratorArrayList
func (rp *RecursiveParser) abstractDeclarator() AbstractDeclarator {
	if rp.currentToken.GetValue() == "*" {
		return rp.abstractDeclaratorPointer()
	} else {
		return rp.abstractDeclaratorArrayList()
	}
}

// AbstractDeclaratorPointer ::= '*' AbstractDeclarator
func (rp *RecursiveParser) abstractDeclaratorPointer() AbstractDeclaratorPointer {
	rp.currentToken = rp.scanner.NextToken()
	ad := rp.abstractDeclarator()
	return NewAbstractDeclaratorPointer(ad)
}

// AbstractDeclaratorArrayList ::= AbstractDeclaratorArray+
func (rp *RecursiveParser) abstractDeclaratorArrayList() AbstractDeclaratorArrayList {
	abstractDeclaratorArrayList := make([]AbstractDeclaratorArray, 0)
	if rp.currentToken.GetValue() != "[" && rp.currentToken.GetValue() != "(" &&
		rp.currentToken.GetType() != lexer.IdentifierTag {
		log.Fatalf("failed parse in abstractDeclaratorArrayList: expected first"+
			"[ or ( or identifier, get: %s", rp.currentToken.GetValue())
	}
	abstractDeclaratorArrayList = append(abstractDeclaratorArrayList,
		rp.abstractDeclaratorArray())

	for rp.currentToken.GetValue() == "[" || rp.currentToken.GetValue() == "(" ||
		rp.currentToken.GetType() == lexer.IdentifierTag {
		abstractDeclaratorArrayList = append(abstractDeclaratorArrayList,
			rp.abstractDeclaratorArray())
	}

	return NewAbstractDeclaratorArrayList(abstractDeclaratorArrayList)
}

// AbstractDeclaratorArray ::= AbstractDeclaratorPrimArray |
// AbstractDeclaratorPrimSimple | AbstractDeclaratorPrimDifficult
func (rp *RecursiveParser) abstractDeclaratorArray() AbstractDeclaratorArray {
	if rp.currentToken.GetValue() == "[" {
		return rp.abstractDeclaratorPrimArray()
	} else if rp.currentToken.GetValue() == "(" {
		return rp.AbstractDeclaratorPrimDifficult()
	} else {
		return rp.abstractDeclaratorPrimSimple()
	}
}

// AbstractDeclaratorPrimArray ::= '[' Expression ']'
func (rp *RecursiveParser) abstractDeclaratorPrimArray() AbstractDeclaratorPrimArray {
	rp.isExpectedToken("[", lexer.SpecSymbolTag)
	expr := rp.expression()
	rp.isExpectedToken("]", lexer.SpecSymbolTag)
	return NewAbstractDeclaratorPrimArray(expr)
}

// AbstractDeclaratorPrimSimple ::= IDENTIFIER
func (rp *RecursiveParser) abstractDeclaratorPrimSimple() AbstractDeclaratorPrimSimple {
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	return NewAbstractDeclaratorPrimSimple(t.GetValue())
}

// AbstractDeclaratorPrimDifficult ::= '(' AbstractDeclarator ')'
func (rp *RecursiveParser) AbstractDeclaratorPrimDifficult() AbstractDeclaratorPrimDifficult {
	rp.isExpectedToken("(", lexer.SpecSymbolTag)
	ad := rp.abstractDeclarator()
	rp.isExpectedToken(")", lexer.SpecSymbolTag)
	return NewAbstractDeclaratorPrimDifficult(ad)
}

// TypeSpecifier ::= SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier
func (rp *RecursiveParser) typeSpecifier() TypeSpecifier {
	if rp.currentToken.GetType() == lexer.KeywordTag && rp.currentToken.GetValue() != "sizeof" {
		if rp.currentToken.GetValue() == "enum" {
			return rp.enumTypeSpecifier()
		} else if rp.currentToken.GetValue() == "struct" || rp.currentToken.GetValue() == "union" {
			return rp.structOrUnionSpecifier()
		} else {
			return rp.simpleTypeSpecifier()
		}
	}
	log.Fatalf("failed parse in typeSpecifier. Expected declarartion keyword, get %s",
		rp.currentToken)
	return SimpleTypeSpecifier{}
}

// SimpleTypeSpecifier ::= char | short | int | long | float | double
func (rp *RecursiveParser) simpleTypeSpecifier() SimpleTypeSpecifier {
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	return NewSimpleTypeSpecifier(t.GetValue())
}

// EnumTypeSpecifier ::= ENUM EnumStatement
func (rp *RecursiveParser) enumTypeSpecifier() EnumTypeSpecifier {
	rp.isExpectedToken("enum", lexer.KeywordTag)
	es := rp.enumStatement()
	return NewEnumTypeSpecifier(es)
}

// EnumStatement ::= IdentEnumStatement | BodyEnumStatement
func (rp *RecursiveParser) enumStatement() EnumStatement {
	if rp.currentToken.GetType() == lexer.IdentifierTag {
		return rp.identEnumStatement()
	} else {
		return rp.bodyEnumStatement()
	}
}

// IdentEnumStatement ::= IDENTIFIER BodyEnumStatement?
func (rp *RecursiveParser) identEnumStatement() IdentEnumStatement {
	ident := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	var bes BodyEnumStatement
	if rp.currentToken.GetValue() == "{" {
		bes = rp.bodyEnumStatement()
	}
	return NewIdentEnumStatement(ident.GetValue(), bes)
}

// BodyEnumStatement ::= '{' EnumeratorList '}'
func (rp *RecursiveParser) bodyEnumStatement() BodyEnumStatement {
	rp.isExpectedToken("{", lexer.SpecSymbolTag)
	el := rp.enumeratorList()
	rp.isExpectedToken("}", lexer.SpecSymbolTag)
	return NewBodyEnumStatement(el)
}

// EnumeratorList ::= Enumerator (',' Enumerator)*
func (rp *RecursiveParser) enumeratorList() EnumeratorList {
	enumeratorList := make([]Enumerator, 0)
	enumeratorList = append(enumeratorList, rp.enumerator())

	for rp.currentToken.GetValue() == "," {
		rp.currentToken = rp.scanner.NextToken()
		enumeratorList = append(enumeratorList, rp.enumerator())
	}

	return NewEnumeratorList(enumeratorList)
}

// Enumerator ::= IDENTIFIER ('=' Expression)?
func (rp *RecursiveParser) enumerator() Enumerator {
	ident := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	if rp.currentToken.GetValue() == "=" {
		rp.currentToken = rp.scanner.NextToken()
		expr := rp.expression()
		return NewEnumerator(ident.GetValue(), expr)
	}
	return NewEnumerator(ident.GetValue(), NewNilExpression())
}

// StructOrUnionSpecifier ::= (struct | union) StructOrUnionStatement
func (rp *RecursiveParser) structOrUnionSpecifier() StructOrUnionSpecifier {
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	sus := rp.structOrUnionStatement()
	return NewStructOrUnionSpecifier(t.GetValue(), sus)
}

// StructOrUnionStatement ::= IdentStructOrUnionStatement | BodyStructOrUnionStatement
func (rp *RecursiveParser) structOrUnionStatement() StructOrUnionStatement {
	if rp.currentToken.GetType() == lexer.IdentifierTag {
		return rp.identStructOrUnionStatement()
	} else {
		return rp.bodyStructOrUnionStatement()
	}
}

// IdentStructOrUnionStatement ::= IDENTIFIER BodyStructOrUnionStatement?
func (rp *RecursiveParser) identStructOrUnionStatement() IdentStructOrUnionStatement {
	ident := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	var bsus BodyStructOrUnionStatement
	if rp.currentToken.GetValue() == "{" {
		bsus = rp.bodyStructOrUnionStatement()
	}
	return NewIdentStructOrUnionStatement(ident.GetValue(), bsus)
}

// BodyStructOrUnionStatement ::= '{' DeclarationList '}'
func (rp *RecursiveParser) bodyStructOrUnionStatement() BodyStructOrUnionStatement {
	rp.isExpectedToken("{", lexer.SpecSymbolTag)
	dl := rp.declarationList()
	rp.isExpectedToken("}", lexer.SpecSymbolTag)
	return NewBodyStructOrUnionStatement(dl)
}

// Expression ::= Term (('+' | '-') Term)*
func (rp *RecursiveParser) expression() Expression {
	expr := rp.term()
	for rp.currentToken.GetValue() == "+" || rp.currentToken.GetValue() == "-" {
		operation := rp.currentToken.GetValue()
		rp.currentToken = rp.scanner.NextToken()
		rightExpr := rp.term()

		expr = NewBinaryOperatorExpression(expr, operation, rightExpr)
	}

	return expr
}

// Term -> Factor (('*' | '/') Factor)*
func (rp *RecursiveParser) term() Expression {
	expr := rp.factor()
	for rp.currentToken.GetValue() == "*" || rp.currentToken.GetValue() == "/" {
		operation := rp.currentToken.GetValue()
		rp.currentToken = rp.scanner.NextToken()
		rightExpr := rp.factor()

		expr = NewBinaryOperatorExpression(expr, operation, rightExpr)
	}

	return expr
}

// Factor ::= sizeof '(' (struct | union | enum) IDENTIFIER ')' 
//    | IDENTIFIER | INTEGER | '(' Expression ')'
func (rp *RecursiveParser) factor() Expression {
	if rp.currentToken.GetValue() == "sizeof" {
		rp.currentToken = rp.scanner.NextToken()
		rp.isExpectedToken("(", lexer.SpecSymbolTag)
		typeDeclaration := rp.currentToken
		rp.currentToken = rp.scanner.NextToken()
		ident := rp.currentToken
		rp.currentToken = rp.scanner.NextToken()
		rp.isExpectedToken(")", lexer.SpecSymbolTag)
		return NewSizeOf(typeDeclaration.GetValue(), ident.GetValue())
	} else if rp.currentToken.GetType() == lexer.IdentifierTag {
		ident := rp.currentToken
		rp.currentToken = rp.scanner.NextToken()
		return NewIdentExpression(ident.GetValue())
	} else if rp.currentToken.GetType() == lexer.IntTag {
		integer := rp.currentToken
		rp.currentToken = rp.scanner.NextToken()
		return NewIntegerExpression(integer.GetValue())
		return SizeOfExpression{}
	} else if rp.currentToken.GetValue() == "(" {
		rp.isExpectedToken("(", lexer.SpecSymbolTag)
		expr := rp.expression()
		rp.isExpectedToken(")", lexer.SpecSymbolTag)
		return NewInsideExpression(expr)
	}

	log.Fatalf("failed to parse expression, get %s", rp.currentToken.GetValue())
	return SizeOfExpression{}
}

```

`tree_nodes.go`
```go
package recdesc

type Program struct {
	declarationList DeclarationList
}

func NewProgram(dl DeclarationList) Program {
	return Program{declarationList: dl}
}

type DeclarationList struct {
	declarations []Declaration
}

func NewDeclarationList(ds []Declaration) DeclarationList {
	return DeclarationList{declarations: ds}
}

type Declaration struct {
	abstractDeclaratorsOpt AbstractDeclaratorsOpt
	typeSpecifier          TypeSpecifier
}

func NewDeclaration(ado AbstractDeclaratorsOpt, ts TypeSpecifier) Declaration {
	return Declaration{abstractDeclaratorsOpt: ado, typeSpecifier: ts}
}

type AbstractDeclaratorsOpt struct {
	abstractDeclarators AbstractDeclarators
}

func NewAbstractDeclaratorsOpt(ads AbstractDeclarators) AbstractDeclaratorsOpt {
	return AbstractDeclaratorsOpt{abstractDeclarators: ads}
}

type AbstractDeclarators struct {
	abstractDeclarators []AbstractDeclarator
}

func NewAbstractDeclarators(ado []AbstractDeclarator) AbstractDeclarators {
	return AbstractDeclarators{abstractDeclarators: ado}
}

type AbstractDeclarator interface {
	printNode(offset int)
	abstractDeclaratorI()
}

type AbstractDeclaratorPointer struct {
	abstractDeclarator AbstractDeclarator
}

func (adp AbstractDeclaratorPointer) abstractDeclaratorI() {}

func NewAbstractDeclaratorPointer(ad AbstractDeclarator) AbstractDeclaratorPointer {
	return AbstractDeclaratorPointer{abstractDeclarator: ad}
}

type AbstractDeclaratorArrayList struct {
	abstractDeclaratorArray []AbstractDeclaratorArray
}

func (adal AbstractDeclaratorArrayList) abstractDeclaratorI() {}

func NewAbstractDeclaratorArrayList(ada []AbstractDeclaratorArray) AbstractDeclaratorArrayList {
	return AbstractDeclaratorArrayList{abstractDeclaratorArray: ada}
}

type AbstractDeclaratorArray interface {
	printNode(offset int)
	abstractDeclaratorArrayI()
}

type AbstractDeclaratorPrimArray struct {
	expression Expression
}

func (adpa AbstractDeclaratorPrimArray) abstractDeclaratorArrayI() {}

func NewAbstractDeclaratorPrimArray(expr Expression) AbstractDeclaratorPrimArray {
	return AbstractDeclaratorPrimArray{expression: expr}
}

type AbstractDeclaratorPrimSimple struct {
	identifier string
}

func (adps AbstractDeclaratorPrimSimple) abstractDeclaratorArrayI() {}

func NewAbstractDeclaratorPrimSimple(ident string) AbstractDeclaratorPrimSimple {
	return AbstractDeclaratorPrimSimple{identifier: ident}
}

type AbstractDeclaratorPrimDifficult struct {
	abstractDeclarator AbstractDeclarator
}

func (adpd AbstractDeclaratorPrimDifficult) abstractDeclaratorArrayI() {}

func NewAbstractDeclaratorPrimDifficult(ad AbstractDeclarator) AbstractDeclaratorPrimDifficult {
	return AbstractDeclaratorPrimDifficult{abstractDeclarator: ad}
}

type TypeSpecifier interface {
	printNode(offset int)
	typeSpecifierI()
}

type SimpleTypeSpecifier struct {
	specType string
}

func (sts SimpleTypeSpecifier) typeSpecifierI() {}

func NewSimpleTypeSpecifier(st string) SimpleTypeSpecifier {
	return SimpleTypeSpecifier{specType: st}
}

type EnumTypeSpecifier struct {
	enumStatement EnumStatement
}

func (ets EnumTypeSpecifier) typeSpecifierI() {}

func NewEnumTypeSpecifier(es EnumStatement) EnumTypeSpecifier {
	return EnumTypeSpecifier{enumStatement: es}
}

type EnumStatement interface {
	printNode(offset int)
	enumStatementI()
}

type IdentEnumStatement struct {
	ident                string
	bodyEnumStatementOpt BodyEnumStatement
}

func (ies IdentEnumStatement) enumStatementI() {}

func NewIdentEnumStatement(ident string, bes BodyEnumStatement) IdentEnumStatement {
	return IdentEnumStatement{ident: ident, bodyEnumStatementOpt: bes}
}

type BodyEnumStatement struct {
	enumeratorList EnumeratorList
}

func (bes BodyEnumStatement) enumStatementI() {}

func NewBodyEnumStatement(el EnumeratorList) BodyEnumStatement {
	return BodyEnumStatement{enumeratorList: el}
}

type EnumeratorList struct {
	enumerators []Enumerator
}

func NewEnumeratorList(es []Enumerator) EnumeratorList {
	return EnumeratorList{enumerators: es}
}

type Enumerator struct {
	ident         string
	expressionOpt Expression
}

func NewEnumerator(ident string, eOpt Expression) Enumerator {
	return Enumerator{ident: ident, expressionOpt: eOpt}
}

type StructOrUnionSpecifier struct {
	typeSpecifier          string
	structOrUnionStatement StructOrUnionStatement
}

func (sus StructOrUnionSpecifier) typeSpecifierI() {}

func NewStructOrUnionSpecifier(
	typeSpecifier string, sus StructOrUnionStatement,
	) StructOrUnionSpecifier {
	return StructOrUnionSpecifier{typeSpecifier: typeSpecifier, structOrUnionStatement: sus}
}

type StructOrUnionStatement interface {
	printNode(offset int)
	structOrUnionStatementI()
}

type IdentStructOrUnionStatement struct {
	ident                         string
	bodyStructOrUnionStatementOpt BodyStructOrUnionStatement
}

func (isus IdentStructOrUnionStatement) structOrUnionStatementI() {}

func NewIdentStructOrUnionStatement(
	ident string, bsus BodyStructOrUnionStatement,
	) IdentStructOrUnionStatement {
	return IdentStructOrUnionStatement{ident: ident, bodyStructOrUnionStatementOpt: bsus}
}

type BodyStructOrUnionStatement struct {
	declarationList DeclarationList
}

func (bsus BodyStructOrUnionStatement) structOrUnionStatementI() {}

func NewBodyStructOrUnionStatement(dl DeclarationList) BodyStructOrUnionStatement {
	return BodyStructOrUnionStatement{declarationList: dl}
}

type Expression interface {
	printNode(offset int)
	expressionI()
}

type BinaryOperatorExpression struct {
	left      Expression
	operation string
	right     Expression
}

func (boe BinaryOperatorExpression) expressionI() {}

func NewBinaryOperatorExpression(
	left Expression, operation string, right Expression,
	) BinaryOperatorExpression {
	return BinaryOperatorExpression{left: left, operation: operation, right: right}
}

type SizeOfExpression struct {
	typeDeclaration string
	ident           string
}

func (soe SizeOfExpression) expressionI() {}

func NewSizeOf(typeDeclaration, ident string) SizeOfExpression {
	return SizeOfExpression{typeDeclaration: typeDeclaration, ident: ident}
}

type IntegerExpression struct {
	integer string
}

func (ie IntegerExpression) expressionI() {}

func NewIntegerExpression(integer string) IntegerExpression {
	return IntegerExpression{integer: integer}
}

type IdentExpression struct {
	ident string
}

func (ie IdentExpression) expressionI() {}

func NewIdentExpression(ident string) IdentExpression {
	return IdentExpression{ident: ident}
}

type InsideExpression struct {
	expr Expression
}

func (ie InsideExpression) expressionI() {}

func NewInsideExpression(expr Expression) InsideExpression {
	return InsideExpression{expr: expr}
}

type NilExpression struct {
}

func (ne NilExpression) expressionI() {}

func NewNilExpression() NilExpression {
	return NilExpression{}
}

```

### main файл
`main.go`
```go
package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/cd/lab9/recdesc"

	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
)

const filepath = "test_files/mixed.txt"

func main() {

	scanner := lexer.NewScanner(filepath)

	scanner.PrintTokens()
	fmt.Println()

	rp := recdesc.NewParser(scanner)
	program := rp.Parse()
	program.Print()

	fmt.Println("\n\n\nfinish")
}

```

# Тестирование

Входные данные

```
struct Coords {
  int x, y;
};

enum ScreenType { aaa = bab } *screenType[5 + 5], **fas;

enum ScreenType { aaa = bab } *screenType[5 + 5], **fas;

enum ScreenType { aaa = bab } a[1], *a[1], (*a)[1];

enum Color {
  COLOR_RED = 1,
  COLOR_GREEN = 2,
  COLOR_BLUE = 2*2,
  COLOR_HIGHLIGHT = 8
};

enum ScreenType {
  SCREEN_TYPE_TEXT,
  SCREEN_TYPE_GRAPHIC
} screen_type;

enum {
  HPIXELS = 480,
  WPIXELS = 640,
  HCHARS = 24,
  WCHARS = 80
};

struct ScreenChar {
  char symbol;
  enum Color sym_color;
  enum Color back_color;
};

struct TextScreen {
  struct ScreenChar chars[HCHARS][WCHARS];
};

struct GrahpicScreen {
  enum Color pixels[HPIXELS][WPIXELS];
};

union Screen {
  struct TextScreen text;
  struct GraphicScreen graphic;
};

enum {
  BUFFER_SIZE = sizeof(union Screen),
  PAGE_SIZE = 4096,
  PAGES_FOR_BUFFER = (BUFFER_SIZE + PAGE_SIZE - 1) / PAGE_SIZE
};

struct Token {
  struct Fragment {
    struct Pos {
      int line;
      int col;
    } starting, following;
  } fragment;

  enum { Ident, IntConst, FloatConst } type;

  union {
    char *name;
    int int_value;
    double float_value;
  } info;
};

struct List {
  struct Token value;
  struct List *next;
};

```

Вывод на `stdout`

```
Program:
..DeclarationList:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: Coords
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................SimpleTypeSpecifier: int
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: x
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: y
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......EnumTypeSpecifier:
........IdentEnumStatement: ScreenType
..........BodyEnumStatement:
............EnumeratorList:
..............Enumerator: aaa
................IdentExpression: bab
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
..........AbstractDeclaratorPointer:
............AbstractDeclaratorArrayList:
..............AbstractDeclaratorPrimSimple: screenType
..............AbstractDeclaratorPrimArray:
................BinaryOperatorExpression:
..................IntegerExpression: 5
................operator: +
..................IntegerExpression: 5
..........AbstractDeclaratorPointer:
............AbstractDeclaratorPointer:
..............AbstractDeclaratorArrayList:
................AbstractDeclaratorPrimSimple: fas
....Declaration:
......EnumTypeSpecifier:
........IdentEnumStatement: ScreenType
..........BodyEnumStatement:
............EnumeratorList:
..............Enumerator: aaa
................IdentExpression: bab
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
..........AbstractDeclaratorPointer:
............AbstractDeclaratorArrayList:
..............AbstractDeclaratorPrimSimple: screenType
..............AbstractDeclaratorPrimArray:
................BinaryOperatorExpression:
..................IntegerExpression: 5
................operator: +
..................IntegerExpression: 5
..........AbstractDeclaratorPointer:
............AbstractDeclaratorPointer:
..............AbstractDeclaratorArrayList:
................AbstractDeclaratorPrimSimple: fas
....Declaration:
......EnumTypeSpecifier:
........IdentEnumStatement: ScreenType
..........BodyEnumStatement:
............EnumeratorList:
..............Enumerator: aaa
................IdentExpression: bab
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
..........AbstractDeclaratorArrayList:
............AbstractDeclaratorPrimSimple: a
............AbstractDeclaratorPrimArray:
..............IntegerExpression: 1
..........AbstractDeclaratorPointer:
............AbstractDeclaratorArrayList:
..............AbstractDeclaratorPrimSimple: a
..............AbstractDeclaratorPrimArray:
................IntegerExpression: 1
..........AbstractDeclaratorArrayList:
............AbstractDeclaratorPrimDifficult:
..............AbstractDeclaratorPointer:
................AbstractDeclaratorArrayList:
..................AbstractDeclaratorPrimSimple: a
............AbstractDeclaratorPrimArray:
..............IntegerExpression: 1
....Declaration:
......EnumTypeSpecifier:
........IdentEnumStatement: Color
..........BodyEnumStatement:
............EnumeratorList:
..............Enumerator: COLOR_RED
................IntegerExpression: 1
..............Enumerator: COLOR_GREEN
................IntegerExpression: 2
..............Enumerator: COLOR_BLUE
................BinaryOperatorExpression:
..................IntegerExpression: 2
................operator: *
..................IntegerExpression: 2
..............Enumerator: COLOR_HIGHLIGHT
................IntegerExpression: 8
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......EnumTypeSpecifier:
........IdentEnumStatement: ScreenType
..........BodyEnumStatement:
............EnumeratorList:
..............Enumerator: SCREEN_TYPE_TEXT
................NilExpression
..............Enumerator: SCREEN_TYPE_GRAPHIC
................NilExpression
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
..........AbstractDeclaratorArrayList:
............AbstractDeclaratorPrimSimple: screen_type
....Declaration:
......EnumTypeSpecifier:
........BodyEnumStatement:
..........EnumeratorList:
............Enumerator: HPIXELS
..............IntegerExpression: 480
............Enumerator: WPIXELS
..............IntegerExpression: 640
............Enumerator: HCHARS
..............IntegerExpression: 24
............Enumerator: WCHARS
..............IntegerExpression: 80
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: ScreenChar
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................SimpleTypeSpecifier: char
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: symbol
..............Declaration:
................EnumTypeSpecifier:
..................IdentEnumStatement: Color
....................BodyEnumStatement:
......................EnumeratorList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: sym_color
..............Declaration:
................EnumTypeSpecifier:
..................IdentEnumStatement: Color
....................BodyEnumStatement:
......................EnumeratorList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: back_color
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: TextScreen
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: ScreenChar
....................BodyStructOrUnionStatement:
......................DeclarationList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: chars
......................AbstractDeclaratorPrimArray:
........................IdentExpression: HCHARS
......................AbstractDeclaratorPrimArray:
........................IdentExpression: WCHARS
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: GrahpicScreen
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................EnumTypeSpecifier:
..................IdentEnumStatement: Color
....................BodyEnumStatement:
......................EnumeratorList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: pixels
......................AbstractDeclaratorPrimArray:
........................IdentExpression: HPIXELS
......................AbstractDeclaratorPrimArray:
........................IdentExpression: WPIXELS
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: union
........IdentStructOrUnionStatement: Screen
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: TextScreen
....................BodyStructOrUnionStatement:
......................DeclarationList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: text
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: GraphicScreen
....................BodyStructOrUnionStatement:
......................DeclarationList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: graphic
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......EnumTypeSpecifier:
........BodyEnumStatement:
..........EnumeratorList:
............Enumerator: BUFFER_SIZE
..............SizeOfExpression: union Screen
............Enumerator: PAGE_SIZE
..............IntegerExpression: 4096
............Enumerator: PAGES_FOR_BUFFER
..............BinaryOperatorExpression:
................InsideExpression:
..................BinaryOperatorExpression:
....................BinaryOperatorExpression:
......................IdentExpression: BUFFER_SIZE
....................operator: +
......................IdentExpression: PAGE_SIZE
..................operator: -
....................IntegerExpression: 1
..............operator: /
................IdentExpression: PAGE_SIZE
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: Token
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: Fragment
....................BodyStructOrUnionStatement:
......................DeclarationList:
........................Declaration:
..........................StructOrUnionSpecifier: struct
............................IdentStructOrUnionStatement: Pos
..............................BodyStructOrUnionStatement:
................................DeclarationList:
..................................Declaration:
....................................SimpleTypeSpecifier: int
....................................AbstractDeclaratorsOpt:
......................................AbstractDeclarators:
........................................AbstractDeclaratorArrayList:
..........................................AbstractDeclaratorPrimSimple: line
..................................Declaration:
....................................SimpleTypeSpecifier: int
....................................AbstractDeclaratorsOpt:
......................................AbstractDeclarators:
........................................AbstractDeclaratorArrayList:
..........................................AbstractDeclaratorPrimSimple: col
..........................AbstractDeclaratorsOpt:
............................AbstractDeclarators:
..............................AbstractDeclaratorArrayList:
................................AbstractDeclaratorPrimSimple: starting
..............................AbstractDeclaratorArrayList:
................................AbstractDeclaratorPrimSimple: following
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: fragment
..............Declaration:
................EnumTypeSpecifier:
..................BodyEnumStatement:
....................EnumeratorList:
......................Enumerator: Ident
........................NilExpression
......................Enumerator: IntConst
........................NilExpression
......................Enumerator: FloatConst
........................NilExpression
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: type
..............Declaration:
................StructOrUnionSpecifier: union
..................BodyStructOrUnionStatement:
....................DeclarationList:
......................Declaration:
........................SimpleTypeSpecifier: char
........................AbstractDeclaratorsOpt:
..........................AbstractDeclarators:
............................AbstractDeclaratorPointer:
..............................AbstractDeclaratorArrayList:
................................AbstractDeclaratorPrimSimple: name
......................Declaration:
........................SimpleTypeSpecifier: int
........................AbstractDeclaratorsOpt:
..........................AbstractDeclarators:
............................AbstractDeclaratorArrayList:
..............................AbstractDeclaratorPrimSimple: int_value
......................Declaration:
........................SimpleTypeSpecifier: double
........................AbstractDeclaratorsOpt:
..........................AbstractDeclarators:
............................AbstractDeclaratorArrayList:
..............................AbstractDeclaratorPrimSimple: float_value
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: info
......AbstractDeclaratorsOpt:
........AbstractDeclarators:
....Declaration:
......StructOrUnionSpecifier: struct
........IdentStructOrUnionStatement: List
..........BodyStructOrUnionStatement:
............DeclarationList:
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: Token
....................BodyStructOrUnionStatement:
......................DeclarationList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorArrayList:
......................AbstractDeclaratorPrimSimple: value
..............Declaration:
................StructOrUnionSpecifier: struct
..................IdentStructOrUnionStatement: List
....................BodyStructOrUnionStatement:
......................DeclarationList:
................AbstractDeclaratorsOpt:
..................AbstractDeclarators:
....................AbstractDeclaratorPointer:
......................AbstractDeclaratorArrayList:
........................AbstractDeclaratorPrimSimple: next
......AbstractDeclaratorsOpt:
........AbstractDeclarators:

```

# Вывод
В этой работе был повторен алгоритм построения парсеров методом рекурсивного спуска.
Мы занимались рекурсивным спуском на 1 семестре, однако сейчас грамматика была сильно сложнее
и благодаря лекциям я шел по основным алгоритмам раскрытия правил и получил полное 
понимание, почему нужно делать именно так, а не иначе. Трудностью у меня было
перевести грамматики из LALR(1) в LL(1), но в итоге это получилось.
Нравится то, как обычное представление грамматики, которое я использовал в 2.2,
сжалось примерно на 35% после того, как я записал ее в РБНФ.

В этой лабораторной, в отличие от 2.2, нужно было руками прописать все переходы,
а не просто описать правила. Плюсом я вижу то, что наше дерево находится под полным 
нашим контролем, и в не зависимости от сложности и размера грамматики мы сможем
сделать рекурсивный спуск (если она LL1), используя ограниченный набор правил.
Времени это заняло достаточно много, но постепенно,
когда многие узлы уже были описаны, нужно повторять похожую логику.
Метод, как мне кажется, один из основных и я рад его закрепить.

