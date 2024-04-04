import abc
import enum
import parser_edsl as pe
import re
import typing
from dataclasses import dataclass
from pprint import pprint


# class SimpleType(enum.Enum):
#     Char = "CHAR"
#     Short = "SHORT"
#     Int = "INT"
#     Long = "LONG"
#     Float = "FLOAT"
#     Double = "DOUBLE"
#     Signed = "SIGNED"
#     Unsigned = "UNSIGNED"
#
#
# @dataclass
# class CommaOpt(enum.Enum):
#     Comma = ","
#     EmptyComma = ""
#
#
# @dataclass
# class IdentifierOpt(abc.ABC):
#     pass
#
#
# @dataclass
# class Identifier(IdentifierOpt):
#     ident: str
#
#
# @dataclass
# class IdentifierEpsilon(IdentifierOpt):
#     pass
#
#
# @dataclass
# class ConstantExpression:
#     constantExpression: int
#
#
# @dataclass
# class EnumeratorConstant:
#     ident: str
#
#
# @dataclass
# class EnumeratorExpressionOpt(abc.ABC):
#     pass
#
#
# @dataclass
# class ConstantExpression(EnumeratorExpressionOpt):
#     pass
#
#
# @dataclass
# class EmptyExpression(EnumeratorExpressionOpt):
#     pass
#
#
# @dataclass
# class Enumerator:
#     enumeratorConstant: EnumeratorConstant
#     enumeratorExpressionOpt: EnumeratorExpressionOpt
#
#
# @dataclass
# class EnumeratorList:
#     enList: list[Enumerator]
#
#
@dataclass
class TypeSpecifier(abc.ABC):
    pass
#
#
# @dataclass
# class SimpleTypeSpecifier(TypeSpecifier):
#     simpleType: SimpleType


@dataclass
class EnumStatement(abc.ABC):
    pass
#
#
# @dataclass
# class FullEnumStatement(EnumStatement):
#     identOpt: IdentifierOpt
#     enumList: EnumeratorList
#     commaOpt: CommaOpt


@dataclass
class EmptyEnumStatement(EnumStatement):
    ident: str


@dataclass
class EnumType:
    enumStatement: EnumStatement


@dataclass
class EnumTypeSpecifier(TypeSpecifier):
    enumType: EnumType
#
#
# @dataclass
# class AbstractDeclaratorOpt(abc.ABC):
#     pass
#
#
# @dataclass
# class AbstractDeclarator(AbstractDeclaratorOpt):
#     ident: str
#
#
# @dataclass
# class AbstractEpsilon(AbstractDeclaratorOpt):
#     pass


@dataclass
class Declaration:
    typeSpecifier: TypeSpecifier
    # abstractDeclaratorOpt: AbstractDeclaratorOpt


@dataclass
class DeclarationList:
    declarations: list[Declaration]


INTEGER = pe.Terminal('INTEGER', '[0-9]+', int, priority=7)
IDENT = pe.Terminal('IDENT', '[A-Za-z][A-Za-z0-9]*', str.upper)


def make_keyword(t):
    return pe.Terminal(t, t, lambda name: None, priority=10)


KW_ENUM = map(make_keyword, 'enum'.split())

# KW_CHAR, KW_SHORT, KW_INT, KW_LONG, KW_FLOAT, KW_DOUBLE, KW_SIGNED, KW_UNSIGNED = \
#     map(make_keyword, 'char short int long float double signed unsigned'.split())
#
#
NDeclarationList, NDeclaration, NTypeSpecifier, NAbstractDeclaratorOpt = \
    map(pe.NonTerminal, 'DeclarationList Declaration TypeSpecifier AbstractDeclaratorOpt'.split())

NAbstractDeclarator, NAbstractEpsilon, NSimpleTypeSpecifier, NSimpleType = \
    map(pe.NonTerminal, 'AbstractDeclarator AbstractEpsilon SimpleTypeSpecifier SimpleType'.split())

NEnumTypeSpecifier, NEnumType, NEnumStatement, NFullEnumStatement, NEmptyEnumStatement = \
    map(pe.NonTerminal, 'EnumTypeSpecifier EnumType EnumStatement FullEnumStatement EmptyEnumStatement'.split())
#
# NIdentifierOpt, NIdentifier, NIdentifierEpsilon, NCommaOpt = \
#     map(pe.NonTerminal, 'IdentifierOpt Identifier IdentifierEpsilon CommaOpt'.split())
#
# NEnumeratorList, NEnumerator, NEnumeratorConstant = \
#     map(pe.NonTerminal, 'EnumeratorList Enumerator EnumeratorConstant'.split())
#
# NEnumeratorExpressionOpt, NConstantExpression, NEmptyExpression = \
#     map(pe.NonTerminal, 'EnumeratorExpressionOpt ConstantExpression EmptyExpression'.split())
#

# TODO: делаю переходы

NDeclarationList |= NDeclarationList, NDeclaration, lambda dlist, d: dlist + [d]
NDeclarationList |= NDeclaration, lambda d: [d]

# NDeclaration |= NTypeSpecifier, NAbstractDeclaratorOpt, ';'

NDeclaration |= NTypeSpecifier, ';'


# NTypeSpecifier |= NSimpleTypeSpecifier
NTypeSpecifier |= NEnumTypeSpecifier

# NSimpleTypeSpecifier |= NSimpleType
#
#
# NAbstractDeclaratorOpt |= NAbstractDeclarator
# NAbstractDeclaratorOpt |= NAbstractEpsilon
#
# NAbstractDeclarator |= IDENT
#

# NSimpleType |= KW_CHAR, lambda: SimpleType.Char
# NSimpleType |= KW_SHORT, lambda: SimpleType.Short
# NSimpleType |= KW_INT, lambda: SimpleType.Int
# NSimpleType |= KW_LONG, lambda: SimpleType.Long
# NSimpleType |= KW_FLOAT, lambda: SimpleType.Float
# NSimpleType |= KW_DOUBLE, lambda: SimpleType.Double
# NSimpleType |= KW_SIGNED, lambda: SimpleType.Signed
# NSimpleType |= KW_UNSIGNED, lambda: SimpleType.Unsigned

NEnumTypeSpecifier |= NEnumType

# NEnumType |= KW_ENUM, NEnumStatement

# NEnumStatement |= NFullEnumStatement
# NEnumStatement |= NEmptyEnumStatement
#
# NFullEnumStatement |= NIdentifierOpt, '{', NEnumeratorList, NCommaOpt, '}'
#
# NEmptyEnumStatement |= IDENT
#
# NCommaOpt |= ','
# NCommaOpt |= ""  # мб это плохо
#
# NEnumeratorList |= EnumeratorList, ',', NEnumerator, lambda elist, e: elist + [e]
# NEnumeratorList |= NEnumerator, lambda e: [e]
#
# NEnumerator |= NEnumeratorConstant, NEnumeratorExpressionOpt
#
# NEnumeratorConstant |= IDENT
#
# NIdentifierOpt |= IDENT
# NIdentifierOpt |= NIdentifierEpsilon
#
# NEnumeratorExpressionOpt |= NConstantExpression
# NEnumeratorExpressionOpt |= NEmptyExpression
#
# ConstantExpression |= INTEGER


def main():
    p = pe.Parser(NDeclarationList)

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
