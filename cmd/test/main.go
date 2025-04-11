package main

import (
	"fmt"

	"github.com/ML1883/GoFigure/pkg/analyzer"
	"github.com/ML1883/GoFigure/pkg/parser"
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
		"The old oak tree in the town square has witnessed 42 generations of children playing beneath its sprawling branches.",
		"Morning light filtered through 7 dusty blinds, casting striped shadows across the worn 1950s hardwood floor.",
		"She carefully added the final 8th brushstroke to her painting, stepping back to admire the 365-day project.",
		"The ancient book's leather binding from 1783 cracked slightly as he opened it, releasing centuries' scent.",
		"Waves crashed rhythmically against the shoreline while 24 seagulls circled overhead in the 37-degree air.",
		"The antique watch from 1862 ticked steadily on his wrist, a family heirloom passed down through 4 generations.",
		"Fresh snow blanketed the 15 houses in the neighborhood, transforming the familiar streets into a 9-degree wonderland.",
		"Steam rose from her coffee cup as she sat by window 302, watching the city slowly come to life at 6:30 am.",
		"Thunder rumbled in the distance, prompting 53 farmers to hurry with their harvest before the 90 percent chance storm.",
		"The melody from the musician's $125 violin echoed through station 42, captivating 57 passing commuters.",
		"Fireflies, about 150 of them, blinked like tiny stars as 10 children chased them with glass jars on summer evenings.",
		"His footprints in the sand were quickly erased by the incoming tide at 3:45, leaving no trace of his passage.",
		"The scent of 36 freshly baked loaves wafted from the bakery, drawing a line of 27 customers down the block.",
		"Raindrops tapped a 4/4 rhythm on the tin roof, creating the perfect soundtrack for 2 hours of afternoon reading.",
		"Autumn leaves crunched beneath their feet as they walked 5.2 miles through the park on October 23rd.",
		"The old piano in room 19 had been silent for 12 years until she sat down and played a 200-year-old sonata.",
		"Moonlight silvered the 400-acre lake surface, reflecting the silhouettes of 97 pine trees along the shore.",
		"The vintage 1976 camera captured moments that digital technology from 2025 could never quite replicate.",
		"Stars appeared one by one in the darkening sky as 13 campers gathered around the 650-degree campfire.",
		"Morning dew glistened on 18 spider webs between fence posts, catching the first rays of 7:15 am sunlight.",
	}

	// Create a distribution model from the training texts
	model, err := analyzer.CreateDistributionFittedModel(trainingTexts, 2.5, 0.8) // 2.5 is the anomaly threshold
	if err != nil {
		fmt.Printf("Error creating model: %v\n", err)
		return
	}

	// Print model summary
	fmt.Println(model.GetModelSummary())

	// Save the model
	// modelPath := filepath.Join(".", "text_model.gob")
	// if err := model.SaveTextModel(modelPath); err != nil {
	// 	fmt.Printf("Error saving model: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Model saved to %s\n", modelPath)

	// // Load the model (just to demonstrate)
	// loadedModel, err := analyzer.LoadTextModel(modelPath)
	// if err != nil {
	// 	fmt.Printf("Error loading model: %v\n", err)
	// 	return
	// }
	// fmt.Println("Model loaded successfully")

	// Test text similar to training data
	normalText := "This is similar to our training texts with numbers 123."
	isAnomaly, score, anomalies, probability := model.IsAnomaly(normalText)
	fmt.Printf("\nNormal Text Analysis:\n")
	fmt.Printf("Is anomaly: %v\n", isAnomaly)
	fmt.Printf("Anomaly score: %.2f\n", score)
	fmt.Printf("Probability of observing this text: %.2f\n", probability)

	if len(anomalies) > 0 {
		fmt.Println("Top anomalies:")
		topAnomalies := model.GetTopAnomalies(normalText, 3)
		for _, a := range topAnomalies {
			fmt.Printf("  %s\n", a)
		}
	}

	// Test text with different distribution
	anomalousText := "ZZZZZZZZZZ999999999XXXXXXXXXX000000000000000000"
	isAnomaly, score, anomalies, probability = model.IsAnomaly(anomalousText)
	fmt.Printf("\nAnomalous Text Analysis:\n")
	fmt.Printf("Is anomaly: %v\n", isAnomaly)
	fmt.Printf("Anomaly score: %.2f\n", score)
	fmt.Printf("Probability of observing this text: %.2f\n", probability)

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
