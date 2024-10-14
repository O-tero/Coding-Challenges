package main

import (
	"fmt"
	"os"
)

func main() {
	// Check for input arguments (filename)
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as input.")
		return
	}

	filename := os.Args[1]

	// Step 1: Calculate character frequency
	frequencyMap, err := calculateCharFrequency(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Step 2: Build the Huffman Tree from character frequency
	huffmanTree := BuildHuffmanTree(frequencyMap)

	// Optionally, print the Huffman Tree for debugging
	fmt.Println("Huffman Tree:")
	PrintHuffmanTree(huffmanTree, "")
}
