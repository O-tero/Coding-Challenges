# Huffman Compression Tool

This is a command-line tool for compressing and decompressing files using Huffman coding. It leverages Huffman trees to reduce file sizes while allowing for efficient data retrieval.

## Features

- Compress files to reduce their size.
- Decompress previously compressed files, restoring the original content.
- Uses Huffman coding for efficient data representation.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.15 or later) must be installed on your machine.

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/O-tero/huffman-compression-tool.git
   cd huffman-compression-tool


### Build the application
```bash
    go build -o huffman


### Command Line Arguments
- compress: Flag to indicate that you want to compress a file.
- decompress: Flag to indicate that you want to decompress a file.
- input: The path to the input file that you want to process.
-o utput: The path where the output file will be saved.

Example Commands
To compress a file:
```bash
    ./huffman -compress -input input.txt -output compressed.huff

To decompress a file:
```bash
    ./huffman -decompress -input compressed.huff -output output.txt

### How It Works - Compressing a File:

- Read the input file and build a frequency map of the characters.
- Construct a Huffman tree based on the frequency map.
- Generate a prefix table that maps characters to their corresponding binary codes.
- Pack the original content using the prefix table and write the compressed data to the specified output file along with the Huffman tree.

Decompressing a File:

- Read the compressed file to retrieve the Huffman tree and the packed bytes.
- Decode the packed bytes using the Huffman tree and write the decompressed data to the specified output file.