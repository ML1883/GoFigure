package parser

import (
	"bufio"
	"os"
	"strings"
)

// Read a .txt file and return its contents as a string
func ReadFile(path string) (string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Read multiline input from console until the text "END" is found
func ReadMultilineInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
