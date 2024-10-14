# 

## Go JSON Parser

### Overview
This project provides a step-by-step guide to building a custom JSON parser in Go from scratch. The parser is designed to handle various JSON structures, including:

- Strings
- Numbers
- Booleans
- Null Values
- Nested Objects
- Arrays

### Features
- Lexical analyzer: Identifies JSON symbols, strings, numbers, and primitives.
- Parser: Validates JSON structure and produces an abstract syntax tree (AST).
- CLI Support: Verifies JSON validity from a string or file with error handling.

### Installation
1. Ensure Go is installed on your system.
2. Clone the repository: `git clone https://github.com/O-tero/go-json-parser.git`
3. Navigate to the project directory: `cd go-json-parser`

### Usage
**1. From Command Line:**
```
go run main.go '{"key": "value", "key2": 123}'
```

**2. From File:**
```
cat test.json | go run main.go
```

### Exit Codes
- 0: Valid JSON
- 1: Invalid JSON (stderr prints error)

### Steps Implemented
- **Step 1:** Parsing Empty JSON Objects
- **Step 2:** String Keys and Values
- **Step 3:** Different Data Types
- **Step 4:** Nested Objects and Arrays

