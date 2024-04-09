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

`DECLARATION ::= IS_AXIOM REWRITING_RULE DECLARATION | epsilon `

`IS_AXIOM ::= AxiomSign | epsilon`

`REWRITING_RULE ::= NonTerminal REWRITING`

`REWRITING ::= OpenBracket BODY CloseBracket REWRITING | epsilon`

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

