package main

import "fmt"

const FILE = "/Users/slavaruswarrior/Documents/GitHub/cd/lab2/source_example/example.go"

const LINE = 7

func main() {
	fmt.Println("FILE", "/Users/slavaruswarrior/Documents/GitHub/cd/lab2/source_example/example.go", "LINE", 10)

	x := 12
	y := "/Users/slavaruswarrior/Documents/GitHub/cd/lab2/source_example/example.go"
	fmt.Println(x, y)
}

//const (
//	FILE = "replace"
//	LINE = 1
//)

//const (
//	FILE, LINE = "replace", 1
//)
