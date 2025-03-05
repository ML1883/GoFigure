package analyzer

import (
	"fmt"
	"unicode"
)

type lettercount struct {
	totalcount        int
	letterNumberArray [36]int //0-9 + 26 letters
}

func CountLetters(textToCount string) *lettercount {
	//Take text and return adress of struct
	lcText := lettercount{}
	for index, char := range textToCount {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			if unicode.IsLetter(char) {
				var letterLower rune = unicode.ToLower(char)
				var letterNumber int = int(letterLower) //This is not correct for our purpose.
				lcText.letterNumberArray[letterNumber+10]++
				fmt.Printf("For index %v, plussed position %v with one with character %v", index, letterNumber+10, char)
			} else {
				lcText.letterNumberArray[int(char)]++
				fmt.Printf("For index %v, plussed position %v with one with character %v", index, int(char), char)

			}

		}
		lcText.totalcount++
	}

	return &lcText
}
