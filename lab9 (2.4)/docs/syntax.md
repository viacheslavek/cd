Приведем описание абстрактного и конкретного синтаксиса для задания

Определения структур, объединений и перечислений языка Си.
В инициализаторах перечислений допустимы знаки операций +, -, *, /, sizeof,
операндами могут служить имена перечислимых значений и целые числа.

Числовые константы могут быть только целочисленными и десятичными.

Примеры синтаксиса можно увидеть в папке tests


## Абстрактный синтаксис:

```
Program -> DeclarationList

DeclarationList -> Declaration*

Declaration -> TypeSpecifier AbstractDeclaratorsOpt ;


AbstractDeclaratorsOpt -> AbstractDeclarators | ε

AbstractDeclarators -> AbstractDeclarator*

AbstractDeclaratorOpt -> AbstractDeclarator | ε


AbstractDeclarator -> NPointerOpt IDENTIFIER


AbstractDeclarator -> NPointerOpt IDENTIFIER ListArraysOpt

NPointerOpt -> POINTER

ListArraysOpt -> [ ModifySimpleType ]*

ModifySimpleType -> SimpleType | IDENTIFIER

TypeSpecifier -> SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier


SimpleTypeSpecifier -> SimpleType

SimpleType -> CHAR | SHORT | INT | LONG | FLOAT | DOUBLE | SIGNED | UNSIGNED


EnumTypeSpecifier -> ENUM EnumStatement

EnumStatement -> FullEnumStatement | EmptyEnumStatement

FullEnumStatement -> IdentifierOpt { EnumeratorList CommaOpt }

IdentifierOpt -> IDENTIFIER | ε

EmptyEnumStatement -> IDENTIFIER

EnumeratorList -> EnumeratorList*

Enumerator -> IDENTIFIER EnumeratorExpressionOpt

EnumeratorExpressionOpt -> = ConstantExpression | ε


CommaOpt -> , | ε


ConstantExpression -> Expression

Expression -> IDENTIFIER | INT | Expression BinOp Expression | UnOp Expression
BinaryOperation -> + | - | ∙ | /
UnaryOperation -> + | - | sizeof ( TypeSpecifier AbstractDeclaratorsOpt )


StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement

StructOrUnion -> STRUCT | UNION

StructOrUnionStatement -> FullStructOrUnionStatement | EmptyStructOrUnionStatement

EmptyStructOrUnionStatement -> IDENTIFIER

FullStructOrUnionStatement -> IdentifierOpt { DeclarationList }


```

## Конкретный синтаксис:

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

Переводим в LL1

Program -> DeclarationList . 

DeclarationList -> .
DeclarationList -> Declaration DeclarationList .

Declaration -> TypeSpecifier AbstractDeclaratorsOpt ";" .


AbstractDeclaratorsOpt -> AbstractDeclarators .
AbstractDeclaratorsOpt -> .

AbstractDeclarators -> AbstractDeclarators "," AbstractDeclarator .
AbstractDeclarators -> AbstractDeclarator .

AbstractDeclarator -> AbstractDeclaratorStar .
AbstractDeclarator -> AbstractDeclaratorArrayListOpt .

AbstractDeclaratorStar -> "*" AbstractDeclarator .


AbstractDeclaratorArrayListOpt -> AbstractDeclaratorArrayList .
AbstractDeclaratorArrayListOpt -> .
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayListOpt AbstractDeclaratorArray .

AbstractDeclaratorArray -> AbstractDeclaratorArray "[" Expression "]" .
AbstractDeclaratorArray -> AbstractDeclaratorPrim .

AbstractDeclaratorPrim -> AbstractDeclaratorPrimSimple . 
AbstractDeclaratorPrim -> AbstractDeclaratorPrimDifficult .

AbstractDeclaratorPrimSimple -> IDENTIFIER .
AbstractDeclaratorPrimDifficult -> "(" AbstractDeclarator ")" .


TypeSpecifier -> SimpleTypeSpecifier .
TypeSpecifier -> EnumTypeSpecifier .
TypeSpecifier -> StructOrUnionSpecifier .

SimpleTypeSpecifier -> SimpleType .

SimpleType -> CHAR .
SimpleType -> SHORT .
SimpleType -> INT .
SimpleType -> LONG .
SimpleType -> FLOAT .
SimpleType -> DOUBLE .
SimpleType -> SIGNED .
SimpleType -> UNSIGNED .

EnumTypeSpecifier -> ENUM EnumStatement .

EnumStatement -> FullEnumStatement .
EnumStatement -> EmptyEnumStatement .

FullEnumStatement -> IdentifierOpt "{" EnumeratorList CommaOpt "}" .

IdentifierOpt -> IDENTIFIER .
IdentifierOpt -> .

EmptyEnumStatement -> IDENTIFIER .

EnumeratorList -> EnumeratorList "," Enumerator . 
EnumeratorList -> Enumerator .

Enumerator -> IDENTIFIER EnumeratorExpressionOpt .

EnumeratorExpressionOpt -> "=" ConstantExpression .
EnumeratorExpressionOpt -> .

CommaOpt -> "," .
CommaOpt -> .

ConstantExpression -> Expression .

Expression -> ArithmeticExpression .

ArithmeticExpression -> Term .
ArithmeticExpression -> "+" Term .
ArithmeticExpression -> "-" Term .
ArithmeticExpression -> ArithmeticExpression AddOperation Term .

AddOperation -> "+" .
AddOperation -> "-" .

Term -> Factor .
Term -> Term MultyOperation Factor .

MultyOperation -> "*" .
MultyOperation -> "/" .

Factor -> sizeof "(" TypeSpecifier AbstractDeclaratorsOpt ")" .
Factor -> IDENTIFIER .
Factor -> INT .
Factor -> "(" Expression ")" .

StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement .

StructOrUnion -> STRUCT . 
StructOrUnion -> UNION .

StructOrUnionStatement -> FullStructOrUnionStatement .
StructOrUnionStatement -> EmptyStructOrUnionStatement .

EmptyStructOrUnionStatement -> IDENTIFIER .

FullStructOrUnionStatement -> IdentifierOpt "{" DeclarationList "}" .




Пытаюсь сделать LL(1) ______________________________

// Вот это не работает - спросить 

EnumeratorList -> EnumeratorList "," Enumerator .

EnumeratorList -> Enumerator .

В

EnumeratorList -> Enumerator InnerEnumerators .
InnerEnumerators -> "," EnumeratorList .
InnerEnumerators -> .

Statements -> Statement InnerStatements .
InnerStatements -> "," Statement InnerStatements .
InnerStatements -> .


// ______________

AbstractDeclarators -> AbstractDeclarators "," AbstractDeclarator .
AbstractDeclarators -> AbstractDeclarator .

В 

AbstractDeclarators -> AbstractDeclarator InnerAbstractDeclarators .
AbstractDeclarators -> "," AbstractDeclarator AbstractDeclarators .
InnerAbstractDeclarators -> .


// ______________

____________________


ArithmeticExpression -> Term .
ArithmeticExpression -> "+" Term .
ArithmeticExpression -> "-" Term .
ArithmeticExpression -> ArithmeticExpression AddOperation Term .

AddOperation -> "+" .
AddOperation -> "-" .

Term -> Factor .
Term -> Term MultyOperation Factor .

MultyOperation -> "*" .
MultyOperation -> "/" .

В

ArithmeticExpression -> Term InnerArithmeticExpression .
InnerArithmeticExpression -> AddOperation Term InnerArithmeticExpression .
InnerArithmeticExpression -> .

AddOperation -> "+" .
AddOperation -> "-" .

Term -> Factor InnerTerm .
InnerTerm -> MultyOperation Factor InnerTerm .
InnerTerm -> .

MultyOperation -> "*" .
MultyOperation -> "/" .


/// ______________

DeclarationList -> Declaration .
DeclarationList -> DeclarationList Declaration .

в



_______

AbstractDeclaratorArray -> AbstractDeclaratorArray "[" Expression "]" .

В

AbstractDeclaratorArray -> "[" Expression "]" AbstractDeclaratorArray  .


/// ____________


AbstractDeclarator -> AbstractDeclaratorStar .
AbstractDeclarator ->  Opt .

AbstractDeclaratorStar -> "*" AbstractDeclarator .


AbstractDeclaratorArrayListOpt -> AbstractDeclaratorArrayList .
AbstractDeclaratorArrayListOpt -> .
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayListOpt AbstractDeclaratorArray .

AbstractDeclaratorArray -> "[" Expression "]" AbstractDeclaratorArray  .
AbstractDeclaratorArray -> AbstractDeclaratorPrim .


в

AbstractDeclarator -> AbstractDeclaratorStar .
AbstractDeclarator -> AbstractDeclaratorArrayList .

AbstractDeclaratorStar -> "*" AbstractDeclarator .

AbstractDeclaratorArrayList -> AbstractDeclaratorArray  .
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayList AbstractDeclaratorArray .

AbstractDeclaratorArray -> "[" Expression "]" .

AbstractDeclaratorArray -> AbstractDeclaratorPrim .
AbstractDeclaratorPrim -> AbstractDeclaratorPrimSimple .

// _______________

AbstractDeclaratorArrayList -> AbstractDeclaratorArray .
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayList AbstractDeclaratorArray .

в

AbstractDeclaratorArrayList -> AbstractDeclaratorArray InnerAbstractDeclaratorArrayList .
InnerAbstractDeclaratorArrayList -> AbstractDeclaratorArray InnerAbstractDeclaratorArrayList .
InnerAbstractDeclaratorArrayList -> .


// _______________


////// ____________________: тогда


Program -> DeclarationList .

DeclarationList -> .
DeclarationList -> Declaration DeclarationList .

Declaration -> TypeSpecifier AbstractDeclaratorsOpt ";" .


AbstractDeclaratorsOpt -> AbstractDeclarators .
AbstractDeclaratorsOpt -> .

AbstractDeclarators -> AbstractDeclarator InnerAbstractDeclarators .
AbstractDeclarators -> "," AbstractDeclarator AbstractDeclarators .
InnerAbstractDeclarators -> .

AbstractDeclarator -> AbstractDeclaratorStar .
AbstractDeclarator -> AbstractDeclaratorArrayList .

AbstractDeclaratorStar -> "*" AbstractDeclarator .

AbstractDeclaratorArrayList -> AbstractDeclaratorArray .
AbstractDeclaratorArrayList -> AbstractDeclaratorArrayList AbstractDeclaratorArray .

AbstractDeclaratorArray -> "[" Expression "]" .

AbstractDeclaratorArray -> AbstractDeclaratorPrim .

AbstractDeclaratorPrim -> AbstractDeclaratorPrimSimple .
AbstractDeclaratorPrim -> AbstractDeclaratorPrimDifficult .

AbstractDeclaratorPrimSimple -> IDENTIFIER .
AbstractDeclaratorPrimDifficult -> "(" AbstractDeclarator ")" .


TypeSpecifier -> SimpleTypeSpecifier .
TypeSpecifier -> EnumTypeSpecifier .
TypeSpecifier -> StructOrUnionSpecifier .

SimpleTypeSpecifier -> SimpleType .

SimpleType -> CHAR .
SimpleType -> SHORT .
SimpleType -> INT .
SimpleType -> LONG .
SimpleType -> FLOAT .
SimpleType -> DOUBLE .
SimpleType -> SIGNED .
SimpleType -> UNSIGNED .

EnumTypeSpecifier -> ENUM EnumStatement .

EnumStatement -> FullEnumStatement .
EnumStatement -> EmptyEnumStatement .

FullEnumStatement -> IdentifierOpt "{" EnumeratorList CommaOpt "}" .

IdentifierOpt -> IDENTIFIER .
IdentifierOpt -> .

EmptyEnumStatement -> IDENTIFIER .

EnumeratorList -> EnumeratorList "," Enumerator .
EnumeratorList -> Enumerator .

Enumerator -> IDENTIFIER EnumeratorExpressionOpt .

EnumeratorExpressionOpt -> "=" ConstantExpression .
EnumeratorExpressionOpt -> .

CommaOpt -> "," .
CommaOpt -> .

ConstantExpression -> Expression .

Expression -> ArithmeticExpression .

ArithmeticExpression -> Term InnerArithmeticExpression .
InnerArithmeticExpression -> AddOperation Term InnerArithmeticExpression .
InnerArithmeticExpression -> .

AddOperation -> "+" .
AddOperation -> "-" .

Term -> Factor InnerTerm .
InnerTerm -> MultyOperation Factor InnerTerm .
InnerTerm -> .

MultyOperation -> "*" .
MultyOperation -> "/" .

Factor -> sizeof "(" TypeSpecifier AbstractDeclaratorsOpt ")" .
Factor -> IDENTIFIER .
Factor -> INT .
Factor -> "(" Expression ")" .

StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement .

StructOrUnion -> STRUCT .
StructOrUnion -> UNION .

StructOrUnionStatement -> FullStructOrUnionStatement .
StructOrUnionStatement -> EmptyStructOrUnionStatement .

EmptyStructOrUnionStatement -> IDENTIFIER .

FullStructOrUnionStatement -> IdentifierOpt "{" DeclarationList "}" .


/// _______________________________________



