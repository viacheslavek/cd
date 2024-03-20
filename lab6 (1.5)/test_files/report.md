% Лабораторная работа № 1.5 «Порождение лексического анализатора с помощью flex»
% 20 марта 2024 г.
% Вячеслав Локшин, ИУ9-61Б

# Цель работы
Целью данной работы является изучение генератора лексических анализаторов flex.

# Индивидуальный вариант
Идентификаторы: последовательности буквенных символов ASCII, цифр и дефисов,
начинающиеся с заглавной буквы.
Директивы: любой знак валюты, после которого следует непустая последовательность заглавных букв.
Знаки операций: «(», «)», «<», «>».

# Реализация

```lex
%option noyywrap bison-bridge bison-locations
%{

#include <stdio.h>
#include <stdlib.h>

#include <vector>
#include <unordered_map>
#include <string>
#include <iostream>

#define TAG_IDENT                     1
#define TAG_DIRECT                    2
#define TAG_OPERATION_OPEN_BRACKET    3
#define TAG_OPERATION_CLOSE_BRACKET   4
#define TAG_OPERATION_LOWER           5
#define TAG_OPERATION_BIGGER          6
#define TAG_ERROR                     7


char *tag_names[] =
{
    "END_OF_PROGRAM", "IDENT", "DIRECT",
    "TAG_OPERATION_OPEN_BRACKET", "TAG_OPERATION_CLOSE_BRACKET",
    "TAG_OPERATION_LOWER", "TAG_OPERATION_BIGGER",
    "ERROR"
};

struct Position
{
    int line, pos, index;
};

typedef struct Position Position;

void print_pos(Position *p)
{
    printf("(%d,%d)",p->line,p->pos);
}

struct Fragment
{
    Position starting, following;
};

typedef struct Fragment Fragment;

typedef struct Fragment YYLTYPE;

void print_frag(Fragment* f)
{
    print_pos(&(f->starting));
    printf("-");
    print_pos(&(f->following));
}

union Token
{
    char *direct;
    int ident_num;
    char *operation;
};

typedef union Token YYSTYPE;

int continued;
struct Position cur;

#define YY_USER_ACTION               \
{                                    \
    int i;                           \
    if (!continued)                  \
        yylloc->starting = cur;      \
    continued = 0;                   \
                                     \
    for (i = 0; i < yyleng; i++)     \
    {                                \
        if (yytext[i] == '\n')       \
        {                            \
            cur.line++;              \
            cur.pos = 1;             \
        }                            \
        else                         \
            cur.pos++;               \
        cur.index++;                 \
    }                                \
                                     \
    yylloc->following = cur;         \
}

void init_scanner (char *program)
{
    continued = 0;
    cur.line = 1;
    cur.pos = 1;
    cur.index = 0;
    yy_scan_string(program);
}

struct Error {
    Position pos;
    char     *message;
};

std::vector<Error> errors;

void add_error(Position pos, char* message) {
    Error error;
    error.pos = pos;
    error.message = message;
    errors.push_back(error);
}

void print_errors() {
    for (auto& error : errors) {
        printf("Error at ");
        print_pos(&(error.pos));
        printf(": %s\n", error.message);
    }
}


class IdentifierTable {
private:
    std::vector<std::string> identifiers;
    std::unordered_map<std::string, int> indexMap;

public:
    int add_identifier(const std::string& identifier) {
        auto it = indexMap.find(identifier);
        if (it != indexMap.end()) {
            return it->second;
        } else {
            int index = identifiers.size();
            identifiers.push_back(identifier);
            indexMap[identifier] = index;
            return index;
        }
    }

    void print_identifiers() const {
        std::cout << "Identifier Table:\n";
        for (int i = 0; i < identifiers.size(); ++i) {
            std::cout << i << ": " << identifiers[i] << std::endl;
        }
    }
};

IdentifierTable table;


%}

OPERATION_OPEN_BRACKET     [(]
OPERATION_CLOSE_BRACKET    [)]
OPERATION_LOWER            [<]
OPERATION_BIGGER           [>]

CAPITAL_LETTER [A-Z]
LETTER         [a-z]

DIRECT_START   [$]
DIRECT         {DIRECT_START}{CAPITAL_LETTER}+

DIGIT          [0-9]
DASH           [-]
IDENT          {CAPITAL_LETTER}({CAPITAL_LETTER}|{LETTER}|{DIGIT}|{DASH})*


%%

{OPERATION_OPEN_BRACKET}   {
                               yylval->operation = yytext;
                               return TAG_OPERATION_OPEN_BRACKET;
                           }
{OPERATION_CLOSE_BRACKET}  {
                               yylval->operation = yytext;
                               return TAG_OPERATION_CLOSE_BRACKET;
                           }
{OPERATION_LOWER}          {
                               yylval->operation = yytext;
                               return TAG_OPERATION_LOWER;
                           }
{OPERATION_BIGGER}         {
                               yylval->operation = yytext;
                               return TAG_OPERATION_BIGGER;
                           }


{DIRECT}  {
              yylval->direct = yytext;
              return TAG_DIRECT;
          }


{IDENT}  {
             yylval->ident_num = table.add_identifier(yytext);
             return TAG_IDENT;
         }


[\n\t ]+

.            add_error(cur, "ERROR unknown symbol");

<<EOF>>      return 0;


%%

int main()
{
    int tag;
    YYSTYPE value;
    YYLTYPE coords;

   	FILE *input;
	long size;
	char *buf;

    {
        input = fopen("test_files/mixed.txt","r");
        fseek(input, 0, SEEK_END);
        size = ftell(input);
        rewind(input);
        buf = (char*)malloc(sizeof(char) * (size + 1));
        fread(buf, sizeof(char), size, input);
    }

    buf[size] = '\0';
    fclose(input);

    init_scanner(buf);

    printf("buf: %s\n", buf);

    printf("START\n");

    do
    {
        tag = yylex(&value,&coords);
        if (tag != 0)
        {
            printf("%s ", tag_names[tag]);
            print_frag(&coords);
            printf(": ");

            if (tag == TAG_IDENT)
            {
                printf("%d", value.ident_num);
            }

            if (tag == TAG_DIRECT)
            {
                printf("%s", value.direct);
            }

            if (tag == TAG_OPERATION_OPEN_BRACKET || tag == TAG_OPERATION_CLOSE_BRACKET
                || tag == TAG_OPERATION_LOWER || tag == TAG_OPERATION_BIGGER)
            {
                printf("%s", value.operation);
            }

            printf("\n");
        }
    }
    while (tag != 0);
    printf("\n");

    table.print_identifiers();
    printf("\n");

    print_errors();
    printf("\n");

    free(buf);

    printf("FINISH\n");

    return 0;
}

```

# Тестирование

Для тестирования было написано несколько тестирующих файлов,
покрывающих большую часть тестовых случаев

Входные данные

Файл `direct_correct.txt`
```
$QWE $ER $AS
$QWERTY
```

Файл `direct_error.txt`
```
$WER $ _ $RUB
$wer $___ $US $
```

Файл `ident_correct.txt`
```
Hello Slava123
Vanya19-12
VOVA VLadiMIr
Anton- Ale----xey
Iv1An5e
```

Файл `ident_error.txt`
```
hello Slava123
Vanya19----12_
VOVA__ VLadiMIr
-Anton- Ale-xey
Iv1An5e
```

Файл `mixed.txt`
```
(())<()> < $ER $AS
$ _ $RUB Vanya19----12_
Anton- Ale----xey
```

Файл `operation_correct.txt`
```
( ) < <
> < ( ( >
(( << >> ))
((()))<(())>()
```

Файл `operation_error.txt`
```
( ) < -
_ ) [ ]
<) (>
(- -)
```

Вывод на `stdout` для *mixed.txt*

```
START
TAG_OPERATION_OPEN_BRACKET (1,1)-(1,2): (
TAG_OPERATION_OPEN_BRACKET (1,2)-(1,3): (
TAG_OPERATION_CLOSE_BRACKET (1,3)-(1,4): )
TAG_OPERATION_CLOSE_BRACKET (1,4)-(1,5): )
TAG_OPERATION_LOWER (1,5)-(1,6): <
TAG_OPERATION_OPEN_BRACKET (1,6)-(1,7): (
TAG_OPERATION_CLOSE_BRACKET (1,7)-(1,8): )
TAG_OPERATION_BIGGER (1,8)-(1,9): >
TAG_OPERATION_LOWER (1,10)-(1,11): <
DIRECT (1,12)-(1,15): $ER
DIRECT (1,16)-(1,19): $AS
DIRECT (2,5)-(2,9): $RUB
IDENT (2,10)-(2,23): 0
IDENT (3,1)-(3,7): 1
IDENT (3,8)-(3,18): 2

Identifier Table:
0: Vanya19----12
1: Anton-
2: Ale----xey

Error at (2,2): ERROR unknown symbol
Error at (2,4): ERROR unknown symbol
Error at (2,24): ERROR unknown symbol

FINISH

```

# Вывод
В ходе данной работы был изучен генератор лексических анализаторов flex и
написана программа, распознающая идентификаторы, директивы и операции по условию задачи.
Использование flex мне понравилось, получаются наглядные структуры 
и совсем немного кода, при этом написание лексического анализатора занимает 
не так много времени. 

