package lab_lexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Scanner struct {
	tokens   []Token
	compiler *Compiler
	position int
}

func NewScanner(filepath string) *Scanner {
	return &Scanner{
		tokens:   ParseFile(filepath),
		compiler: NewCompiler(),
		position: 0,
	}
}

func (s *Scanner) NextToken() Token {
	// TODO: тут встраиваю все для компилятора
	if s.position < len(s.tokens) {
		token := s.tokens[s.position]
		s.position++
		return token
	}
	// TODO: возвращаю специальный токен EOF
	return Token{}

}

func (s *Scanner) GetCompiler() *Compiler {
	return s.compiler
}

func ParseFile(filepath string) []Token {
	file, errO := os.Open(filepath)
	if errO != nil {
		log.Fatalf("can't open file in parser %e", errO)
	}
	defer func() {
		_ = file.Close()
	}()

	tokens := make([]Token, 0)

	scanner := bufio.NewScanner(file)
	runeScanner := NewRunePosition(scanner)

	for runeScanner.GetRune() != -1 {
		for runeScanner.IsWhiteSpace() {
			runeScanner.NextRune()
		}
		start := runeScanner.GetCurrentPosition()

		switch runeScanner.GetRune() {
		case '"':
			tokens = append(tokens, processString())
		default:
			if runeScanner.IsLetter() {
				tokens = append(tokens, processIdentifier())
			} else {
				tokens = append(tokens, processStartError())
			}
		}

		// TODO: это потереть
		fmt.Printf("s: %d, %s | ", start, string(runeScanner.GetRune()))
		runeScanner.NextRune()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed read file by line in parser %e", err)
	}

	fmt.Println()

	return tokens
}

// TODO: делаю обработку строки
func processString() Token {
	// TODO: функция разбора строки или ошибки
	// Если message.isError true -> что-то пошло не так, мы получили ошибку, пропустили ее и добавляем ее
	// иначе мы нашли токен
	// в любом случае идем дальше
	return Token{}
}

// TODO: делаю обработку индефикатора
func processIdentifier() Token {
	return Token{}
}

// TODO: делаю обработку стартовой ошибки
func processStartError() Token {
	fmt.Println("Разбираем ошибку, что у меня нет таких токенов, с которых могло бы начинаться")
	return Token{}
}
