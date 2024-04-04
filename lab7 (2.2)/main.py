import abc
import enum
import parser_edsl as pe
from dataclasses import dataclass
from pprint import pprint


@dataclass
class ConstantExpression:
    body: str


@dataclass
class Enumerator:
    identifier: str
    constantExpression: ConstantExpression


@dataclass
class EnumeratorList:
    body: list[Enumerator]


@dataclass
class TypeSpecifier(abc.ABC):
    pass


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


class SimpleType(enum.Enum):
    Char = "CHAR"
    Short = "SHORT"
    Int = "INT"
    Long = "LONG"
    Float = "FLOAT"
    Double = "DOUBLE"
    Signed = "SIGNED"
    Unsigned = "UNSIGNED"


@dataclass
class AbstractDeclarator:
    pointer: str
    declarator: str


@dataclass
class AbstractDeclaratorsOpt:
    abstractDeclaratorList: list[AbstractDeclarator]


@dataclass
class SimpleTypeSpecifier(TypeSpecifier):
    simpleType: SimpleType


@dataclass
class Declaration:
    declarationBody: TypeSpecifier
    varName: AbstractDeclaratorsOpt


@dataclass
class Program:
    declarationList: list[Declaration]


NProgram = pe.NonTerminal('Program')

NDeclarationList = pe.NonTerminal('DeclarationList')
NDeclaration = pe.NonTerminal('Declaration')

NAbstractDeclaratorsOpt = pe.NonTerminal('AbstractDeclaratorsOpt')
NAbstractDeclarators = pe.NonTerminal('AbstractDeclarators')
NAbstractDeclarator = pe.NonTerminal('AbstractDeclarator')

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
NPointerOpt = pe.NonTerminal('PointerOpt')


def make_keyword(image):
    return pe.Terminal(image, image, lambda _: None, priority=10)


KW_ENUM = make_keyword('enum')

KW_CHAR, KW_SHORT, KW_INT, KW_LONG, KW_FLOAT, KW_DOUBLE, KW_SIGNED, KW_UNSIGNED = \
    map(make_keyword, 'char short int long float double signed unsigned'.split())

KW_IDENTIFIER = pe.Terminal('IDENTIFIER', r'[A-Za-z_]([A-Za-z_0-9])*', str)

KW_POINTER = pe.Terminal('POINTER', r'(\*)*', str)


NProgram |= NDeclarationList, Program

NDeclarationList |= lambda: []
NDeclarationList |= NDeclarationList, NDeclaration, lambda dl, d: dl + [d]

NDeclaration |= NTypeSpecifier, NAbstractDeclaratorsOpt, ';', Declaration

NAbstractDeclaratorsOpt |= lambda: []
NAbstractDeclaratorsOpt |= NAbstractDeclarators

NAbstractDeclarators |= NAbstractDeclarator, lambda a: [a]
NAbstractDeclarators |= NAbstractDeclarators, ',', NAbstractDeclarator, lambda ads, a: ads + [a]

NAbstractDeclarator |= NPointerOpt, KW_IDENTIFIER, AbstractDeclarator

NPointerOpt |= KW_POINTER
NPointerOpt |= lambda: ""

NTypeSpecifier |= NEnumTypeSpecifier, EnumTypeSpecifier
NTypeSpecifier |= NSimpleTypeSpecifier, SimpleTypeSpecifier

NSimpleTypeSpecifier |= NSimpleType

NSimpleType |= KW_CHAR, lambda: SimpleType.Char
NSimpleType |= KW_SHORT, lambda: SimpleType.Short
NSimpleType |= KW_INT, lambda: SimpleType.Int
NSimpleType |= KW_LONG, lambda: SimpleType.Long
NSimpleType |= KW_FLOAT, lambda: SimpleType.Float
NSimpleType |= KW_DOUBLE, lambda: SimpleType.Double
NSimpleType |= KW_SIGNED, lambda: SimpleType.Signed
NSimpleType |= KW_UNSIGNED, lambda: SimpleType.Unsigned


NEnumTypeSpecifier |= KW_ENUM, NEnumStatement

NEnumStatement |= NFullEnumStatement

NFullEnumStatement |= NIdentifierOpt, '{', NEnumeratorList, NCommaOpt, '}', FullEnumStatement

NIdentifierOpt |= lambda: ""
NIdentifierOpt |= KW_IDENTIFIER

NEnumStatement |= NEmptyEnumStatement, EmptyEnumStatement

NEmptyEnumStatement |= KW_IDENTIFIER

NEnumeratorList |= NEnumerator, lambda e: [e]
NEnumeratorList |= NEnumeratorList, ',', NEnumerator, lambda el, e: el + [e]

NEnumerator |= KW_IDENTIFIER, NEnumeratorExpressionOpt, Enumerator

NEnumeratorExpressionOpt |= '=', NConstantExpression
NEnumeratorExpressionOpt |= lambda: ""

# TODO: потом тут будет вычисление выражения
NConstantExpression |= KW_IDENTIFIER

NCommaOpt |= ',', lambda: True
NCommaOpt |= lambda: False


def main():
    p = pe.Parser(NProgram)

    p.print_table()

    assert p.is_lalr_one()

    p.add_skipped_domain('\\s')

    files = ["tests/enum_prev.txt"]

    for filename in files:
        print("file:", filename)
        try:
            with open(filename) as f:
                tree = p.parse(f.read())
                pprint(tree)
        except pe.Error as e:
            print(f'Ошибка {e.pos}: {e.message}')
        except Exception as e:
            print(e)


main()
