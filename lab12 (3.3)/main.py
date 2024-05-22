import abc
import enum
from abc import ABC

import parser_edsl as pe
from dataclasses import dataclass
from pprint import pprint


class SemanticError(pe.Error, ABC):
    pass


class RepeatedIdentifier(SemanticError):
    def __init__(self, pos, ident):
        self.pos = pos
        self.ident = ident

    @property
    def message(self):
        return f'Повторный идентификатор {self.ident}'


@dataclass
class TypeSpecifier(abc.ABC):
    pass


class SimpleType(enum.Enum):
    Char = "CHAR"
    Short = "SHORT"
    Int = "INT"
    Long = "LONG"
    Float = "FLOAT"
    Double = "DOUBLE"
    Signed = "SIGNED"
    Unsigned = "UNSIGNED"


class Expression(abc.ABC):
    pass


@dataclass
class ConstantExpression:
    expression: Expression


@dataclass
class Enumerator:
    identifier: str
    constantExpression: ConstantExpression


@dataclass
class EnumeratorList:
    body: list[Enumerator]


@dataclass
class StructOrUnionStatement(abc.ABC):
    pass


@dataclass
class EmptyStructOrUnionStatement(StructOrUnionStatement):
    identifier: str


@dataclass
class StructOrUnionSpecifier(TypeSpecifier):
    type: str
    structOrUnionSpecifier: StructOrUnionStatement


@dataclass
class EnumStatement(abc.ABC):
    pass


@dataclass
class EnumTypeSpecifier(TypeSpecifier):
    enumStatement: EnumStatement


@dataclass
class FullEnumStatement(EnumStatement):
    identifier: str
    enumeratorList: EnumeratorList
    isEndComma: bool


@dataclass
class EmptyEnumStatement(EnumStatement):
    identifier: str


@dataclass
class IdentifierExpression(Expression):
    identifier: str


@dataclass
class IntExpression(Expression):
    value: int


@dataclass
class BinaryOperationExpression(Expression):
    left: Expression
    operation: str
    right: Expression


@dataclass
class UnaryOperationExpression(Expression):
    operation: str
    expression: Expression


@dataclass
class ListArraysOpt:
    listArraysOpt: list[str]


@dataclass
class AbstractDeclarator(abc.ABC):
    def check(self):
        pass


@dataclass
class AbstractDeclaratorPointer(AbstractDeclarator):
    declarator: AbstractDeclarator

    def check(self):
        print("abstractDeclaratorPointer", self.declarator)
        self.declarator.check()


@dataclass
class AbstractDeclaratorArrayList:
    arrays: list[AbstractDeclarator]

    def check(self):
        print("arrays", self.arrays)
        for ad in self.arrays:
            print("ad in arrays", ad)
            ad.check()


@dataclass
class AbstractDeclaratorArray(AbstractDeclarator):
    declarator: Expression

    def check(self):
        print("abstractDeclaratorArray", self.declarator)
        # TODO: здесь уже смогу смотреть expression


@dataclass
class AbstractDeclaratorsOpt:
    abstractDeclaratorList: list[AbstractDeclarator]

    def check(self):
        for ad in self.abstractDeclaratorList:
            ad.check()
            print()


@dataclass
class SimpleTypeSpecifier(TypeSpecifier):
    simpleType: SimpleType


@dataclass
class SizeofExpression(Expression):
    declarationBody: TypeSpecifier
    varName: AbstractDeclaratorsOpt


@dataclass
class Declaration:
    declarationBody: TypeSpecifier
    varName: AbstractDeclaratorsOpt

    def check(self):
        # Проверяю первый пункт
        self.varName.check()


@dataclass
class AbstractDeclaratorPrim(abc.ABC):
    pass


@dataclass
class AbstractDeclaratorPrimSimple(AbstractDeclaratorPrim):
    identifier: str
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        print("in create", self, coords, res_coord)
        ident, = self
        idc, = coords
        print("ident", ident, "idc", idc)
        return AbstractDeclaratorPrimSimple(ident, idc.start)

    def check(self):
        print("in AbstractDeclaratorPrimSimple", self.identifier)
        if check_and_add_to_map(self.identifier):
            raise RepeatedIdentifier(self.identifier_pos, self.identifier)


@dataclass
class AbstractDeclaratorPrimDifficult(AbstractDeclaratorPrim):
    identifier: AbstractDeclarator

    def check(self):
        print("in AbstractDeclaratorPrimDifficult")
        self.identifier.check()


@dataclass
class FullStructOrUnionStatement(StructOrUnionStatement):
    identifierOpt: str
    declarationList: list[Declaration]


esuIdent = {}


def check_and_add_to_map(s):
    if s in esuIdent:
        print("повтор")
        return True
    else:
        esuIdent[s] = True
        return False


constName = {}


@dataclass
class Program:
    declarationList: list[Declaration]

    def check(self):
        for declaration in self.declarationList:
            declaration.check()


NProgram = pe.NonTerminal('Program')

NDeclarationList = pe.NonTerminal('DeclarationList')
NDeclaration = pe.NonTerminal('Declaration')

NAbstractDeclaratorsOpt = pe.NonTerminal('AbstractDeclaratorsOpt')
NAbstractDeclarators = pe.NonTerminal('AbstractDeclarators')
NAbstractDeclarator = pe.NonTerminal('AbstractDeclarator')

NAbstractDeclaratorStar = pe.NonTerminal('AbstractDeclaratorStar')

NAbstractDeclaratorArrayList = pe.NonTerminal('AbstractDeclaratorArrayList')

NAbstractDeclaratorArray = pe.NonTerminal('AbstractDeclaratorArray')

NAbstractDeclaratorPrim = pe.NonTerminal('AbstractDeclaratorPrim')

NAbstractDeclaratorPrimSimple = pe.NonTerminal('AbstractDeclaratorPrimSimple')
NAbstractDeclaratorPrimDifficult = pe.NonTerminal('AbstractDeclaratorPrimDifficult')

NTypeSpecifier = pe.NonTerminal('TypeSpecifier')

NEnumTypeSpecifier = pe.NonTerminal('EnumTypeSpecifier')

NSimpleTypeSpecifier = pe.NonTerminal('SimpleTypeSpecifier')
NSimpleType = pe.NonTerminal('SimpleType')

NEnumStatement = pe.NonTerminal('EnumStatement')

NFullEnumStatement = pe.NonTerminal('FullEnumStatement')
NEmptyEnumStatement = pe.NonTerminal('EmptyEnumStatement')

NEnumeratorList = pe.NonTerminal('EnumeratorList')
NEnumerator = pe.NonTerminal('Enumerator')

NEnumeratorExpressionOpt = pe.NonTerminal('EnumeratorExpressionOpt')
NConstantExpression = pe.NonTerminal('ConstantExpression')

NIdentifierOpt = pe.NonTerminal('IdentifierOpt')


NCommaOpt = pe.NonTerminal('CommaOpt')


NExpression = pe.NonTerminal('Expression')

NArithmeticExpression = pe.NonTerminal('ArithmeticExpression')
NTerm = pe.NonTerminal('Term')
NAddOperation = pe.NonTerminal('AddOperation')

NFactor = pe.NonTerminal('Factor')
NMultyOperation = pe.NonTerminal('MultyOperation')


NStructOrUnionSpecifier = pe.NonTerminal('StructOrUnionSpecifier')

NStructOrUnion = pe.NonTerminal('StructOrUnion')

NStructOrUnionStatement = pe.NonTerminal('StructOrUnionStatement')

NFullStructOrUnionStatement = pe.NonTerminal('FullStructOrUnionStatement')
NEmptyStructOrUnionStatement = pe.NonTerminal('EmptyStructOrUnionStatement')


def make_keyword(image):
    return pe.Terminal(image, image, lambda _: None, priority=10)


KW_ENUM = make_keyword('enum')
KW_STRUCT = make_keyword('struct')
KW_UNION = make_keyword('union')

KW_SIZEOF = make_keyword('sizeof')

KW_CHAR, KW_SHORT, KW_INT, KW_LONG, KW_FLOAT, KW_DOUBLE, KW_SIGNED, KW_UNSIGNED = \
    map(make_keyword, 'char short int long float double signed unsigned'.split())

INTEGER = pe.Terminal('IDENTIFIER', r'[0-9]*', str)

IDENTIFIER = pe.Terminal('IDENTIFIER', r'[A-Za-z_]([A-Za-z_0-9])*', str)


NProgram |= NDeclarationList, Program

NDeclarationList |= lambda: []
NDeclarationList |= NDeclarationList, NDeclaration, lambda dl, d: dl + [d]

NDeclaration |= NTypeSpecifier, NAbstractDeclaratorsOpt, ';', Declaration

NAbstractDeclaratorsOpt |= lambda: AbstractDeclaratorsOpt(list())
NAbstractDeclaratorsOpt |= NAbstractDeclarators, AbstractDeclaratorsOpt

NAbstractDeclarators |= NAbstractDeclarator, lambda a: [a]
NAbstractDeclarators |= NAbstractDeclarators, ',', NAbstractDeclarator, lambda ads, a: ads + [a]

NAbstractDeclarator |= NAbstractDeclaratorStar, AbstractDeclaratorPointer
NAbstractDeclarator |= NAbstractDeclaratorArrayList, AbstractDeclaratorArrayList

NAbstractDeclaratorStar |= '*', NAbstractDeclarator, AbstractDeclaratorPointer

NAbstractDeclaratorArrayList |= NAbstractDeclaratorArray, lambda a: [a]
NAbstractDeclaratorArrayList |= (NAbstractDeclaratorArrayList, NAbstractDeclaratorArray,
                                 lambda adalo, a: adalo + [a])

NAbstractDeclaratorArray |= '[', NExpression, ']', AbstractDeclaratorArray

NAbstractDeclaratorArray |= NAbstractDeclaratorPrim
NAbstractDeclaratorPrim |= NAbstractDeclaratorPrimSimple, AbstractDeclaratorPrimSimple.create

NAbstractDeclaratorPrim |= NAbstractDeclaratorPrimDifficult, AbstractDeclaratorPrimDifficult

NAbstractDeclaratorPrimSimple |= IDENTIFIER

NAbstractDeclaratorPrimDifficult |= '(', NAbstractDeclarator, ')'


NTypeSpecifier |= NEnumTypeSpecifier
NTypeSpecifier |= NSimpleTypeSpecifier
NTypeSpecifier |= NStructOrUnionSpecifier

NSimpleTypeSpecifier |= NSimpleType, SimpleTypeSpecifier

NSimpleType |= KW_CHAR, lambda: SimpleType.Char
NSimpleType |= KW_SHORT, lambda: SimpleType.Short
NSimpleType |= KW_INT, lambda: SimpleType.Int
NSimpleType |= KW_LONG, lambda: SimpleType.Long
NSimpleType |= KW_FLOAT, lambda: SimpleType.Float
NSimpleType |= KW_DOUBLE, lambda: SimpleType.Double
NSimpleType |= KW_SIGNED, lambda: SimpleType.Signed
NSimpleType |= KW_UNSIGNED, lambda: SimpleType.Unsigned


NEnumTypeSpecifier |= KW_ENUM, NEnumStatement, EnumTypeSpecifier

NEnumStatement |= NFullEnumStatement

NFullEnumStatement |= NIdentifierOpt, '{', NEnumeratorList, NCommaOpt, '}', FullEnumStatement

NIdentifierOpt |= lambda: ""
NIdentifierOpt |= IDENTIFIER

NEnumStatement |= NEmptyEnumStatement, EmptyEnumStatement

NEmptyEnumStatement |= IDENTIFIER

NEnumeratorList |= NEnumerator, lambda e: [e]
NEnumeratorList |= NEnumeratorList, ',', NEnumerator, lambda el, e: el + [e]

NEnumerator |= IDENTIFIER, NEnumeratorExpressionOpt, Enumerator

NEnumeratorExpressionOpt |= '=', NConstantExpression, ConstantExpression
NEnumeratorExpressionOpt |= lambda: ""

NCommaOpt |= ',', lambda: True
NCommaOpt |= lambda: False


NConstantExpression |= NExpression

NExpression |= NArithmeticExpression

NArithmeticExpression |= NTerm
NArithmeticExpression |= '+', NTerm, lambda t: UnaryOperationExpression('+', t)
NArithmeticExpression |= '-', NTerm, lambda t: UnaryOperationExpression('-', t)
NArithmeticExpression |= NArithmeticExpression, NAddOperation, NTerm, BinaryOperationExpression

NAddOperation |= '+', lambda: '+'
NAddOperation |= '-', lambda: '-'

NTerm |= NFactor
NTerm |= NTerm, NMultyOperation, NFactor, BinaryOperationExpression

NMultyOperation |= '*', lambda: '*'
NMultyOperation |= '/', lambda: '/'

NFactor |= '(', NExpression, ')'

NFactor |= INTEGER, IntExpression

NFactor |= IDENTIFIER, IdentifierExpression

NFactor |= KW_SIZEOF, '(', NTypeSpecifier, NAbstractDeclaratorsOpt, ')', SizeofExpression


NStructOrUnionSpecifier |= NStructOrUnion, NStructOrUnionStatement, StructOrUnionSpecifier

NStructOrUnion |= KW_STRUCT, lambda: "struct"
NStructOrUnion |= KW_UNION, lambda: "union"

NStructOrUnionStatement |= NFullStructOrUnionStatement
NStructOrUnionStatement |= NEmptyStructOrUnionStatement

NEmptyStructOrUnionStatement |= IDENTIFIER, EmptyStructOrUnionStatement

NFullStructOrUnionStatement |= NIdentifierOpt,  '{', NDeclarationList, '}', FullStructOrUnionStatement


def main():
    p = pe.Parser(NProgram)

    assert p.is_lalr_one()

    p.add_skipped_domain('\\s')

    # files = ["tests/sem_first.txt"]
    files = ["tests/mixed.txt"]

    for filename in files:
        print("file:", filename)
        try:
            with open(filename) as f:
                tree = p.parse(f.read())
                pprint(tree)
                print()
                tree.check()
        except pe.Error as e:
            print(f'Ошибка {e.pos}: {e.message}')
        except Exception as e:
            print(e)


main()


print("esuIdent:", esuIdent)
