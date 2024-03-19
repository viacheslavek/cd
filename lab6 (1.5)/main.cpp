%option noyywrap bison-bridge bison-locations
%{

#include <stdio.h>
#include <stdlib.h>


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

void err (char *msg)
{
    // TODO: кладу ошибки в список ошибок
    printf("Error");
    print_pos(&cur);
    printf(":%s\n",msg);
}

// TODO: здесь описываю свои функции работы с таблицами идентификаторов
// мне нужны create_ident_table, add_ident_table, print_ident_table

%}

OPERATION_OPEN_BRACKET     [(]
OPERATION_CLOSE_BRACKET    [)]
OPERATION_LOWER            [<]
OPERATION_BIGGER           [>]

%%

[\n\t ]+

{OPERATION_OPEN_BRACKET}    {
                                yylval->operation = yytext;
                                return TAG_OPERATION_OPEN_BRACKET;
                            }

{OPERATION_CLOSE_BRACKET}   {
                                yylval->operation = yytext;
                                return TAG_OPERATION_CLOSE_BRACKET;
                            }

{OPERATION_LOWER}           {
                                yylval->operation = yytext;
                                return TAG_OPERATION_LOWER;
                            }

{OPERATION_BIGGER}          {
                                yylval->operation = yytext;
                                return TAG_OPERATION_BIGGER;
                            }

.                           err("ERROR unknown symbol");

<<EOF>>                     return 0;


%%

// TODO: добавить - два состояния - доллар, потом просто заглавные буквы, но не пустые!
// Директивы: любой знак валюты ($), после которого следует непустая последовательность заглавных букв.

// TODO: добавить - два состояния - заглавная буква, потом буквы, цифры и дефис
// Идентификаторы: последовательности буквенных символов ASCII, цифр и дефисов, начинающиеся с заглавной буквы.

int main()
{
    int tag;
    YYSTYPE value;
    YYLTYPE coords;

   	FILE *input;
	long size;
	char *buf;

    {
        input = fopen("test_files/operation_error.txt","r");
        fseek(input, 0, SEEK_END);
        size = ftell(input);
        rewind(input);
        buf = (char*)malloc(sizeof(char) * (size + 1));
        size_t n = fread(buf, sizeof(char), size, input);
    }

    buf[size] = '\0';
    fclose(input);

    init_scanner(buf);

    printf("START\n");

    printf("buf: %s\n", buf);

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

    free(buf);

    printf("FINISH\n");

    return 0;
}
