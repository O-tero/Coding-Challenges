package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

// Function to calculate character frequencies
func calculateCharFrequency(filename string) (map[rune]int, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Initialize frequency map
	frequencyMap := make(map[rune]int)

	// Read the file content
	reader := bufio.NewReader(file)
	for {
		char, _, err := reader.ReadRune() 
		if err != nil {
			break 
		}
		if unicode.IsPrint(char) || unicode.IsSpace(char) {
			frequencyMap[char]++
		}
	}

	return frequencyMap, nil
}

// Function to print the character frequency map
func printFrequencyTable(frequencyMap map[rune]int) {
	fmt.Println("Character Frequency Table:")
	for char, freq := range frequencyMap {
		fmt.Printf("%q: %d\n", char, freq)
	}
}

func main() {
	// Check for input arguments (filename)
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as input.")
		return
	}

	filename := os.Args[1]

	// Calculate character frequency
	frequencyMap, err := calculateCharFrequency(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// print the frequency table for debugging
	printFrequencyTable(frequencyMap)

	
}
