// parser.go
package main

import (
	"fmt"
)

type Parser struct {
	lexer  *Lexer
	current Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.lexer.NextToken()
}

// Parse parses the input JSON
func (p *Parser) Parse() error {
	if p.current.Type == OpenBrace {
		p.nextToken() // Consume '{'
		return p.parseObject()
	}
	return fmt.Errorf("Invalid JSON")
}

func (p *Parser) parseObject() error {
	for p.current.Type != CloseBrace && p.current.Type != EOF {
		if p.current.Type != String {
			return fmt.Errorf("Expected string key")
		}
		p.nextToken() // Consume key

		if p.current.Type != Colon {
			return fmt.Errorf("Expected colon after key")
		}
		p.nextToken() // Consume ':'

		if err := p.parseValue(); err != nil {
			return err
		}

		if p.current.Type == Comma {
			p.nextToken() // Consume ','
		}
	}
	if p.current.Type != CloseBrace {
		return fmt.Errorf("Expected '}' at the end of object")
	}
	p.nextToken() // Consume '}'
	return nil
}

func (p *Parser) parseValue() error {
	switch p.current.Type {
	case String, Number, True, False, Null:
		p.nextToken() // Consume value
		return nil
	case OpenBrace:
		p.nextToken() // Consume '{'
		return p.parseObject()
	case 0:
		return fmt.Errorf("Unexpected end of input")
	default:
		return fmt.Errorf("unexpected token: %v", p.current)
	}
}
