package main

import (
	"GoFigure/pkg/analyzer"
	"GoFigure/pkg/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting main script.")
	var parsed_text = parser.ParseStringToAlphanumeric("This is a test string")
	var return_struct = analyzer.CountLetters(parsed_text)
	fmt.Printf("Returned struct: %v\n", return_struct)
}
