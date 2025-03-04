package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func ParseStringToAlphanumeric(input string) string {
	//Parses a string to Alphanumeric characters and spaces only.
	var result strings.Builder

	for index, char := range input {
		// Allow alphanumeric characters and spaces
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) {
			result.WriteRune(char)
			fmt.Printf("Character appended: %v", index)
		}
	}

	fmt.Printf("Final result: %v", result.String())
	return result.String()
}
