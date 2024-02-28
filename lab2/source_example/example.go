package main

import "fmt"

const FILE = "replace"

const LINE = 1

func main() {
	fmt.Println("FILE", FILE, "LINE", LINE)

	x := LINE
	y := FILE
	fmt.Println(x, y)
}

//const (
//	FILE = "replace"
//	LINE = 1
//)

//const (
//	FILE, LINE = "replace", 1
//)
