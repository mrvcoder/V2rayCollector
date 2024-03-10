package collector

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func ChangeUrlToTelegramWebUrl(input string) string {
	// Check if the input URL already contains "/s/", if not, add it
	if !strings.Contains(input, "/s/") {
		// Find the position of "/t.me/" in the URL
		index := strings.Index(input, "/t.me/")
		if index != -1 {
			// Insert "/s/" after "/t.me/"
			modifiedURL := input[:index+len("/t.me/")] + "s/" + input[index+len("/t.me/"):]
			return modifiedURL
		}
	}
	// If "/s/" already exists or "/t.me/" is not found, return the input as is
	return input
}

func ReadFileContent(filePath string) (string, error) {
	// Read the entire file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Convert the content to a string and return
	return string(content), nil
}

func Reverse(lines []string) []string {
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - i - 1
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}

func RemoveDuplicate(config string) string {
	lines := strings.Split(config, "\n")
	slices.Sort(lines)
	lines = slices.Compact(lines)
	// Join unique lines into a string
	uniqueString := strings.Join(lines, "\n")
	return uniqueString
}

func WriteToFile(fileContent string, filePath string) {

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// If the file exists, clear its content
		err = os.WriteFile(filePath, []byte{}, 0644)
		if err != nil {
			fmt.Println("Error clearing file:", err)
			return
		}
	} else if os.IsNotExist(err) {
		// If the file does not exist, create it
		_, err = os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	} else {
		// If there was some other error, print it and return
		fmt.Println("Error checking file:", err)
		return
	}

	// Write the new content to the file
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully")
}
