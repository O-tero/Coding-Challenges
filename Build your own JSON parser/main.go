// main.go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <json-string>")
		return
	}

	input := os.Args[1]
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	if err := parser.Parse(); err != nil {
		fmt.Println("Invalid JSON:", err)
		os.Exit(1)
	}
	fmt.Println("Valid JSON")
	os.Exit(0)
}
