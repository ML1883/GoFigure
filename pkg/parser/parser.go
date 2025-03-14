package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseStringToAlphanumeric(textToParse string) string {
	/*Parses a string to alphanumeric characters and spaces only.
	 */
	var result strings.Builder

	for _, char := range textToParse {
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) {
			result.WriteRune(char)
			// fmt.Printf("Character appended: %c at index %v\n", char, index)
		}
	}

	fmt.Printf("Final result: %v\n", result.String())
	return result.String()
}
