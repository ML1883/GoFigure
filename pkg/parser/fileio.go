package parser

import (
	"bufio"
	"os"
	"path/filepath"
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

// ReadTextFilesFromFolder reads all .txt files from a folder and returns their contents as a slice of strings
func ReadTextFilesFromFolder(folderPath string) ([]string, []string, error) {
	var textContents []string
	var filenames []string

	// Read all files in the directory
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, nil, err
	}

	// Process each .txt file
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".txt") {
			filePath := filepath.Join(folderPath, file.Name())
			content, err := ReadFile(filePath)
			if err != nil {
				return nil, nil, err
			}

			// Add the file content to our slice
			textContents = append(textContents, content)
			filenames = append(filenames, file.Name())
		}
	}

	return textContents, filenames, nil
}
