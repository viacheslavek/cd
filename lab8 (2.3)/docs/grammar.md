# 2.1 Индивидуальный вариант

```
/* аксиома помечена звёздочкой */
F ("n") ("(" E ")")
T (F T')
T' ("*" F T') ()
* E (T E')
E' ("+" T E') ()
```

# Реализация

`DECLARATIONS ::= REWRITING_RULE DECLARATIONS | epsilon `

`REWRITING_RULE ::= AxiomSign NonTerminal REWRITING | NonTerminal REWRITING`

`REWRITING ::= OpenBracket BODY CloseBracket REWRITING_OPT`

`REWRITING_OPT ::= OpenBracket BODY CloseBracket REWRITING_OPT | epsilon`

`BODY ::= Terminal BODY | NonTerminal BODY | epsilon`

```
AxiomSign ::= *
OpenBracket ::= (
CloseBracket ::= )
Whitespace ::= [ \t\n\r]+
Comment ::= /*([^(*/)]*)*/
NonTerminal ::= [a-zA-Z][a-zA-Z0-9]*(')?
Terminal    ::= "[^"]+"
```

