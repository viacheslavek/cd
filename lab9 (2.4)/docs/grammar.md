Program -> DeclarationList .

DeclarationList -> .
DeclarationList -> Declaration DeclarationList .

Declaration -> TypeSpecifier AbstractDeclaratorsOpt ";" .


AbstractDeclaratorsOpt -> AbstractDeclarators .
AbstractDeclaratorsOpt -> .

AbstractDeclarators -> AbstractDeclarator InnerAbstractDeclarators .
AbstractDeclarators -> "," AbstractDeclarator InnerAbstractDeclarators .
InnerAbstractDeclarators -> .

AbstractDeclarator -> AbstractDeclaratorStar .
AbstractDeclarator -> AbstractDeclaratorArrayList .

AbstractDeclaratorStar -> "*" AbstractDeclarator .


AbstractDeclaratorArrayList -> AbstractDeclaratorArray InnerAbstractDeclaratorArrayList .
InnerAbstractDeclaratorArrayList -> AbstractDeclaratorArray AbstractDeclaratorArrayList.  
InnerAbstractDeclaratorArrayList -> .


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

EnumStatement -> IDENTIFIER FullEnumStatementOpt .
EnumStatement -> FullEnumStatement .

FullEnumStatementOpt -> FullEnumStatement .
FullEnumStatementOpt -> .

FullEnumStatement -> "{" EnumeratorList "}" .


EnumeratorList -> Enumerator InnerEnumeratorList .    
InnerEnumeratorList -> "," Enumerator .             
InnerEnumeratorList -> .


Enumerator -> IDENTIFIER EnumeratorExpressionOpt .



EnumeratorExpressionOpt -> "=" ConstantExpression .
EnumeratorExpressionOpt -> .

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

Factor -> sizeof "(" TypeSizeofSpecifier IDENTIFIER ")" .
Factor -> IDENTIFIER .
Factor -> INT .
Factor -> "(" Expression ")" .

TypeSizeofSpecifier -> STRUCT .
TypeSizeofSpecifier -> UNION .
TypeSizeofSpecifier -> ENUM .



StructOrUnionSpecifier -> StructOrUnion StructOrUnionStatement .

StructOrUnion -> STRUCT .
StructOrUnion -> UNION .

StructOrUnionStatement -> IDENTIFIER FullStructOrUnionStatementOpt .
StructOrUnionStatement -> FullStructOrUnionStatement .

FullStructOrUnionStatementOpt -> FullStructOrUnionStatement .
FullStructOrUnionStatementOpt -> .

FullStructOrUnionStatement -> "{" DeclarationList "}" .


## В РБНФ

```
Program ::= DeclarationList
DeclarationList ::= Declaration*
Declaration ::= TypeSpecifier AbstractDeclaratorsOpt ';' 


AbstractDeclaratorsOpt ::= AbstractDeclarators?

AbstractDeclarators ::= AbstractDeclarator (',' AbstractDeclarator)* 

AbstractDeclarator ::= AbstractDeclaratorPointer | AbstractDeclaratorArrayList 

AbstractDeclaratorPointer ::= '*' AbstractDeclarator 
AbstractDeclaratorArrayList ::= AbstractDeclaratorArray+

AbstractDeclaratorArray ::= AbstractDeclaratorPrimArray |
 AbstractDeclaratorPrimSimple | AbstractDeclaratorPrimDifficult

AbstractDeclaratorPrimArray ::= '[' Expression ']'
AbstractDeclaratorPrimSimple ::= IDENTIFIER
AbstractDeclaratorPrimDifficult ::= '(' AbstractDeclarator ')'


TypeSpecifier ::= SimpleTypeSpecifier | EnumTypeSpecifier | StructOrUnionSpecifier

SimpleTypeSpecifier ::= char | short | int | long | float | double 


EnumTypeSpecifier ::= ENUM EnumStatement 

EnumStatement ::= IdentEnumStatement | BodyEnumStatement
IdentEnumStatement ::= IDENTIFIER BodyEnumStatement?
BodyEnumStatement ::= '{' EnumeratorList '}'

EnumeratorList ::= Enumerator (',' Enumerator)*  
Enumerator ::= IDENTIFIER ('=' Expression)? 


StructOrUnionSpecifier ::= (struct | union) StructOrUnionStatement 

StructOrUnionStatement ::= IdentStructOrUnionStatement | BodyStructOrUnionStatement
IdentStructOrUnionStatement ::= IDENTIFIER BodyStructOrUnionStatement?
BodyStructOrUnionStatement ::= '{' DeclarationList '}' 


Expression ::= Term (('+' | '-') Term)*  

Term ::= Factor (('*' | '/') Factor)*

Factor ::= sizeof '(' (struct | union | enum) IDENTIFIER ')' | IDENTIFIER | INTEGER | '(' Expression ')'

```


