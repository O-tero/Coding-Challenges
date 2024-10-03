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
	countLines := flag.Bool("l", false, "count lines")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Please provide a filename.")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)

	lineCount := len(strings.Split(content, "\n"))

	wordCount := len(strings.Fields(content))

	byteCount := len(data)

	charCount := utf8.RuneCountInString(content)

	if *countBytes {
		fmt.Printf("%8d %s\n", byteCount, filename)
	}

	if *countWords {
		fmt.Printf("%8d %s\n", wordCount, filename)
	}

	if *countChars {
		fmt.Printf("%8d %s\n", charCount, filename)
	}

	if *countLines {
		fmt.Printf("%8d %s\n", lineCount, filename)
	}

	if !*countBytes && !*countWords && !*countChars && !*countLines {
		fmt.Printf("%8d %8d %8d %s\n", lineCount, wordCount, byteCount, filename)
	}
}
