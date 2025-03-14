package parser

import (
	"bufio"
	"os"
	"strings"
)

func ReadFile(path string) (string, error) {
	/*Read a .txt file and return its contents as a string
	 */
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadMultilineInput() string {
	/*Read multiline input from console until the text "END is found"
	 */
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
