package main

import (
	"bufio"
	"os"
)

// CalculateCharFrequency reads a file and returns a map of character frequencies
func CalculateCharFrequency(filename string) (map[rune]int, error) { 
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frequencyMap := make(map[rune]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes) // Scan runes (characters)

	for scanner.Scan() {
		char := []rune(scanner.Text())[0]
		frequencyMap[char]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return frequencyMap, nil
}
