package main

import (
	"GoFigure/pkg/analyzer"
	"GoFigure/pkg/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting main script.")
	var parsedText = parser.ParseStringToAlphanumeric("This is a test stringz0123")
	var parsedText2 = parser.ParseStringToAlphanumeric("This is a test stringz0123")
	var returnStruct = analyzer.CountLetters(parsedText)
	var returnStruct2 = analyzer.CountLetters(parsedText2)
	totalValue, err := analyzer.IntVectorMultiplication(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Returned struct: %v\n Return totalcount: %v\n", returnStruct.LetterNumberArray, returnStruct.TotalCount)
		fmt.Printf("Returned struct2: %v\n Return totalcount2: %v\n", returnStruct2.LetterNumberArray, returnStruct2.TotalCount)
		fmt.Printf("Total value %v\n", totalValue)
	}

}
