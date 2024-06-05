package recdesc

import (
	"fmt"
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
	abstractDeclaratorArrayList = append(abstractDeclaratorArrayList, rp.abstractDeclaratorArray())

	for rp.currentToken.GetValue() == "[" || rp.currentToken.GetValue() == "(" ||
		rp.currentToken.GetType() == lexer.IdentifierTag {
		abstractDeclaratorArrayList = append(abstractDeclaratorArrayList, rp.abstractDeclaratorArray())
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

// TODO: делаю блок 5
//
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

// TODO: делаю
// Factor ::= sizeof '(' (struct | union | enum) IDENTIFIER ')' | IDENTIFIER | INTEGER | '(' Expression ')'
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

	fmt.Println("WTF?")

	log.Fatalf("failed to parse expression, get %s", rp.currentToken.GetValue())
	return SizeOfExpression{}
}
