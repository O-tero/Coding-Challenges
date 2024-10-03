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

	var data []byte
	var err error
	if len(flag.Args()) == 0 {
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Failed to read from standard input: %v", err)
		}
	} else {
		filename := flag.Arg(0)
		data, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
	}

	content := string(data)

	lineCount := len(strings.Split(content, "\n"))

	wordCount := len(strings.Fields(content))

	byteCount := len(data)

	charCount := utf8.RuneCountInString(content)

	if *countBytes {
		fmt.Printf("%8d\n", byteCount)
	}

	if *countWords {
		fmt.Printf("%8d\n", wordCount)
	}

	if *countChars {
		fmt.Printf("%8d\n", charCount)
	}

	if *countLines {
		fmt.Printf("%8d\n", lineCount)
	}

	if !*countBytes && !*countWords && !*countChars && !*countLines {
		fmt.Printf("%8d %8d %8d\n", lineCount, wordCount, byteCount)
	}
}
