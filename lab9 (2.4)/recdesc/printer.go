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

// TODO: блок 5
//
// TODO: правлю вывод Expression
func (e Expression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"Expression:", e.todo)

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

// TODO: правлю вывод
func (bes BodyEnumStatement) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"BodyEnumStatement:", bes.todo)
}

// TODO: блок 4 - struct or union

// TODO: правлю вывод StructOrUnionSpecifier
func (sus StructOrUnionSpecifier) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"StructOrUnionSpecifier:", sus.todo)
}
