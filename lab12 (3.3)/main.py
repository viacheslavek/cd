import abc
import enum
from abc import ABC

import parser_edsl as pe
from dataclasses import dataclass
from pprint import pprint


class SemanticError(pe.Error, ABC):
    pass


class RepeatedTag(SemanticError):
    def __init__(self, pos, ident):
        self.pos = pos
        self.ident = ident

    @property
    def message(self):
        return f'Повторный тег {self.ident}'


class RepeatedIdentifier(SemanticError):
    def __init__(self, pos, ident):
        self.pos = pos
        self.ident = ident

    @property
    def message(self):
        return f'Повторный идентификатор {self.ident}'


class RepeatedConstant(SemanticError):
    def __init__(self, pos, ident):
        self.pos = pos
        self.ident = ident

    @property
    def message(self):
        return f'Повторная константа {self.ident}'


class UnannouncedConstant(SemanticError):
    def __init__(self, pos, ident):
        self.pos = pos
        self.ident = ident

    @property
    def message(self):
        return f'Необъявленная константа {self.ident}'


@dataclass
class TypeSpecifier(abc.ABC):
    def check(self):
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
    def check(self):
        pass


@dataclass
class ConstantExpression:
    expression: Expression

    def check(self):
        pass


@dataclass
class Enumerator:
    identifier: str
    constantExpression: ConstantExpression
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, constExpr = self
        idc, _ = coords
        return Enumerator(ident, constExpr, idc.start)

    def check(self, enum_pos):
        if check_const(self.identifier):
            raise RepeatedConstant(self.identifier_pos, self.identifier)
        if isinstance(self.constantExpression, ConstantExpression):
            add_to_const(self.identifier, self.constantExpression.expression)
            print("in Enumerator full содержит", self.constantExpression.expression)
            self.constantExpression.expression.check()

        # TODO: где-то тут у меня будет функция "посчитать константное выражение"
        else:
            add_to_const(self.identifier, enum_pos)


@dataclass
class StructOrUnionStatement(abc.ABC):
    def check(self):
        pass


@dataclass
class EmptyStructOrUnionStatement(StructOrUnionStatement):
    identifier: str
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, = self
        idc, = coords
        return EmptyStructOrUnionStatement(ident, idc.start)

    def check(self):
        if check_and_add_to_map(self.identifier, esu_tag):
            raise RepeatedTag(self.identifier_pos, self.identifier)


@dataclass
class StructOrUnionSpecifier(TypeSpecifier):
    type: str
    structOrUnionSpecifier: StructOrUnionStatement

    def check(self):
        self.structOrUnionSpecifier.check()


@dataclass
class EnumStatement(abc.ABC):
    def check(self):
        pass


@dataclass
class EnumTypeSpecifier(TypeSpecifier):
    enumStatement: EnumStatement

    def check(self):
        self.enumStatement.check()


@dataclass
class FullEnumStatement(EnumStatement):
    identifier: str
    enumeratorList: list[Enumerator]
    isEndComma: bool
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, enList, IsComma = self
        idc, _, enc, _, icc = coords
        return FullEnumStatement(ident, enList, IsComma, idc.start)

    def check(self):
        if len(self.identifier) != 0 and check_and_add_to_map(self.identifier, esu_tag):
            raise RepeatedTag(self.identifier_pos, self.identifier)

        for idx, enumerator in enumerate(self.enumeratorList):
            enumerator.check(idx)


@dataclass
class EmptyEnumStatement(EnumStatement):
    identifier: str
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, = self
        idc, = coords
        return EmptyEnumStatement(ident, idc.start)

    def check(self):
        if check_and_add_to_map(self.identifier, esu_tag):
            raise RepeatedTag(self.identifier_pos, self.identifier)


@dataclass
class IdentifierExpression(Expression):
    identifier: str
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, = self
        idc, = coords
        return IdentifierExpression(ident, idc.start)

    def check(self):
        if not check_const(self.identifier):
            raise UnannouncedConstant(self.identifier_pos, self.identifier)


@dataclass
class IntExpression(Expression):
    value: int

    def check(self):
        # Начало вычисления для expression
        print("IntExpression")


@dataclass
class BinaryOperationExpression(Expression):
    left: Expression
    operation: str
    right: Expression

    def check(self):
        print("BinaryOperationExpression")
        self.left.check()
        self.right.check()


@dataclass
class UnaryOperationExpression(Expression):
    operation: str
    expression: Expression

    def check(self):
        print("UnaryOperationExpression")
        self.expression.check()


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
        self.declarator.check()


@dataclass
class AbstractDeclaratorArrayList:
    arrays: list[AbstractDeclarator]

    def check(self):
        for ad in self.arrays:
            ad.check()


@dataclass
class AbstractDeclaratorArray(AbstractDeclarator):
    declarator: Expression

    def check(self):
        pass
        # print("abstractDeclaratorArray", self.declarator)
        # TODO: здесь уже смогу смотреть expression
        # И здесь проверяю 2 пункт


@dataclass
class AbstractDeclaratorsOpt:
    abstractDeclaratorList: list[AbstractDeclarator]

    def check(self):
        for ad in self.abstractDeclaratorList:
            ad.check()


@dataclass
class SimpleTypeSpecifier(TypeSpecifier):
    simpleType: SimpleType

    def check(self):
        # вроде не пригодится
        pass


@dataclass
class SizeofExpression(Expression):
    declarationBody: TypeSpecifier
    varName: AbstractDeclaratorsOpt

    def check(self):
        # TODO: до него надо еще дойти
        print("это sizeof")


@dataclass
class Declaration:
    declarationBody: TypeSpecifier
    varName: AbstractDeclaratorsOpt

    def check(self):
        self.varName.check()
        # print("in declarationBody", self.declarationBody)
        self.declarationBody.check()
        print()


@dataclass
class AbstractDeclaratorPrim(abc.ABC):
    pass


@dataclass
class AbstractDeclaratorPrimSimple(AbstractDeclaratorPrim):
    identifier: str
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, = self
        idc, = coords
        return AbstractDeclaratorPrimSimple(ident, idc.start)

    def check(self):
        if check_and_add_to_map(self.identifier, esu_ident):
            raise RepeatedIdentifier(self.identifier_pos, self.identifier)


@dataclass
class AbstractDeclaratorPrimDifficult(AbstractDeclaratorPrim):
    identifier: AbstractDeclarator

    def check(self):
        self.identifier.check()


@dataclass
class FullStructOrUnionStatement(StructOrUnionStatement):
    identifierOpt: str
    declarationList: list[Declaration]
    identifier_pos: pe.Position

    @pe.ExAction
    def create(self, coords, res_coord):
        ident, declList = self
        idc, _, dc, _ = coords
        return FullStructOrUnionStatement(ident, declList, idc.start)

    def check(self):
        if len(self.identifierOpt) != 0 and check_and_add_to_map(self.identifierOpt, esu_tag):
            raise RepeatedTag(self.identifier_pos, self.identifierOpt)

        for declaration in self.declarationList:
            declaration.check()


esu_ident = {}

esu_tag = {}


def check_and_add_to_map(s, in_map):
    if s in in_map:
        print("повтор")
        return True
    else:
        in_map[s] = True
        return False


const_name = {}


def check_const(ident):
    return ident in const_name


def add_to_const(ident, expr):
    const_name[ident] = expr


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

NFullEnumStatement |= NIdentifierOpt, '{', NEnumeratorList, NCommaOpt, '}', FullEnumStatement.create

NIdentifierOpt |= lambda: ""
NIdentifierOpt |= IDENTIFIER

NEnumStatement |= NEmptyEnumStatement, EmptyEnumStatement.create

NEmptyEnumStatement |= IDENTIFIER

NEnumeratorList |= NEnumerator, lambda e: [e]
NEnumeratorList |= NEnumeratorList, ',', NEnumerator, lambda el, e: el + [e]

NEnumerator |= IDENTIFIER, NEnumeratorExpressionOpt, Enumerator.create

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

NFactor |= IDENTIFIER, IdentifierExpression.create

NFactor |= KW_SIZEOF, '(', NTypeSpecifier, NAbstractDeclaratorsOpt, ')', SizeofExpression

NStructOrUnionSpecifier |= NStructOrUnion, NStructOrUnionStatement, StructOrUnionSpecifier

NStructOrUnion |= KW_STRUCT, lambda: "struct"
NStructOrUnion |= KW_UNION, lambda: "union"

NStructOrUnionStatement |= NFullStructOrUnionStatement
NStructOrUnionStatement |= NEmptyStructOrUnionStatement

NEmptyStructOrUnionStatement |= IDENTIFIER, EmptyStructOrUnionStatement.create

NFullStructOrUnionStatement |= NIdentifierOpt, '{', NDeclarationList, '}', FullStructOrUnionStatement.create


def main():
    p = pe.Parser(NProgram)

    assert p.is_lalr_one()

    p.add_skipped_domain('\\s')

    files = ["tests/sem_first.txt"]
    # files = ["tests/mixed.txt"]

    for filename in files:
        print("file:", filename)
        try:
            with open(filename) as f:
                tree = p.parse(f.read())
                pprint(tree)
                print()
                tree.check()
                print("Семантических ошибок не найдено")
        except pe.Error as e:
            print(f'Ошибка {e.pos}: {e.message}')
        except Exception as e:
            print(e)


main()

print("esuIdent:", esu_ident)
print("esuTag:", esu_tag)
print("const_name:", const_name)
