import abc
import enum
import parser_edsl as pe
import re
import typing
from dataclasses import dataclass
from pprint import pprint


# Еще сильнее упрощаю грамматику
#
# DeclarationList -> Declaration | DeclarationList Declaration   // готово
#
# Declaration -> EnumStatement VarnameOpt ;
#
# // пока пустой enum
#
# EnumStatement -> ENUM
#
# VarnameOpt -> Varname | epsi


@dataclass
class TypeSpecifier(abc.ABC):
    pass


@dataclass
class EnumTypeSpecifier(TypeSpecifier):
    pass


@dataclass
class SimpleTypeSpecifier(TypeSpecifier):
    pass


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

NEnumStatement = pe.NonTerminal('EnumStatement')

NAbstractDeclaratorOpt = pe.NonTerminal('AbstractDeclaratorOpt')
NAbstractDeclarator = pe.NonTerminal('AbstractDeclarator')

NTypeSpecifier = pe.NonTerminal('TypeSpecifier')
NEnumTypeSpecifier = pe.NonTerminal('EnumTypeSpecifier')
NSimpleTypeSpecifier = pe.NonTerminal('SimpleTypeSpecifier')


def make_keyword(image):
    return pe.Terminal(image, image, lambda _: None, priority=10)


ENUM = make_keyword('enum')
INT = make_keyword('int')

IDENT = pe.Terminal('IDENT', '[A-Za-z_][A-Za-z_0-9]*', str)

NProgram |= NDeclarationList, Program

NDeclarationList |= lambda: []
NDeclarationList |= NDeclarationList, NDeclaration, lambda dl, d: dl + [d]

NDeclaration |= NTypeSpecifier, NAbstractDeclaratorOpt, ';', Declaration

NAbstractDeclaratorOpt |= lambda: ""
NAbstractDeclaratorOpt |= NAbstractDeclarator

NAbstractDeclarator |= IDENT

NTypeSpecifier |= NEnumTypeSpecifier, EnumTypeSpecifier
NTypeSpecifier |= NSimpleTypeSpecifier, SimpleTypeSpecifier

NSimpleTypeSpecifier |= INT

NEnumTypeSpecifier |= NEnumStatement

NEnumStatement |= ENUM


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
