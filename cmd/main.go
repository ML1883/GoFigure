package main

import (
	"GoFigure/pkg/analyzer"
	"GoFigure/pkg/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting main script.")
	var parsedText = parser.ParseStringToAlphanumeric("This is a test stringz0123")
	var parsedText2 = parser.ParseStringToAlphanumeric("How this does is the a cosine test react stringz0123 to changing order")
	var returnStruct = analyzer.CountLetters(parsedText)
	var returnStruct2 = analyzer.CountLetters(parsedText2)
	totalValue, err := analyzer.IntVectorMultiplication(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Returned struct: %v\n Return totalcount: %v\n Return position array: %v\n", returnStruct.LetterNumberArray, returnStruct.TotalCount, returnStruct.PositionArray)
		fmt.Printf("Returned struct2: %v\n Return totalcount2: %v\n Return position array2: %v\n", returnStruct2.LetterNumberArray, returnStruct2.TotalCount, returnStruct2.PositionArray)
		fmt.Printf("Total value %v\n", totalValue)
	}
	cosineSimilarity := analyzer.CosineSimilarityVectors(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	fmt.Printf("Cosine similairity of the two arrays: %v\n", cosineSimilarity)

}
