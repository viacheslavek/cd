package recdesc

import (
	"fmt"
	"strings"
)

const offsetString = ".."

type treeNodePrinter interface {
	printNode(offset int)
}

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

// TODO: правлю вывод Expression
func (e Expression) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"Expression:", e.todo)

}

// TODO: блок 3 и 4
//
// TODO: правлю вывод TypeSpecifier
func (ts TypeSpecifier) printNode(offset int) {
	fmt.Println(strings.Repeat(offsetString, offset)+"TypeSpecifier:", ts.tsInTODO)
}
