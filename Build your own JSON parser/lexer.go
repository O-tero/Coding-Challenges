// lexer.go
package main

import (
	"unicode"
)

type TokenType int

const (
	Illegal TokenType = iota
	EOF
	Whitespace
	OpenBrace
	CloseBrace
	Colon
	Comma
	String
	Number
	True
	False
	Null
)

// Token represents a single token in the input
type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int  // current position in the input (points to current char)
	readPosition int  // current reading position (after current char)
	currentChar  rune // current char under examination
}

// NewLexer initializes a new lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character from the input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.currentChar = 0 // 0 represents EOF
	} else {
		l.currentChar = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var token Token

	switch l.currentChar {
	case '{':
		token = Token{Type: OpenBrace, Value: string(l.currentChar)}
	case '}':
		token = Token{Type: CloseBrace, Value: string(l.currentChar)}
	case ':':
		token = Token{Type: Colon, Value: string(l.currentChar)}
	case ',':
		token = Token{Type: Comma, Value: string(l.currentChar)}
	case 0:
		token.Type = EOF
	default:
		if unicode.IsSpace(l.currentChar) {
			token.Type = Whitespace
			token.Value = string(l.currentChar)
			l.readChar()
			return token
		} else if l.currentChar == '"' {
			token.Type = String
			token.Value = l.readString()
			return token
		} else if unicode.IsDigit(l.currentChar) || l.currentChar == '-' {
			token.Type = Number
			token.Value = l.readNumber()
			return token
		} else {
			// Handle true, false, and null
			switch l.currentChar {
			case 't':
				token.Type = True
				token.Value = l.readKeyword("true")
				return token
			case 'f':
				token.Type = False
				token.Value = l.readKeyword("false")
				return token
			case 'n':
				token.Type = Null
				token.Value = l.readKeyword("null")
				return token
			}
			token = Token{Type: Illegal, Value: string(l.currentChar)}
		}
	}
	l.readChar()
	return token
}

// readString reads a string token
func (l *Lexer) readString() string {
	start := l.position + 1
	l.readChar()
	for l.currentChar != '"' && l.currentChar != 0 {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber reads a number token
func (l *Lexer) readNumber() string {
	start := l.position
	for unicode.IsDigit(l.currentChar) || l.currentChar == '.' {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readKeyword reads a keyword (true, false, null)
func (l *Lexer) readKeyword(keyword string) string {
	start := l.position
	for l.position < len(l.input) && l.currentChar == rune(keyword[len(keyword)-1]) {
		l.readChar()
	}
	return l.input[start:l.position]
}
