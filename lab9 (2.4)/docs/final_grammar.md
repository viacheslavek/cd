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

