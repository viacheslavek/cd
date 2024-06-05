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

//
//
//
//
//
// TODO: делаю блок 5 - выражения

// Expression TODO: раскрываю expression
type Expression struct {
	todo string
}

//
//
//
//

func NewExpression(todo string) Expression {
	return Expression{todo: todo}
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

//
//
//
//
// TODO: делаю блок 4 struct and union

type StructOrUnionSpecifier struct {
	typeSpecifier          string
	structOrUnionStatement StructOrUnionStatement
}

func (sus StructOrUnionSpecifier) typeSpecifierI() {}

func NewStructOrUnionSpecifier(typeSpecifier string, sus StructOrUnionStatement) StructOrUnionSpecifier {
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

func NewIdentStructOrUnionStatement(ident string, bsus BodyStructOrUnionStatement) IdentStructOrUnionStatement {
	return IdentStructOrUnionStatement{ident: ident, bodyStructOrUnionStatementOpt: bsus}
}

type BodyStructOrUnionStatement struct {
	declarationList DeclarationList
}

func (bsus BodyStructOrUnionStatement) structOrUnionStatementI() {}

func NewBodyStructOrUnionStatement(dl DeclarationList) BodyStructOrUnionStatement {
	return BodyStructOrUnionStatement{declarationList: dl}
}
