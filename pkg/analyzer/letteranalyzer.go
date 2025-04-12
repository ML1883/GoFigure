package analyzer

import (
	"fmt"
	"math"
	"unicode"
)

type LetterData struct {
	TotalCount        int
	LetterCount       int
	LetterNumberArray [36]int //0-9 + 26 letters
	PositionArray     [36][]int
}

// Takes text and return adress of the letterdata struct
func AnalyzeLettersFromText(textToCount string) *LetterData {
	lcText := LetterData{}
	for index, char := range textToCount {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			if unicode.IsLetter(char) {
				var letterLower rune = unicode.ToLower(char)
				var letterNumber int = int(letterLower - 'a')
				lcText.LetterNumberArray[letterNumber+10]++
				lcText.PositionArray[letterNumber+10] = append(lcText.PositionArray[letterNumber+10], index)
			} else {
				lcText.LetterNumberArray[int(char-'0')]++
				lcText.PositionArray[int(char-'0')] = append(lcText.PositionArray[int(char-'0')], index)

			}
			lcText.LetterCount++
		}
		lcText.TotalCount++
	}

	return &lcText
}

// Function for vector multiplication for our specific use case
func IntVectorMultiplication(array1 []int, array2 []int) (int, error) {

	var lengthArray1 int = len(array1)
	var lengthArray2 int = len(array2)

	if lengthArray1 > 0 && lengthArray2 > 0 && lengthArray1 == lengthArray2 {
		var total int = 0
		for index := range array1 {
			total += (array1[index] * array2[index])
		}
		return total, nil
	} else {
		return 0, fmt.Errorf("array with zero length or unequal array length detected. Length array1: %v Length array2: %v", lengthArray1, lengthArray2)
	}

}

// Performs consine similarity calculation on two arrays of integers.
// Resulting in a cosine similairity. We do not check if the arrays have correct sizes.
// Return range: [-1,1].
// Where 1 is complete similairity, 0 is no similairity, -1 is complete opposites
func CosineSimilarityVectors(array1 []int, array2 []int) float64 {
	var dotProduct int = 0
	var magnitudeArray1 int = 0
	var magnitudeArray2 int = 0

	for i := range array1 {
		dotProduct += array1[i] * array2[i]
		magnitudeArray1 += array1[i] * array1[i]
		magnitudeArray2 += array2[i] * array2[i]
	}

	return float64(dotProduct) / (math.Sqrt(float64(magnitudeArray1)) * math.Sqrt(float64(magnitudeArray2))) //Can be above two because of float magic

}

// Calculate the Jaccard idnex of two arrays
// Return range: [0,1]
// Where 0 is not overlap at all
// and 1 is complete overlap
func JaccardIndexVectors(array1 []int, array2 []int) float64 {
	var intersection int = 0
	var union int = 0

	for i := range array1 {
		intersection += min(array1[i], array2[i]) //overlap
		union += max(array1[i], array2[i])        //area
	}

	if union == 0 {
		return 1
	}

	return float64(intersection) / float64(union)
}

// Calculate average difference of each number or letter.
// Calculate the average of that difference across the whole spectrum of numbers and letters.
// Divide this difference over the max length of either one or two, minus one to normalize it.
// Return range: [0,1]
// Where 0 is complete similairity of positions
// And nearing 1 is no similairity of positions.
// The larger a string is, the greater the chance that it nears one.
// That still needs to be fixed somehow.
func PositionDifferenceVectors(array1 [][]int, array2 [][]int, totalLength1 int, totalLength2 int) float64 {

	var lengthArray1 int = 0
	var lengthArray2 int = 0
	var totalAvgDifference float64 = 0
	var elementsCalculated int = 0

	for i := range array1 {
		var remainderTotal float64 = 0
		var totalAbsDiff float64 = 0
		subArray1 := array1[i]
		subArray2 := array2[i]
		lengthArray1 = len(subArray1)
		lengthArray2 = len(subArray2)

		if lengthArray1 == 0 && lengthArray2 == 0 {
			continue
		}

		elementsCalculated++
		switch {
		case lengthArray1 == 0: //Calculate the average position and add unchallenged
			for j := range lengthArray2 {
				remainderTotal += float64(subArray2[j])
			}
			totalAvgDifference += (remainderTotal / float64(lengthArray2))

		case lengthArray2 == 0:
			for j := range lengthArray1 {
				remainderTotal += float64(subArray1[j])
			}
			totalAvgDifference += (remainderTotal / float64(lengthArray1))

		case lengthArray1 > lengthArray2:
			for j := range lengthArray2 {
				totalAbsDiff += math.Abs(float64(subArray1[j]) - float64(subArray2[j]))
			}
			totalAvgDifference += (totalAbsDiff / float64(lengthArray2))

			for j := lengthArray2; j < lengthArray1; j++ {
				remainderTotal += float64(subArray1[j])
			}
			totalAvgDifference += (remainderTotal / float64(lengthArray1-lengthArray2))

		case lengthArray1 < lengthArray2:
			for j := range lengthArray1 {
				totalAbsDiff += math.Abs(float64(subArray1[j]) - float64(subArray2[j]))
			}
			totalAvgDifference += (totalAbsDiff / float64(lengthArray1))

			for j := lengthArray1; j < lengthArray2; j++ {
				remainderTotal += float64(subArray2[j])
			}
			totalAvgDifference += (remainderTotal / float64(lengthArray2-lengthArray1))

		case lengthArray1 == lengthArray2:
			for j := range lengthArray1 {
				totalAbsDiff += math.Abs(float64(subArray1[j]) - float64(subArray2[j]))
			}
			totalAvgDifference += (totalAbsDiff / float64(lengthArray1))

		default:
			continue
		}

	}
	var grandTotalAvgDifference float64 = totalAvgDifference / float64(elementsCalculated)
	return grandTotalAvgDifference / float64(max(totalLength1, totalLength2)-1)
}
