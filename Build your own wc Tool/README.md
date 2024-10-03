# ccwc - A Custom `wc` Tool in Go

This is a custom version of the Unix `wc` (word count) tool, written in Go, called `ccwc`. The tool can count the number of lines, words, bytes, and characters in a file or from standard input.

## Features

- **`-c`**: Count the number of bytes in a file or input.
- **`-w`**: Count the number of words.
- **`-l`**: Count the number of lines.
- **`-m`**: Count the number of characters (supports multibyte characters).
- Supports **reading from standard input** when no filename is provided.
- **Efficient**: Uses Go's `bufio.Scanner` for line-by-line reading.

## Installation

To install the tool, you can clone this repository and build the Go project.

```bash
git clone <https://github.com/O-tero/Coding-Challenges.git>
cd Build your own wc tool
go build main.go
