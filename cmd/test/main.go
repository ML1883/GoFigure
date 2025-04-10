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
	var returnStruct = analyzer.AnalyzeLettersFromText(parsedText)
	var returnStruct2 = analyzer.AnalyzeLettersFromText(parsedText2)
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

	trainingTexts := []string{
		"This is an example of normal text that follows certain patterns. It has numbers like 123 and 456.",
		"Another normal sample with similar distribution of characters and numbers 789.",
		"The quick brown fox jumps over the lazy dog. It contains all letters of the alphabet.",
		"We expect these texts to establish a baseline for what's considered normal in our model.",
		"One more sample text that will help define our distribution parameters.",
	}

	// Create a distribution model from the training texts
	model, err := analyzer.CreateDistributionModelWithFitting(trainingTexts, 2.5) // 2.5 is the anomaly threshold
	if err != nil {
		fmt.Printf("Error creating model: %v\n", err)
		return
	}

	// Print model summary
	fmt.Println(model.ModelSummary())
	fmt.Println(model.GetDistributionSummary())

	// // Save the model
	// modelPath := filepath.Join(".", "text_model.gob")
	// if err := model.SaveModel(modelPath); err != nil {
	// 	fmt.Printf("Error saving model: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Model saved to %s\n", modelPath)

	// // Load the model (just to demonstrate)
	// loadedModel, err := analyzer.LoadModel(modelPath)
	// if err != nil {
	// 	fmt.Printf("Error loading model: %v\n", err)
	// 	return
	// }
	// fmt.Println("Model loaded successfully")

	// Test text similar to training data
	normalText := "This is similar to our training texts with numbers 123."
	isAnomaly, score, anomalies := model.IsAnomaly(normalText)
	fmt.Printf("\nNormal Text Analysis:\n")
	fmt.Printf("Is anomaly: %v\n", isAnomaly)
	fmt.Printf("Anomaly score: %.2f\n", score)

	if len(anomalies) > 0 {
		fmt.Println("Top anomalies:")
		topAnomalies := model.GetTopAnomalies(normalText, 3)
		for _, a := range topAnomalies {
			fmt.Printf("  %s\n", a)
		}
	}

	// Test text with different distribution
	anomalousText := "ZZZZZZZZZZ999999999XXXXXXXXXX000000000000000000"
	isAnomaly, score, anomalies = model.IsAnomaly(anomalousText)
	fmt.Printf("\nAnomalous Text Analysis:\n")
	fmt.Printf("Is anomaly: %v\n", isAnomaly)
	fmt.Printf("Anomaly score: %.2f\n", score)

	if len(anomalies) > 0 {
		fmt.Println("Top anomalies:")
		topAnomalies := model.GetTopAnomalies(anomalousText, 5)
		for _, a := range topAnomalies {
			fmt.Printf("  %s\n", a)
		}
	}

	// Clean up
	// os.Remove(modelPath)

}
