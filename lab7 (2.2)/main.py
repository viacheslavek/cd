import abc
import enum
import parser_edsl as pe
import re
import typing
from dataclasses import dataclass
from pprint import pprint


# TODO: пока делаю все для enum
# TODO: делаю так же, как в примере, но по своему синтаксису
class SimpleType(enum.Enum):
    Char = "CHAR"
    Short = "SHORT"
    Int = "INT"
    Long = "LONG"
    Float = "FLOAT"
    Double = "DOUBLE"
    Signed = "SIGNED"
    Unsigned = "UNSIGNED"


class Declaration(abc.ABC):
    pass


@dataclass
class BlockDeclarationSpecifier(Statement):
    body: list[Statement]


# TODO: делаю так же, как в примере, но по своему синтаксису
NProgram = map(pe.NonTerminal, "Program")


def main():
    p = pe.Parser(NProgram)
    assert p.is_lalr_one()

    p.add_skipped_domain('\\s')

    files = ["tests/enum.txt"]

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
