import abc
import enum
import parser_edsl as pe
from dataclasses import dataclass
from pprint import pprint


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
class AbstractDeclarator:
    pointer: str
    declarator: str
    arrayList: ListArraysOpt


@dataclass
class AbstractDeclaratorsOpt:
    abstractDeclaratorList: list[AbstractDeclarator]


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


@dataclass
class FullStructOrUnionStatement(StructOrUnionStatement):
    identifierOpt: str
    declarationList: list[Declaration]


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

NListArraysOpt = pe.NonTerminal('ListArraysOpt')
NModifySimpleType = pe.NonTerminal('ModifySimpleType')

NCommaOpt = pe.NonTerminal('CommaOpt')
NPointerOpt = pe.NonTerminal('PointerOpt')


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

KW_POINTER = pe.Terminal('POINTER', r'(\*)*', str)


NProgram |= NDeclarationList, Program

NDeclarationList |= lambda: []
NDeclarationList |= NDeclarationList, NDeclaration, lambda dl, d: dl + [d]

NDeclaration |= NTypeSpecifier, NAbstractDeclaratorsOpt, ';', Declaration

NAbstractDeclaratorsOpt |= lambda: []
NAbstractDeclaratorsOpt |= NAbstractDeclarators

NAbstractDeclarators |= NAbstractDeclarator, lambda a: [a]
NAbstractDeclarators |= NAbstractDeclarators, ',', NAbstractDeclarator, lambda ads, a: ads + [a]

NAbstractDeclarator |= NPointerOpt, IDENTIFIER, NListArraysOpt, AbstractDeclarator

NPointerOpt |= KW_POINTER
NPointerOpt |= lambda: ""


NListArraysOpt |= lambda: []
NListArraysOpt |= '[', NModifySimpleType, ']', NListArraysOpt, lambda mst, lao: lao + [mst]
NModifySimpleType |= IDENTIFIER
NModifySimpleType |= NSimpleType, lambda st: str(st)


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

NMultyOperation |= '∙', lambda: '∙'
NMultyOperation |= '/', lambda: '/'

NFactor |= '(', NExpression, ')'

NFactor |= INTEGER, IntExpression

NFactor |= IDENTIFIER, IdentifierExpression

NFactor |= KW_SIZEOF, '(', NTypeSpecifier, NAbstractDeclaratorsOpt, ')', SizeofExpression


# тут может возникнуть проблема с конструктором
NStructOrUnionSpecifier |= NStructOrUnion, NStructOrUnionStatement, StructOrUnionSpecifier

NStructOrUnion |= KW_STRUCT, lambda: "struct"
NStructOrUnion |= KW_UNION, lambda: "union"

NStructOrUnionStatement |= NFullStructOrUnionStatement
NStructOrUnionStatement |= NEmptyStructOrUnionStatement

NEmptyStructOrUnionStatement |= IDENTIFIER, EmptyStructOrUnionStatement

NFullStructOrUnionStatement |= NIdentifierOpt,  '{', NDeclarationList, '}', FullStructOrUnionStatement


def main():
    p = pe.Parser(NProgram)

    p.print_table()

    assert p.is_lalr_one()

    p.add_skipped_domain('\\s')

    files = ["tests/mixed.txt"]

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
