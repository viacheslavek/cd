package recdesc

import (
	"github.com/VyacheslavIsWorkingNow/cd/lab9/lexer"
	"log"
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

// TODO: делаю блок 5

// TODO: поправить
// Expression ::= ArithmeticExpression
// у меня пока это int
func (rp *RecursiveParser) expression() Expression {
	if rp.currentToken.GetType() != lexer.IntTag {
		log.Fatalf("expression isn`t int. expr is %s", rp.currentToken)
	}
	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	return NewExpression(t.GetValue())
}

//
//
//
//
//
//
//
//
//
//
//
// TODO: блок 3 и 4

// TODO: переделываю TypeSpecifier
// TypeSpecifier ::= SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier
func (rp *RecursiveParser) typeSpecifier() TypeSpecifier {

	t := rp.currentToken
	rp.currentToken = rp.scanner.NextToken()
	return NewTypeSpecifier(t.GetValue())
}
