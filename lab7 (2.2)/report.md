% Лабораторная работа № 2.2 «Абстрактные синтаксические деревья»
% 10 апреля 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является получение навыков составления грамматик и проектирования синтаксических деревьев.

# Индивидуальный вариант
Определения структур, объединений и перечислений языка Си.
В инициализаторах перечислений допустимы знаки операций +, -, *, /, sizeof,
операндами могут служить имена перечислимых значений и целые числа.

Числовые константы могут быть только целочисленными и десятичными.

```
struct Coords {
  int x, y;
};

enum Color {
  COLOR_RED = 1,
  COLOR_GREEN = 2,
  COLOR_BLUE = 4,
  COLOR_HIGHLIGHT = 8, // запятая после последнего необязательна
};

enum ScreenType {
  SCREEN_TYPE_TEXT,
  SCREEN_TYPE_GRAPHIC
} screen_type;  // ← объявили переменную

enum {
  HPIXELS = 480,
  WPIXELS = 640,
  HCHARS = 24,
  WCHARS = 80,
};

struct ScreenChar {
  char symbol;
  enum Color sym_color;
  enum Color back_color;
};

struct TextScreen {
  struct ScreenChar chars[HCHARS][WCHARS];
};

struct GrahpicScreen {
  enum Color pixels[HPIXELS][WPIXELS];
};

union Screen {
  struct TextScreen text;
  struct GraphicScreen graphic;
};

enum {
  BUFFER_SIZE = sizeof(union Screen),
  PAGE_SIZE = 4096,
  PAGES_FOR_BUFFER = (BUFFER_SIZE + PAGE_SIZE - 1) / PAGE_SIZE
};

/* допустимы и вложенные определения */
struct Token {
  struct Fragment {
    struct Pos {
      int line;
      int col;
    } starting, following;
  } fragment;

  enum { Ident, IntConst, FloatConst } type;

  union {
    char *name;
    int int_value;
    double float_value;
  } info;
};

struct List {
  struct Token value;
  struct List *next;
};

```

В структурах, объединениях и перечислениях могут отсутствовать тег (имя) типа,
перечисление полей и объявляемые переменные. При выполнении этой лабораторной допустимо считать,
что все три компонента могут отсутствовать одновременно, т.е. следующий код будет верным:

```
struct;
union *p;
enum x[];
```
Это существенно сократит описание синтаксиса.


# Реализация

Во время построения синтаксиса я обращался к стандарту ISO/IEC 9899:1999 (E), где почерпнул некоторые идеи.

## Абстрактный синтаксис

Программа - это последовательность объявлений

```
Program -> DeclarationList

DeclarationList -> Declaration*
```

Объявление состоит из типа и возможных сложных имен идентификаторов

```
Declaration -> TypeSpecifier AbstractDeclaratorsOpt ;

AbstractDeclaratorsOpt -> (AbstractDeclarators)?

AbstractDeclarators -> AbstractDeclarator (, AbstractDeclarator)*

AbstractDeclaratorOpt -> (AbstractDeclarator)?
```

Сложное имя идентификатора может быть указателем, массивом или указателей массивов, или массивом указателей.
Так же могут быть скобки в объявлении. Или в любом случае есть имя идентификатора,
если AbstractDeclarator не пуст

```
AbstractDeclarator -> AbstractDeclaratorStar | AbstractDeclaratorArrayListOpt
AbstractDeclaratorStar -> STAR AbstractDeclarator

AbstractDeclaratorArrayListOpt -> (AbstractDeclaratorArrayList)?
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayListOpt AbstractDeclaratorArray

AbstractDeclaratorArray -> AbstractDeclaratorArray [ Expression ] | AbstractDeclaratorPrim

AbstractDeclaratorPrim -> AbstractDeclaratorPrimSimple | AbstractDeclaratorPrimDifficult

AbstractDeclaratorPrimSimple -> IDENTIFIER
AbstractDeclaratorPrimDifficult -> ( AbstractDeclarator )
```

Тип может быть enum, struct, union или простым типом Си

```
TypeSpecifier -> SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier

SimpleTypeSpecifier -> SimpleType

SimpleType -> CHAR | SHORT | INT | LONG | FLOAT | DOUBLE | SIGNED | UNSIGNED
```

Разберем то, как работает enum. 
Есть ключевое слово enum, после этого может быть его полная или неполная форма
(немного усложнил грамматику, можно было это не учитывать по условию задания)
Полная форма - возможное имя идентификатора и выражения в '{' ... '}'

Особое внимание надо уделить запятой, которая может как стоять,
так и не стоять в конце последнего выражения внутри скобок.

Внутри выражения есть идентификатор и возможное константное выражение.

```
EnumTypeSpecifier -> ENUM EnumStatement

EnumStatement -> FullEnumStatement | EmptyEnumStatement

FullEnumStatement -> IdentifierOpt { EnumeratorList CommaOpt }

IdentifierOpt -> (IDENTIFIER)?

EmptyEnumStatement -> IDENTIFIER

EnumeratorList -> Enumerator (, Enumerator)*

Enumerator -> IDENTIFIER EnumeratorExpressionOpt

EnumeratorExpressionOpt -> (= ConstantExpression)?

CommaOpt -> (,)?
```

Константное выражение вычисляется по стандартным математическим выражениям,
но с добавлением ключевого слова sizeof в унарных операторах

```
ConstantExpression -> Expression

Expression -> IDENTIFIER | INT | Expression BinOp Expression | UnOp Expression
BinaryOperation -> + | - | * | /
UnaryOperation -> + | - | sizeof ( TypeSpecifier AbstractDeclaratorsOpt )
```

Структуры и объединения похожи по структуре и отличаются ключевым словом.

```
StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement

StructOrUnion -> STRUCT | UNION

StructOrUnionStatement -> FullStructOrUnionStatement | EmptyStructOrUnionStatement

EmptyStructOrUnionStatement -> IDENTIFIER
```

Полное выражение может иметь внутри себя любые вложенные объявления, это показывается в синтаксисе.

```
FullStructOrUnionStatement -> IdentifierOpt { DeclarationList }
```

## Лексическая структура и конкретный синтаксис
```
Program -> DeclarationList

DeclarationList -> Declaration | DeclarationList Declaration

Declaration -> TypeSpecifier AbstractDeclaratorsOpt ;


AbstractDeclaratorsOpt -> AbstractDeclarators | ε

AbstractDeclarators -> AbstractDeclarators , AbstractDeclarator | AbstractDeclarator

AbstractDeclaratorOpt -> AbstractDeclarator | ε

AbstractDeclarator -> AbstractDeclaratorStar | AbstractDeclaratorArrayListOpt
AbstractDeclaratorStar -> STAR AbstractDeclarator


AbstractDeclaratorArrayListOpt -> AbstractDeclaratorArrayList | ε
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayListOpt AbstractDeclaratorArray

AbstractDeclaratorArray -> AbstractDeclaratorArray [ Expression ] | AbstractDeclaratorPrim

AbstractDeclaratorPrim -> AbstractDeclaratorPrimSimple | AbstractDeclaratorPrimDifficult

AbstractDeclaratorPrimSimple -> IDENTIFIER
AbstractDeclaratorPrimDifficult -> ( AbstractDeclarator )


TypeSpecifier -> SimpleTypeSpecifier | EnumTypeSpecifier


SimpleTypeSpecifier -> SimpleType

SimpleType -> CHAR | SHORT | INT | LONG | FLOAT | DOUBLE | SIGNED | UNSIGNED


EnumTypeSpecifier -> ENUM EnumStatement

EnumStatement -> FullEnumStatement | EmptyEnumStatement

FullEnumStatement -> IdentifierOpt { EnumeratorList CommaOpt }

IdentifierOpt -> IDENTIFIER | ε

EmptyEnumStatement -> IDENTIFIER

EnumeratorList -> EnumeratorList , Enumerator | Enumerator

Enumerator -> IDENTIFIER EnumeratorExpressionOpt

EnumeratorExpressionOpt -> = ConstantExpression | ε


CommaOpt -> , | ε


ConstantExpression -> Expression

Expression -> ArithmeticExpression

ArithmeticExpression -> Term | + Term | - Term | ArithmeticExpression AddOperation Term

AddOperation -> + | -

Term -> Factor | Term MultyOperation Factor
MultyOperation -> * | /

Factor -> sizeof ( TypeSpecifier AbstractDeclaratorsOpt ) | IDENTIFIER | INT | ( Expression )


StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement

StructOrUnion -> STRUCT | UNION

StructOrUnionStatement -> FullStructOrUnionStatement | EmptyStructOrUnionStatement

EmptyStructOrUnionStatement -> IDENTIFIER

FullStructOrUnionStatement -> IdentifierOpt { DeclarationList }

```

```
INT = [0-9]*
IDENTIFIER = [A-Za-z_]([A-Za-z_0-9])*
```

## Программная реализация

```python
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
class AbstractDeclarator(abc.ABC):
    pass


@dataclass
class AbstractDeclaratorPointer(AbstractDeclarator):
    declarator: str


@dataclass
class AbstractDeclaratorArrayList:
    arrays: list[str]


@dataclass
class AbstractDeclaratorArray(AbstractDeclarator):
    declarator: str


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
class AbstractDeclaratorPrim(abc.ABC):
    pass


@dataclass
class AbstractDeclaratorPrimSimple(AbstractDeclaratorPrim):
    identifier: str


@dataclass
class AbstractDeclaratorPrimDifficult(AbstractDeclaratorPrim):
    identifier: AbstractDeclarator


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

NAbstractDeclaratorsOpt |= lambda: []
NAbstractDeclaratorsOpt |= NAbstractDeclarators

NAbstractDeclarators |= NAbstractDeclarator, lambda a: [a]
NAbstractDeclarators |= NAbstractDeclarators, ',', NAbstractDeclarator, lambda ads, a: ads + [a]

NAbstractDeclarator |= NAbstractDeclaratorStar
NAbstractDeclarator |= NAbstractDeclaratorArrayList, AbstractDeclaratorArrayList

NAbstractDeclaratorStar |= '*', NAbstractDeclarator, AbstractDeclaratorPointer

NAbstractDeclaratorArrayList |= NAbstractDeclaratorArray, lambda a: [a]
NAbstractDeclaratorArrayList |= (NAbstractDeclaratorArrayList, NAbstractDeclaratorArray,
                                 lambda adalo, a: adalo + [a])

NAbstractDeclaratorArray |= '[', NExpression, ']', AbstractDeclaratorArray

NAbstractDeclaratorArray |= NAbstractDeclaratorPrim
NAbstractDeclaratorPrim |= NAbstractDeclaratorPrimSimple, AbstractDeclaratorPrimSimple

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

```

# Тестирование

## Входные данные

```
struct Coords {
  int x, y;
};

enum ScreenType { aaa = bab } *screenType[5 + 5], **fas;

enum ScreenType { aaa = bab } *screenType[5 + 5], **fas;

enum ScreenType { aaa = bab } a[1], *a[1], (*a)[1];

enum Color {
  COLOR_RED = 1,
  COLOR_GREEN = 2,
  COLOR_BLUE = 2*2,
  COLOR_HIGHLIGHT = 8,
};

enum ScreenType {
  SCREEN_TYPE_TEXT,
  SCREEN_TYPE_GRAPHIC
} screen_type;

enum {
  HPIXELS = 480,
  WPIXELS = 640,
  HCHARS = 24,
  WCHARS = 80,
};

struct ScreenChar {
  char symbol;
  enum Color sym_color;
  enum Color back_color;
};

struct TextScreen {
  struct ScreenChar chars[HCHARS][WCHARS];
};

struct GrahpicScreen {
  enum Color pixels[HPIXELS][WPIXELS];
};

union Screen {
  struct TextScreen text;
  struct GraphicScreen graphic;
};

enum {
  BUFFER_SIZE = sizeof(union Screen),
  PAGE_SIZE = 4096,
  PAGES_FOR_BUFFER = (BUFFER_SIZE + PAGE_SIZE - 1) / PAGE_SIZE
};

struct Token {
  struct Fragment {
    struct Pos {
      int line;
      int col;
    } starting, following;
  } fragment;

  enum { Ident, IntConst, FloatConst } type;

  union {
    char *name;
    int int_value;
    double float_value;
  } info;
};

struct List {
  struct Token value;
  struct List *next;
};

```

## Вывод на `stdout`

<!-- ENABLE LONG LINES -->

```
file: tests/mixed.txt
Program(declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='Coords',
                                                                                                                              declarationList=[Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Int: 'INT'>),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='x')]),
                                                                                                                                                                    AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='y')])])])),
                                     varName=[]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='ScreenType',
                                                                                                       enumeratorList=[Enumerator(identifier='aaa',
                                                                                                                                  constantExpression=ConstantExpression(expression=IdentifierExpression(identifier='bab')))],
                                                                                                       isEndComma=False)),
                                     varName=[AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='screenType'),
                                                                                                                       AbstractDeclaratorArray(declarator=BinaryOperationExpression(left=IntExpression(value='5'),
                                                                                                                                                                                    operation='+',
                                                                                                                                                                                    right=IntExpression(value='5')))])),
                                              AbstractDeclaratorPointer(declarator=AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='fas')])))]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='ScreenType',
                                                                                                       enumeratorList=[Enumerator(identifier='aaa',
                                                                                                                                  constantExpression=ConstantExpression(expression=IdentifierExpression(identifier='bab')))],
                                                                                                       isEndComma=False)),
                                     varName=[AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='screenType'),
                                                                                                                       AbstractDeclaratorArray(declarator=BinaryOperationExpression(left=IntExpression(value='5'),
                                                                                                                                                                                    operation='+',
                                                                                                                                                                                    right=IntExpression(value='5')))])),
                                              AbstractDeclaratorPointer(declarator=AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='fas')])))]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='ScreenType',
                                                                                                       enumeratorList=[Enumerator(identifier='aaa',
                                                                                                                                  constantExpression=ConstantExpression(expression=IdentifierExpression(identifier='bab')))],
                                                                                                       isEndComma=False)),
                                     varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='a'),
                                                                                  AbstractDeclaratorArray(declarator=IntExpression(value='1'))]),
                                              AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='a'),
                                                                                                                       AbstractDeclaratorArray(declarator=IntExpression(value='1'))])),
                                              AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimDifficult(identifier=AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='a')]))),
                                                                                  AbstractDeclaratorArray(declarator=IntExpression(value='1'))])]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='Color',
                                                                                                       enumeratorList=[Enumerator(identifier='COLOR_RED',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='1'))),
                                                                                                                       Enumerator(identifier='COLOR_GREEN',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='2'))),
                                                                                                                       Enumerator(identifier='COLOR_BLUE',
                                                                                                                                  constantExpression=ConstantExpression(expression=BinaryOperationExpression(left=IntExpression(value='2'),
                                                                                                                                                                                                             operation='*',
                                                                                                                                                                                                             right=IntExpression(value='2')))),
                                                                                                                       Enumerator(identifier='COLOR_HIGHLIGHT',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='8')))],
                                                                                                       isEndComma=True)),
                                     varName=[]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='ScreenType',
                                                                                                       enumeratorList=[Enumerator(identifier='SCREEN_TYPE_TEXT',
                                                                                                                                  constantExpression=''),
                                                                                                                       Enumerator(identifier='SCREEN_TYPE_GRAPHIC',
                                                                                                                                  constantExpression='')],
                                                                                                       isEndComma=False)),
                                     varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='screen_type')])]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='',
                                                                                                       enumeratorList=[Enumerator(identifier='HPIXELS',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='480'))),
                                                                                                                       Enumerator(identifier='WPIXELS',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='640'))),
                                                                                                                       Enumerator(identifier='HCHARS',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='24'))),
                                                                                                                       Enumerator(identifier='WCHARS',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='80')))],
                                                                                                       isEndComma=True)),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='ScreenChar',
                                                                                                                              declarationList=[Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Char: 'CHAR'>),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='symbol')])]),
                                                                                                                                               Declaration(declarationBody=EnumTypeSpecifier(enumStatement=EmptyEnumStatement(identifier='Color')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='sym_color')])]),
                                                                                                                                               Declaration(declarationBody=EnumTypeSpecifier(enumStatement=EmptyEnumStatement(identifier='Color')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='back_color')])])])),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='TextScreen',
                                                                                                                              declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='ScreenChar')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='chars'),
                                                                                                                                                                                                        AbstractDeclaratorArray(declarator=IdentifierExpression(identifier='HCHARS')),
                                                                                                                                                                                                        AbstractDeclaratorArray(declarator=IdentifierExpression(identifier='WCHARS'))])])])),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='GrahpicScreen',
                                                                                                                              declarationList=[Declaration(declarationBody=EnumTypeSpecifier(enumStatement=EmptyEnumStatement(identifier='Color')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='pixels'),
                                                                                                                                                                                                        AbstractDeclaratorArray(declarator=IdentifierExpression(identifier='HPIXELS')),
                                                                                                                                                                                                        AbstractDeclaratorArray(declarator=IdentifierExpression(identifier='WPIXELS'))])])])),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='union',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='Screen',
                                                                                                                              declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='TextScreen')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='text')])]),
                                                                                                                                               Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='GraphicScreen')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='graphic')])])])),
                                     varName=[]),
                         Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='',
                                                                                                       enumeratorList=[Enumerator(identifier='BUFFER_SIZE',
                                                                                                                                  constantExpression=ConstantExpression(expression=SizeofExpression(declarationBody=StructOrUnionSpecifier(type='union',
                                                                                                                                                                                                                                           structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='Screen')),
                                                                                                                                                                                                    varName=[]))),
                                                                                                                       Enumerator(identifier='PAGE_SIZE',
                                                                                                                                  constantExpression=ConstantExpression(expression=IntExpression(value='4096'))),
                                                                                                                       Enumerator(identifier='PAGES_FOR_BUFFER',
                                                                                                                                  constantExpression=ConstantExpression(expression=BinaryOperationExpression(left=BinaryOperationExpression(left=BinaryOperationExpression(left=IdentifierExpression(identifier='BUFFER_SIZE'),
                                                                                                                                                                                                                                                                           operation='+',
                                                                                                                                                                                                                                                                           right=IdentifierExpression(identifier='PAGE_SIZE')),
                                                                                                                                                                                                                                            operation='-',
                                                                                                                                                                                                                                            right=IntExpression(value='1')),
                                                                                                                                                                                                             operation='/',
                                                                                                                                                                                                             right=IdentifierExpression(identifier='PAGE_SIZE'))))],
                                                                                                       isEndComma=False)),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='Token',
                                                                                                                              declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='Fragment',
                                                                                                                                                                                                                                                    declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                                                                                                                                        structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='Pos',
                                                                                                                                                                                                                                                                                                                                                                          declarationList=[Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Int: 'INT'>),
                                                                                                                                                                                                                                                                                                                                                                                                       varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='line')])]),
                                                                                                                                                                                                                                                                                                                                                                                           Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Int: 'INT'>),
                                                                                                                                                                                                                                                                                                                                                                                                       varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='col')])])])),
                                                                                                                                                                                                                                                                                 varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='starting')]),
                                                                                                                                                                                                                                                                                          AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='following')])])])),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='fragment')])]),
                                                                                                                                               Declaration(declarationBody=EnumTypeSpecifier(enumStatement=FullEnumStatement(identifier='',
                                                                                                                                                                                                                             enumeratorList=[Enumerator(identifier='Ident',
                                                                                                                                                                                                                                                        constantExpression=''),
                                                                                                                                                                                                                                             Enumerator(identifier='IntConst',
                                                                                                                                                                                                                                                        constantExpression=''),
                                                                                                                                                                                                                                             Enumerator(identifier='FloatConst',
                                                                                                                                                                                                                                                        constantExpression='')],
                                                                                                                                                                                                                             isEndComma=False)),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='type')])]),
                                                                                                                                               Declaration(declarationBody=StructOrUnionSpecifier(type='union',
                                                                                                                                                                                                  structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='',
                                                                                                                                                                                                                                                    declarationList=[Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Char: 'CHAR'>),
                                                                                                                                                                                                                                                                                 varName=[AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='name')]))]),
                                                                                                                                                                                                                                                                     Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Int: 'INT'>),
                                                                                                                                                                                                                                                                                 varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='int_value')])]),
                                                                                                                                                                                                                                                                     Declaration(declarationBody=SimpleTypeSpecifier(simpleType=<SimpleType.Double: 'DOUBLE'>),
                                                                                                                                                                                                                                                                                 varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='float_value')])])])),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='info')])])])),
                                     varName=[]),
                         Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                            structOrUnionSpecifier=FullStructOrUnionStatement(identifierOpt='List',
                                                                                                                              declarationList=[Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='Token')),
                                                                                                                                                           varName=[AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='value')])]),
                                                                                                                                               Declaration(declarationBody=StructOrUnionSpecifier(type='struct',
                                                                                                                                                                                                  structOrUnionSpecifier=EmptyStructOrUnionStatement(identifier='List')),
                                                                                                                                                           varName=[AbstractDeclaratorPointer(declarator=AbstractDeclaratorArrayList(arrays=[AbstractDeclaratorPrimSimple(identifier='next')]))])])),
                                     varName=[])])

Process finished with exit code 0

```

# Вывод
Лабораторная работа оказалась одной из самых сложных, что я делал.
Но за время выполнения я хорошо понял то,
как строить абстрактный и конкретный синтаксис и лучше углубился в структуру языка Си.
А так же совершил много ошибок в процессе написания кода: основная была в том,
что я сразу после того, как написал синтаксис, преобразовал его в код
(хотя и разобрался с тем, как работает библиотека), но код не запустился, выдавая ошибки.
Многие часы дебагинга не перевели к успеху и я принял другу тактику — записывать синтаксис в код постепенно,
проверяя работоспособность каждого правила. И это дало свои плоды, работа была сделана.
Так же постепенно я исправлял недочеты в синтаксисе, что было бы куда сложнее,
если бы я решил это делать в процессе дебага. Но после выполнения эта работа вызывает у меня чувство радости,
очень необычный и полезный опыт.

