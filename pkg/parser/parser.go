package parser

import (
	"fmt"
	"strings"
	"unicode"
)

// Parses a string to alphanumeric characters and spaces only.
func ParseStringToAlphanumeric(textToParse string) string {
	var result strings.Builder

	for _, char := range textToParse {
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) {
			result.WriteRune(char)
		}
	}

	fmt.Printf("Final result: %v\n", result.String())
	return result.String()
}
