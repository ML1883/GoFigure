package analyzer

import (
	"fmt"
	"math"
	"unicode"
)

type lettercount struct {
	TotalCount        int
	LetterNumberArray [36]int //0-9 + 26 letters
	PositionArray     [36][]int
}

func CountLetters(textToCount string) *lettercount {
	//Take text and return adress of struct
	lcText := lettercount{}
	for index, char := range textToCount {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			if unicode.IsLetter(char) {
				var letterLower rune = unicode.ToLower(char)
				var letterNumber int = int(letterLower - 'a')
				lcText.LetterNumberArray[letterNumber+10]++
				lcText.PositionArray[letterNumber+10] = append(lcText.PositionArray[letterNumber+10], index)
				fmt.Printf("For index %v, plussed position %v with one with character %c\n", index, letterNumber+10, char)
			} else {
				lcText.LetterNumberArray[int(char-'0')]++
				lcText.PositionArray[int(char-'0')] = append(lcText.PositionArray[int(char-'0')], index)
				fmt.Printf("For index %v, plussed position %v with one with character %c\n", index, int(char), char)

			}

		}
		lcText.TotalCount++
	}

	return &lcText
}

func IntVectorMultiplication(array1 []int, array2 []int) (int, error) {
	//Optimized function for vector multiplication
	var lengthArray1 int = len(array1)
	var lengthArray2 int = len(array2)

	if lengthArray1 > 0 && lengthArray2 > 0 && lengthArray1 == lengthArray2 {
		var total int = 0
		for index := range array1 {
			total += (array1[index] * array2[index])
			fmt.Printf("Array 1 value: %v and array2 value: %v\n", array1[index], array2[index])
		}
		return total, nil
	} else {
		return 0, fmt.Errorf("array with zero length or unequal array length detected. Length array1: %v Length array2: %v", lengthArray1, lengthArray2)
	}

}

func CosineSimilarityVectors(array1 []int, array2 []int) float64 {
	/*Performs consine similarity calculation on two arrays of integers.
	Resulting in a cosine similairity. We do not check if the arrays have correct sizes.
	TODO: make this incorporate the position somehow*/
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
