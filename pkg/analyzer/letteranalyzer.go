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
				// fmt.Printf("For index %v, plussed position %v with one with character %c\n", index, letterNumber+10, char)
			} else {
				lcText.LetterNumberArray[int(char-'0')]++
				lcText.PositionArray[int(char-'0')] = append(lcText.PositionArray[int(char-'0')], index)
				// fmt.Printf("For index %v, plussed position %v with one with character %c\n", index, int(char), char)

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
			// fmt.Printf("Array 1 value: %v and array2 value: %v\n", array1[index], array2[index])
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

func JaccardIndexVectors(array1 []int, array2 []int) float64 {
	/*Calculate the Jaccard idnex of two arrays
	No error handling.*/
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

func PositionDifferenceVectors(array1 [][]int, array2 [][]int, totalLength int) float64 {
	/*
		There aren't really any specific formulas for this problem so we will have to use the following algorithm:
		1. Check if a letter is present in array1 or array2
		2. If it is present in both, calculate difference
		3. If it is present in just one, leave it(??)
		4. If the sizes differ, average out the positions that 'stick out'
		5. If the sizes don't differ, just calculate differences (absolute numbers here or squaring of numbers depending on how much we want to punish letters being on outlier positions)
		6. If its present in none we leave it.
		7. Normalize the difference calculations using the total length of both (or single) arrays such that the difference becomes reduced to an index.
			-> we will have to find out what the index is that we want here

		Maybe do something with the frequency calc as well to show importance akin to levensthein Algo? Or perhaps include this in different version.
	*/

	// var averageLeftover float64 = 0
	var lengthArray1 int = 0
	var lengthArray2 int = 0
	var totalDifference float64 = 0
	//var letterCount int = 0 //How many letters do w have

	for i := range array1 {
		subArray1 := array1[i]
		subArray2 := array2[i]
		lengthArray1 = len(subArray1)
		lengthArray2 = len(subArray2)

		if lengthArray1 == 0 && lengthArray2 == 0 {
			continue
		}

		switch {
		case lengthArray1 > lengthArray2:
			continue
		case lengthArray1 < lengthArray2:
			continue
		case lengthArray1 == lengthArray2:
			for j := range lengthArray1 {
				totalDifference += math.Abs(float64(subArray1[j]) - float64(subArray2[j]))
			}
		default:
			continue
		}

	}

	return float64(totalDifference) / float64(totalLength)
}
