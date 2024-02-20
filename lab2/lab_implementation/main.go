package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
)

// TODO: вернуть к вводу как параметр
const filename = "/Users/slavaruswarrior/Documents/GitHub/cd/lab2/source_example/example.go"

/*
Заменить вхождения глобальной константы LINE на номер текущей строки, FILE — на имя текущего файла.
В исходной программе эти константы должны быть объявлены в глобальной области видимости с типами int и string,
соответственно (в противном случае следует выдать ошибку), значения этих констант могут быть произвольными.
В новой программе они должны отсутствовать.
*/

func findConstFileAndLineAndModify(file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		if genDecl, okDecl := node.(*ast.GenDecl); okDecl && checkConst(genDecl) {
			log.Println("find const")
			modifyFileAndLine(genDecl)
		}
		return true
	})
}

func modifyFileAndLine(gd *ast.GenDecl) {
	ast.Inspect(gd, func(node ast.Node) bool {
		if valueSpec, okValueSpec := node.(*ast.ValueSpec); okValueSpec {
			log.Println("find value")
			if isFileConst(valueSpec) {
				log.Println("find file")
				errFile := modifyFile(valueSpec)
				if errFile != nil {
					log.Fatalf("failed to parse ast: err in modify file: %e", errFile)
				}
			}
			if isLineConst(valueSpec) {
				log.Println("find line")
				errLine := modifyLine(valueSpec)
				if errLine != nil {
					log.Fatalf("failed to parse ast: err in modify file: %e", errLine)
				}
			}
		}
		return true
	})
}

func checkConst(gd *ast.GenDecl) bool {
	return gd.Tok.String() == "const"
}

func isFileConst(vs *ast.ValueSpec) bool {
	for _, name := range vs.Names {
		if name.String() == "FILE" {
			return true
		}
	}
	return false
}

func modifyFile(vs *ast.ValueSpec) error {
	return modify(vs, modifyValueFile, "FILE")
}

func modify(vs *ast.ValueSpec, fModifyValue func(value *ast.BasicLit) error, name string) error {
	filePosition, errPosition := findPosition(name, vs.Names, vs.Values)
	if errPosition != nil {
		return errPosition
	}

	value, ok := (vs.Values[filePosition]).(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("interface type does not match. needed ast.BasicLit")
	}

	if errValue := fModifyValue(value); errValue != nil {
		return errValue
	}

	vs.Values[filePosition] = value
	return nil
}

func findPosition(name string, Names []*ast.Ident, spec []ast.Expr) (int, error) {
	for i, n := range Names {
		if n.String() == name {
			if i > len(spec) {
				return -1, fmt.Errorf("can't len names %d bigger len values %d", i, len(spec))
			}
			return i, nil
		}
	}

	return -1, fmt.Errorf("failed to find %s", name)
}

// TODO: обрабатываю FILE value
func modifyValueFile(value *ast.BasicLit) error {
	log.Println("start modify value file")
	fmt.Println("value", value)
	return nil
}

func isLineConst(vs *ast.ValueSpec) bool {
	for _, name := range vs.Names {
		if name.String() == "LINE" {
			return true
		}
	}
	return false
}

func modifyLine(vs *ast.ValueSpec) error {
	return modify(vs, modifyValueLine, "LINE")
}

// TODO: обрабатываю FILE value
func modifyValueLine(value *ast.BasicLit) error {
	log.Println("start modify value line")
	fmt.Println("value", value)

	return nil
}

func main() {

	//if len(os.Args) != 2 {
	//	fmt.Printf("usage: demo <filename.go>\n")
	//	return
	//}

	astFile, err := os.Create("output/ast_output.txt")
	if err != nil {
		log.Fatalf("failed create file with err %e", err)
	}

	defer func() {
		_ = astFile.Close()
	}()

	newGoFile, err := os.Create("output/new_program.go")
	if err != nil {
		log.Fatalf("failed create file with err %e", err)
	}

	defer func() {
		_ = newGoFile.Close()
	}()

	fset := token.NewFileSet()
	if file, errPF := parser.ParseFile(fset, filename, nil, parser.ParseComments); errPF == nil {

		findConstFileAndLineAndModify(file)

		if errN := format.Node(io.Writer(newGoFile), fset, file); errN != nil {
			fmt.Printf("Formatter error: %v\n", errN)
		}
		errF := ast.Fprint(io.Writer(astFile), fset, file, nil)
		if errF != nil {
			log.Fatalf("failed write go file with err %e", errF)
		}
	} else {
		fmt.Printf("Errors in %s\n as %e", filename, errPF)
	}

	fmt.Println("finish")
}
