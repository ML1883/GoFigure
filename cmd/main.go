package main

import (
	"GoFigure/pkg/analyzer"
	"GoFigure/pkg/parser"
	"fmt"
)

func main() {
	fmt.Println("Starting main script.")
	var parsedText = parser.ParseStringToAlphanumeric("         aaaa")
	var parsedText2 = parser.ParseStringToAlphanumeric("aaaa         ")
	// var parsedText = parser.ParseStringToAlphanumeric("This is a test string123aaaa")
	// var parsedText2 = parser.ParseStringToAlphanumeric("aaaaThis is a test string123")
	var returnStruct = analyzer.CountLetters(parsedText)
	var returnStruct2 = analyzer.CountLetters(parsedText2)
	// totalValue, err := analyzer.IntVectorMultiplication(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// } else {
	fmt.Printf("Returned struct: %v\n Return totalcount: %v\n Return position array: %v\n", returnStruct.LetterNumberArray, returnStruct.TotalCount, returnStruct.PositionArray)
	fmt.Printf("Returned struct2: %v\n Return totalcount2: %v\n Return position array2: %v\n", returnStruct2.LetterNumberArray, returnStruct2.TotalCount, returnStruct2.PositionArray)
	// fmt.Printf("Total value %v\n", totalValue)
	// }
	cosineSimilarity := analyzer.CosineSimilarityVectors(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	JaccardIndex := analyzer.JaccardIndexVectors(returnStruct.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	fmt.Printf("Cosine similairity of the two arrays: %v\n", cosineSimilarity)
	fmt.Printf("Jaccard index of the two arrays: %v\n", JaccardIndex)
	positionalCalc := analyzer.PositionDifferenceVectors(returnStruct.PositionArray[:], returnStruct2.PositionArray[:], returnStruct.TotalCount, returnStruct2.TotalCount)
	fmt.Printf("Position index is: %v\n", positionalCalc)

}
