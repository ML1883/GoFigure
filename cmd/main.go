package main

import (
	"GoFigure/pkg/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting main script.")
	parser.ParseStringToAlphanumeric("This is a test string")
}
