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
	fmt.Println(strings.Repeat(offsetString, offset)+"AbstractDeclaratorPrimSimple:", adps.identifier)
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
