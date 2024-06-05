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

// TODO: делаю блок 5 - выражения

// Expression TODO: раскрываю expression
type Expression struct {
	todo string
}

func NewExpression(todo string) Expression {
	return Expression{todo: todo}
}

//
//
//
//
// TODO: делаю блок 3 и 4

// TypeSpecifier TODO: дальше раскрываю TypeSpecifier -> блок номер 3 и 4
type TypeSpecifier struct {
	tsInTODO string
}

func NewTypeSpecifier(s string) TypeSpecifier {
	return TypeSpecifier{tsInTODO: s}
}
