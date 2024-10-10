package main

import (
	"fmt"
	"os"
)

// Token types
const (
	TokenCurlyOpen = iota
	TokenCurlyClose
	TokenEOF
	TokenInvalid
)

// Token structure
type Token struct {
	Type  int
	Value string
}

// Lexer struct
type Lexer struct {
	input string
	pos   int
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	switch l.input[l.pos] {
	case '{':
		l.pos++
		return Token{Type: TokenCurlyOpen, Value: "{"}
	case '}':
		l.pos++
		return Token{Type: TokenCurlyClose, Value: "}"}
	default:
		return Token{Type: TokenInvalid, Value: string(l.input[l.pos])}
	}
}

// Parser struct
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

// NewParser creates a new parser
func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer}
}

// Parse parses the JSON input
func (p *Parser) Parse() bool {
	p.currentToken = p.lexer.NextToken()

	if p.currentToken.Type == TokenCurlyOpen {
		p.currentToken = p.lexer.NextToken()
		if p.currentToken.Type == TokenCurlyClose {
			return true // Valid JSON: {}
		}
	}
	return false // Invalid JSON
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run json_parser.go <json_string>")
		os.Exit(1)
	}

	input := os.Args[1]
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	if parser.Parse() {
		fmt.Println("Valid JSON")
		os.Exit(0) // Valid JSON
	} else {
		fmt.Println("Invalid JSON")
		os.Exit(1) // Invalid JSON
	}
}
