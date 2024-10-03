package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	countBytes := flag.Bool("c", false, "count bytes")
	countWords := flag.Bool("w", false, "count words")
	countChars := flag.Bool("m", false, "count characters")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a filename.")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	// Read the file contents
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)

	// Check if the -c flag is set (count bytes)
	if *countBytes {
		fmt.Printf("%8d %s\n", len(data), filename)
	}

	// Check if the -w flag is set (count words)
	if *countWords {
		words := strings.Fields(content)
		fmt.Printf("%8d %s\n", len(words), filename)
	}

	// Check if the -m flag is set (count characters)
	if *countChars {
		charCount := utf8.RuneCountInString(content)
		fmt.Printf("%8d %s\n", charCount, filename)
	}

	if !*countBytes && !*countWords && !*countChars {
		fmt.Println("Usage: ccwc [-c|-w|-m] <filename>")
	}
}
