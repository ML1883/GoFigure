package main

import (
	"flag"
	"fmt"

	"github.com/ML1883/GoFigure/pkg/analyzer"
	"github.com/ML1883/GoFigure/pkg/parser"
)

/*
Simple file to setup a CLI interface for putting in two texts.
*/
func main() {
	fileMode := flag.Bool("file", false, "Use file input mode instead of direct text input")
	file1 := flag.String("text1", "", "Path to first text file (when using -file)")
	file2 := flag.String("text2", "", "Path to second text file (when using -file)")
	output := flag.Bool("output", false, "Output detailed vectors and arrays")
	help := flag.Bool("help", false, "Show help information")

	flag.Parse()

	if *help {
		fmt.Println("============================")
		fmt.Println("Text Similarity Analysis Tool")
		fmt.Println("============================")
		fmt.Println("This tool analyzes the similarity between two text inputs.")
		fmt.Println("\nUsage:")
		fmt.Println("  Interactive mode (default): Just run the program without flags")
		fmt.Println("  File mode: Use -file flag with -text1 and -text2 to specify input files")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		return
	}

	var text1, text2 string

	if *fileMode {
		if *file1 == "" || *file2 == "" {
			fmt.Println("Error: In file mode, you must specify both -text1 and -text2 file paths")
			flag.PrintDefaults()
			return
		}

		var err error
		text1, err = parser.ReadFile(*file1)
		if err != nil {
			fmt.Printf("Error reading first file: %v\n", err)
			return
		}

		text2, err = parser.ReadFile(*file2)
		if err != nil {
			fmt.Printf("Error reading second file: %v\n", err)
			return
		}
	} else {
		fmt.Println("=======================")
		fmt.Println("Text Similarity Analysis")
		fmt.Println("=======================")
		fmt.Println("Enter first text (type 'END' on a new line when finished):")
		text1 = parser.ReadMultilineInput()

		fmt.Println("\nEnter second text (type 'END' on a new line when finished):")
		text2 = parser.ReadMultilineInput()
	}

	parsedText1 := parser.ParseStringToAlphanumeric(text1)
	parsedText2 := parser.ParseStringToAlphanumeric(text2)

	returnStruct1 := analyzer.AnalyzeLettersFromText(parsedText1)
	returnStruct2 := analyzer.AnalyzeLettersFromText(parsedText2)

	cosineSimilarity := analyzer.CosineSimilarityVectors(returnStruct1.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	jaccardIndex := analyzer.JaccardIndexVectors(returnStruct1.LetterNumberArray[:], returnStruct2.LetterNumberArray[:])
	positionalCalc := analyzer.PositionDifferenceVectors(returnStruct1.PositionArray[:], returnStruct2.PositionArray[:], returnStruct1.TotalCount, returnStruct2.TotalCount)

	if *output {
		fmt.Printf("\nText 1 Letter Counts: %v\n", returnStruct1.LetterNumberArray)
		fmt.Printf("Text 1 Total Count: %v\n", returnStruct1.TotalCount)
		fmt.Printf("Text 1 Position Array: %v\n", returnStruct1.PositionArray)

		fmt.Printf("\nText 2 Letter Counts: %v\n", returnStruct2.LetterNumberArray)
		fmt.Printf("Text 2 Total Count: %v\n", returnStruct2.TotalCount)
		fmt.Printf("Text 2 Position Array: %v\n", returnStruct2.PositionArray)
	}

	fmt.Println("==================")
	fmt.Println("\nSimilarity Results:")
	fmt.Println("==================")
	fmt.Printf("Cosine Similarity: %v\n", cosineSimilarity)
	fmt.Printf("Jaccard Index: %v\n", jaccardIndex)
	fmt.Printf("Position Index: %v\n", positionalCalc)

	// Calculate combined similarity measures
	combinedSim := (cosineSimilarity + jaccardIndex + (1.0 - positionalCalc)) / 3
	weightedSim := (0.4 * cosineSimilarity) + (0.3 * jaccardIndex) + (0.3 * (1.0 - positionalCalc))

	fmt.Printf("\nEqually weighted Similarity (average): %v\n", combinedSim)
	fmt.Printf("Weighted Similarity (40-30-30 Cosine-Jaccard-Position): %v\n", weightedSim)
}
