package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ML1883/GoFigure/pkg/analyzer"
	"github.com/ML1883/GoFigure/pkg/parser"
)

func main() {
	// Common flags
	helpFlag := flag.Bool("help", false, "Show help information")
	outputFlag := flag.Bool("output", false, "Output detailed vectors and arrays")

	// Mode selection flags
	compareFlag := flag.Bool("compare", false, "Compare two texts for similarity")
	distributionFlag := flag.Bool("distribution", false, "Create or use a statistical distribution model")

	// Comparison mode flags
	fileModeFlag := flag.Bool("file", false, "Use file input mode instead of direct text input")
	file1Flag := flag.String("text1", "", "Path to first text file (when using -file)")
	file2Flag := flag.String("text2", "", "Path to second text file (when using -file)")

	// Distribution mode flags
	createModelFlag := flag.Bool("create-model", false, "Create a new distribution model")
	useModelFlag := flag.Bool("use-model", false, "Use an existing distribution model for analysis")
	folderFlag := flag.String("folder", "", "Path to folder containing training text files")
	modelFileFlag := flag.String("model-file", "text_model.gob", "Path to save/load model file")
	checkTextFlag := flag.String("check-text", "", "Path to text file to check against model")
	anomalyThresholdFlag := flag.Float64("threshold", 2.0, "Threshold for anomaly detection (higher = more strict)")
	fitThresholdFlag := flag.Float64("fit-threshold", 0.8, "Threshold for distribution fitting (higher = more empirical)")

	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	// If no mode is specified, default to comparison mode
	if !*compareFlag && !*distributionFlag {
		*compareFlag = true
	}

	if *compareFlag {
		runComparisonMode(*fileModeFlag, *file1Flag, *file2Flag, *outputFlag)
	}

	if *distributionFlag {
		if *createModelFlag {
			createDistributionModel(*folderFlag, *modelFileFlag, *anomalyThresholdFlag, *fitThresholdFlag, *outputFlag)
		} else if *useModelFlag {
			useDistributionModel(*modelFileFlag, *checkTextFlag, *outputFlag)
		} else {
			fmt.Println("Error: In distribution mode, you must specify either -create-model or -use-model")
			flag.PrintDefaults()
		}
	}
}

func showHelp() {
	fmt.Println("========================================")
	fmt.Println("Text Analysis Toolkit")
	fmt.Println("========================================")
	fmt.Println("This tool provides text analysis capabilities:")
	fmt.Println("1. Text similarity comparison")
	fmt.Println("2. Statistical distribution modeling")
	fmt.Println("3. Anomaly detection")
	fmt.Println("\nUsage Modes:")
	fmt.Println(" Comparison Mode (default): -compare")
	fmt.Println(" Distribution Mode: -distribution with -create-model or -use-model")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println(" Compare two files:")
	fmt.Println("   ./program -compare -file -text1=file1.txt -text2=file2.txt")
	fmt.Println(" Create distribution model from folder of texts:")
	fmt.Println("   ./program -distribution -create-model -folder=./training_texts -model-file=model.gob")
	fmt.Println(" Check text against model:")
	fmt.Println("   ./program -distribution -use-model -model-file=model.gob -check-text=sample.txt")
}

func runComparisonMode(fileMode bool, file1 string, file2 string, outputDetails bool) {
	var text1, text2 string

	if fileMode {
		if file1 == "" || file2 == "" {
			fmt.Println("Error: In file mode, you must specify both -text1 and -text2 file paths")
			flag.PrintDefaults()
			return
		}

		var err error
		text1, err = parser.ReadFile(file1)
		if err != nil {
			fmt.Printf("Error reading first file: %v\n", err)
			return
		}

		text2, err = parser.ReadFile(file2)
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

	if outputDetails {
		fmt.Printf("\nText 1 Letter Counts: %v\n", returnStruct1.LetterNumberArray)
		fmt.Printf("Text 1 Total Count: %v\n", returnStruct1.TotalCount)
		fmt.Printf("Text 1 Position Array: %v\n", returnStruct1.PositionArray)
		fmt.Printf("\nText 2 Letter Counts: %v\n", returnStruct2.LetterNumberArray)
		fmt.Printf("Text 2 Total Count: %v\n", returnStruct2.TotalCount)
		fmt.Printf("Text 2 Position Array: %v\n", returnStruct2.PositionArray)
	}

	fmt.Println("==================")
	fmt.Println("Similarity Results:")
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

func createDistributionModel(folderPath string, modelFilePath string, anomalyThreshold float64, fitThreshold float64, outputDetails bool) {
	if folderPath == "" {
		fmt.Println("Error: You must specify a folder path (-folder) containing training text files")
		return
	}

	// Ensure folder exists
	folderInfo, err := os.Stat(folderPath)
	if err != nil || !folderInfo.IsDir() {
		fmt.Printf("Error: Folder path '%s' does not exist or is not a directory\n", folderPath)
		return
	}

	fmt.Printf("Reading text files from folder: %s\n", folderPath)
	textSamples, filenames, err := parser.ReadTextFilesFromFolder(folderPath)
	if err != nil {
		fmt.Printf("Error reading text files: %v\n", err)
		return
	}

	if len(textSamples) == 0 {
		fmt.Println("Error: No .txt files found in the specified folder")
		return
	}

	fmt.Printf("Found %d text files for training\n", len(textSamples))
	if outputDetails {
		for i, filename := range filenames {
			fmt.Printf("  %d: %s (%d characters)\n", i+1, filename, len(textSamples[i]))
		}
	}

	// Parse each text sample
	parsedSamples := make([]string, len(textSamples))
	for i, sample := range textSamples {
		parsedSamples[i] = parser.ParseStringToAlphanumeric(sample)
	}

	fmt.Println("Creating distribution model...")
	model, err := analyzer.CreateDistributionFittedModel(parsedSamples, anomalyThreshold, fitThreshold)
	if err != nil {
		fmt.Printf("Error creating model: %v\n", err)
		return
	}

	if outputDetails {
		fmt.Println("\nModel Summary:")
		fmt.Println(model.GetModelSummary())
	}

	// Create directory if it doesn't exist
	modelDir := filepath.Dir(modelFilePath)
	if modelDir != "" && modelDir != "." {
		err = os.MkdirAll(modelDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory for model file: %v\n", err)
			return
		}
	}

	err = model.SaveTextModel(modelFilePath)
	if err != nil {
		fmt.Printf("Error saving model: %v\n", err)
		return
	}

	fmt.Printf("Model successfully created and saved to: %s\n", modelFilePath)
}

func useDistributionModel(modelFilePath string, checkTextFilePath string, outputDetails bool) {
	if modelFilePath == "" {
		fmt.Println("Error: You must specify a model file path (-model-file)")
		return
	}

	// Load the model
	fmt.Printf("Loading model from: %s\n", modelFilePath)
	model, err := analyzer.LoadTextModel(modelFilePath)
	if err != nil {
		fmt.Printf("Error loading model: %v\n", err)
		return
	}

	fmt.Println("Model loaded successfully")
	if outputDetails {
		fmt.Println("\nModel Summary:")
		fmt.Println(model.GetModelSummary())
	}

	if checkTextFilePath == "" {
		fmt.Println("\nNo text specified for checking. Use -check-text to analyze a sample against this model.")
		fmt.Println("Alternatively, you can input text directly:")
		fmt.Println("Enter text to check (type 'END' on a new line when finished):")

		inputText := parser.ReadMultilineInput()
		analyzeTextWithModel(model, inputText)
	} else {
		_, err := os.Stat(checkTextFilePath)
		if err != nil {
			fmt.Printf("Error: File '%s' does not exist or cannot be accessed\n", checkTextFilePath)
			return
		}

		fmt.Printf("Reading text from: %s\n", checkTextFilePath)
		textContent, err := parser.ReadFile(checkTextFilePath)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		analyzeTextWithModel(model, textContent)
	}
}

func analyzeTextWithModel(model *analyzer.TextDistributionFittedModel, text string) {
	parsedText := parser.ParseStringToAlphanumeric(text)

	isAnomalyFrequency, scoreFrequency, _, probabilityFrequency, isAnomalyPositions, scorePositions, _, probabilityPositions := model.IsAnomaly(parsedText)
	topAnomaliesFrequency, topAnomaliesPositions := model.GetTopAnomalies(parsedText, 10)

	fmt.Println("\n=========================")
	fmt.Println("Analysis Results Frequency")
	fmt.Println("=========================")

	if isAnomalyFrequency {
		fmt.Printf("ANOMALY DETECTED with score %.4f (threshold: %.4f)\n",
			scoreFrequency, model.AnomalyThreshold)
	} else {
		fmt.Printf("Text appears normal with score %.4f (threshold: %.4f)\n",
			scoreFrequency, model.AnomalyThreshold)
	}

	fmt.Printf("Probability: %.10f\n", probabilityFrequency)

	fmt.Println("\nTop anomalous characters:")
	if len(topAnomaliesFrequency) > 0 {
		for _, anomaly := range topAnomaliesFrequency {
			fmt.Printf("  %s\n", anomaly)
		}
	} else {
		fmt.Println("  No significant anomalies detected")
	}

	fmt.Println("\n=========================")
	fmt.Println("Analysis Results Positions")
	fmt.Println("=========================")

	if isAnomalyPositions {
		fmt.Printf("ANOMALY DETECTED with score %.4f (threshold: %.4f)\n",
			scorePositions, model.AnomalyThreshold)
	} else {
		fmt.Printf("Text appears normal with score %.4f (threshold: %.4f)\n",
			scorePositions, model.AnomalyThreshold)
	}

	fmt.Printf("Probability: %.10f\n", probabilityPositions)

	fmt.Println("\nTop anomalous characters:")
	if len(topAnomaliesPositions) > 0 {
		for _, anomaly := range topAnomaliesPositions {
			fmt.Printf("  %s\n", anomaly)
		}
	} else {
		fmt.Println("  No significant anomalies detected")
	}

	letterData := analyzer.AnalyzeLettersFromText(parsedText)

	fmt.Printf("Total characters: %d\n", letterData.TotalCount)
}
