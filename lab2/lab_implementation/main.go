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
	"strconv"
)

// TODO: вернуть к вводу как параметр
const exampleFilename = "/Users/slavaruswarrior/Documents/GitHub/cd/lab2/source_example/example.go"

/*
Заменить вхождения глобальной константы LINE на номер текущей строки, FILE — на имя текущего файла.
В исходной программе эти константы должны быть объявлены в глобальной области видимости с типами int и string,
соответственно (в противном случае следует выдать ошибку), значения этих констант могут быть произвольными.
В новой программе они должны отсутствовать.
*/

func findConstFileAndLineAndModify(fset *token.FileSet, file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		if genDecl, okDecl := node.(*ast.GenDecl); okDecl && checkConst(genDecl) {
			log.Println("find const")
			modifyFileAndLine(fset, genDecl)
		}
		return true
	})
}

func modifyFileAndLine(fset *token.FileSet, gd *ast.GenDecl) {
	ast.Inspect(gd, func(node ast.Node) bool {
		if valueSpec, okValueSpec := node.(*ast.ValueSpec); okValueSpec {
			log.Println("find value")
			if isFileConst(valueSpec) {
				log.Println("find file")
				errFile := modifyFile(fset, valueSpec)
				if errFile != nil {
					log.Fatalf("failed to parse ast: err in modify file: %e", errFile)
				}
			}
			if isLineConst(valueSpec) {
				log.Println("find line")
				errLine := modifyLine(fset, valueSpec)
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

func modifyFile(fset *token.FileSet, vs *ast.ValueSpec) error {
	return modify(fset, vs, modifyValueFile, "FILE")
}

func modify(fset *token.FileSet, vs *ast.ValueSpec,
	fModifyValue func(fset *token.FileSet, value *ast.BasicLit) error, name string) error {
	filePosition, errPosition := findPosition(name, vs.Names, vs.Values)
	if errPosition != nil {
		return errPosition
	}

	value, ok := (vs.Values[filePosition]).(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("interface type does not match. needed ast.BasicLit")
	}

	if errValue := fModifyValue(fset, value); errValue != nil {
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

func modifyValueFile(fset *token.FileSet, value *ast.BasicLit) error {
	log.Println("start modify value file")

	if value.Kind.String() != "STRING" {
		return fmt.Errorf("type does not match: need %s, given: %s", "STRING", value.Kind.String())
	}

	value.Value = getFilename(fset.Position(value.ValuePos).Filename)

	return nil
}

func getFilename(filename string) string {
	return fmt.Sprintf("\"%s\"", filename)
}

func isLineConst(vs *ast.ValueSpec) bool {
	for _, name := range vs.Names {
		if name.String() == "LINE" {
			return true
		}
	}
	return false
}

func modifyLine(fset *token.FileSet, vs *ast.ValueSpec) error {
	return modify(fset, vs, modifyValueLine, "LINE")
}

func modifyValueLine(fset *token.FileSet, value *ast.BasicLit) error {
	log.Println("start modify value line")

	if value.Kind.String() != "INT" {
		return fmt.Errorf("type does not match: need %s, given: %s", "INT", value.Kind.String())
	}

	value.Value = strconv.Itoa(fset.Position(value.ValuePos).Line)

	return nil
}

func findVarAndModify(fset *token.FileSet, file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		if assignStmt, okAssignStmt := node.(*ast.AssignStmt); okAssignStmt {
			log.Println("find *ast.AssignStmt")
			replaceResultInArr(fset, assignStmt)
		}
		return true
	})
}

func replaceResultInArr(fset *token.FileSet, as *ast.AssignStmt) {
	for i, rhsI := range as.Rhs {
		if rhs, okRhs := rhsI.(*ast.Ident); okRhs {
			if rhs.Name == "LINE" {
				rhs.Name = strconv.Itoa(fset.Position(rhs.Pos()).Line)
				as.Rhs[i] = rhs
			}
			if rhs.Name == "FILE" {
				rhs.Name = getFilename(fset.Position(rhs.Pos()).Filename)
				as.Rhs[i] = rhs
			}
		}
	}
}

func findVarInFunc(fset *token.FileSet, file *ast.File) {
	ast.Inspect(file, func(node ast.Node) bool {
		if callExpr, okCallExpr := node.(*ast.CallExpr); okCallExpr {
			log.Println("find *CallExpr")
			forArgsInFunc(fset, callExpr)
		}
		return true
	})
}

func forArgsInFunc(fset *token.FileSet, ce *ast.CallExpr) {
	for i, a := range ce.Args {
		if ident, okIdent := a.(*ast.Ident); okIdent {
			if ident.Name == "FILE" {
				ident.Name = getFilename(fset.Position(ident.Pos()).Filename)
				ce.Args[i] = ident
			}
			if ident.Name == "LINE" {
				ident.Name = strconv.Itoa(fset.Position(ident.Pos()).Line)
				ce.Args[i] = ident
			}
		}
	}
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
	if file, errPF := parser.ParseFile(fset, exampleFilename, nil, parser.ParseComments); errPF == nil {

		findConstFileAndLineAndModify(fset, file)

		findVarAndModify(fset, file)

		findVarInFunc(fset, file)

		if errN := format.Node(io.Writer(newGoFile), fset, file); errN != nil {
			fmt.Printf("Formatter error: %v\n", errN)
		}
		errF := ast.Fprint(io.Writer(astFile), fset, file, nil)
		if errF != nil {
			log.Fatalf("failed write go file with err %e", errF)
		}
	} else {
		fmt.Printf("Errors in %s\n as %e", exampleFilename, errPF)
	}

	fmt.Println("finish")
}
