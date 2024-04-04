import abc
import enum
import parser_edsl as pe
import re
import typing
from dataclasses import dataclass
from pprint import pprint


@dataclass
class EnumeratorList:
    body: str


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
    endComma: bool


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
class SimpleTypeSpecifier(TypeSpecifier):
    simpleType: SimpleType


@dataclass
class Declaration:
    declarationBody: TypeSpecifier
    varName: str


@dataclass
class Program:
    declarationList: list[Declaration]


NProgram = pe.NonTerminal('Program')

NDeclarationList = pe.NonTerminal('DeclarationList')
NDeclaration = pe.NonTerminal('Declaration')

NAbstractDeclaratorOpt = pe.NonTerminal('AbstractDeclaratorOpt')
NAbstractDeclarator = pe.NonTerminal('AbstractDeclarator')

NTypeSpecifier = pe.NonTerminal('TypeSpecifier')

NEnumTypeSpecifier = pe.NonTerminal('EnumTypeSpecifier')

NSimpleTypeSpecifier = pe.NonTerminal('SimpleTypeSpecifier')
NSimpleType = pe.NonTerminal('SimpleType')

NEnumStatement = pe.NonTerminal('EnumStatement')

NFullEnumStatement = pe.NonTerminal('FullEnumStatement')
NEmptyEnumStatement = pe.NonTerminal('EmptyEnumStatement')

NEnumeratorList = pe.NonTerminal('EnumeratorList')

NIdentifierOpt = pe.NonTerminal('IdentifierOpt')

NCommaOpt = pe.NonTerminal('CommaOpt')


def make_keyword(image):
    return pe.Terminal(image, image, lambda _: None, priority=10)


KW_ENUM = make_keyword('enum')

KW_CHAR, KW_SHORT, KW_INT, KW_LONG, KW_FLOAT, KW_DOUBLE, KW_SIGNED, KW_UNSIGNED = \
    map(make_keyword, 'char short int long float double signed unsigned'.split())

KW_IDENTIFIER = pe.Terminal('IDENT', '[A-Za-z_][A-Za-z_0-9]*', str)


NProgram |= NDeclarationList, Program

NDeclarationList |= lambda: []
NDeclarationList |= NDeclarationList, NDeclaration, lambda dl, d: dl + [d]

NDeclaration |= NTypeSpecifier, NAbstractDeclaratorOpt, ';', Declaration

NAbstractDeclaratorOpt |= lambda: ""
NAbstractDeclaratorOpt |= NAbstractDeclarator

NAbstractDeclarator |= KW_IDENTIFIER

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

NEnumeratorList |= KW_IDENTIFIER

NCommaOpt |= ',', lambda: False
NCommaOpt |= lambda: True


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
