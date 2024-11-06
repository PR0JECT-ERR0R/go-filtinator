package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	var useClipboard string
	var clipboardContent string
	var err error

	// ask if the user wants to use clipboard content
	fmt.Print("Would you like to use clipboard content as input? (yes/no): ")
	// passes input into var
	fmt.Scan(&useClipboard)

	// initialize an empty slice to store lines
	lines := make([]string, 0)

	if strings.ToLower(useClipboard) == "yes" || strings.ToLower(useClipboard) == "y" {
		// attempt to get clipboard content
		clipboardContent, err = getClipboard()
		if err != nil {
			fmt.Println("Failed to read from clipboard:", err)
			return
		}

		// split clipboard content into lines
		lines = strings.Split(clipboardContent, "\n")
	} else {
		fmt.Println("Enter strings (type '!END' on a new line to finish):")
		// call function to get user input
		lines, err = getUserInput()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// get substring filter
	prompt := getPromptForSubstring()

	// filter duplicates and remove lines containing the prompt
	filteredLines := filterLines(lines, prompt)

	// count remaining lines
	lineCount := len(filteredLines)

	// output result
	fmt.Println("\nFiltered Lines:")
	for _, line := range filteredLines {
		fmt.Println(line)
	}
	fmt.Printf("\nTotal lines after filtering: %d\n", lineCount)
}

// getClipboard reads the clipboard content based on the OS
func getClipboard() (string, error) {
	var cmd *exec.Cmd

	// set the appropriate command based on OS
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard", "-o")
	case "darwin":
		cmd = exec.Command("pbpaste")
	case "windows":
		cmd = exec.Command("powershell", "Get-Clipboard")
	default:
		return "", fmt.Errorf("unsupported OS")
	}

	// run the command and capture the output
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error reading clipboard: %v", err)
	}

	return string(out), nil
}

// function to get user input
func getUserInput() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	// read lines from user input
	for scanner.Scan() {
		line := scanner.Text()
		if line == "!END" {
			break
		}
		lines = append(lines, line)
	}

	// handle any scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	return lines, nil
}

// function to get prompt to filter
func getPromptForSubstring() string {
	var prompt string
	fmt.Print("Enter substring to remove lines containing it (or press Enter to skip): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // read the whole line

	prompt = scanner.Text() // get the line from the scanner
	return prompt
}

// filterLines filters out duplicates and lines containing the prompt substring
func filterLines(lines []string, prompt string) []string {
	uniqueLines := make(map[string]bool)
	filteredLines := make([]string, 0)

	// iterate through lines and filter based on the prompt
	for _, line := range lines {
		// trim leading and trailing spaces
		line = strings.TrimSpace(line)

		if prompt != "" && strings.Contains(line, prompt) {
			continue // skip lines with the prompt substring
		}
		if !uniqueLines[line] {
			uniqueLines[line] = true
			filteredLines = append(filteredLines, line)
		}
	}

	return filteredLines
}
