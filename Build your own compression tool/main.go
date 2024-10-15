package main

import (
	"bufio"
	"container/heap"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

// Node represents a node in the Huffman tree
type Node struct {
	Char      rune
	Frequency int
	Left      *Node
	Right     *Node
}

// PriorityQueue implements a priority queue for Nodes
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Frequency < pq[j].Frequency
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// CalculateCharFrequency reads a file and returns a map of character frequencies
func CalculateCharFrequency(filename string) (map[rune]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frequencyMap := make(map[rune]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes) // Scan runes (characters)

	for scanner.Scan() {
		char := []rune(scanner.Text())[0]
		frequencyMap[char]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return frequencyMap, nil
}

// BuildHuffmanTree constructs the Huffman tree based on character frequencies
func BuildHuffmanTree(freqMap map[rune]int) *Node {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for char, freq := range freqMap {
		heap.Push(&pq, &Node{Char: char, Frequency: freq})
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)

		merged := &Node{
			Char:      0,
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}
		heap.Push(&pq, merged)
	}

	return heap.Pop(&pq).(*Node)
}

// GeneratePrefixCodeTable recursively generates the prefix-code table by traversing the tree
func GeneratePrefixCodeTable(node *Node, prefix string, table map[rune]string) {
	if node == nil {
		return
	}
	if node.Char != 0 {
		table[node.Char] = prefix
	} else {
		GeneratePrefixCodeTable(node.Left, prefix+"0", table)
		GeneratePrefixCodeTable(node.Right, prefix+"1", table)
	}
}

// WriteHeader writes the character frequency table as a header to the output file
func WriteHeader(frequencyMap map[rune]int, outputFile *os.File) error {
	encoder := gob.NewEncoder(outputFile)
	err := encoder.Encode(frequencyMap)
	if err != nil {
		return err
	}
	return nil
}

// PackBitsIntoBytes converts a bit string into bytes and writes them to the output file
func PackBitsIntoBytes(bitString string, outputFile *os.File) error {
	writer := bufio.NewWriter(outputFile)
	var currentByte byte
	var bitCount int

	for _, bit := range bitString {
		if bit == '1' {
			currentByte = (currentByte << 1) | 1
		} else {
			currentByte = currentByte << 1
		}
		bitCount++

		if bitCount == 8 {
			err := writer.WriteByte(currentByte)
			if err != nil {
				return err
			}
			currentByte = 0
			bitCount = 0
		}
	}

	// If there are remaining bits, pad with zeros and write the final byte
	if bitCount > 0 {
		currentByte = currentByte << (8 - bitCount)
		err := writer.WriteByte(currentByte)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// CompressFile compresses the file using the prefix-code table and writes the compressed data to the output file
func CompressFile(inputFile string, prefixCodeTable map[rune]string, outputFile *os.File) error {
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanRunes)

	var compressedBits strings.Builder

	for scanner.Scan() {
		char := []rune(scanner.Text())[0]
		compressedBits.WriteString(prefixCodeTable[char])
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	err = PackBitsIntoBytes(compressedBits.String(), outputFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide input and output filenames.")
		return
	}

	inputFilename := os.Args[1]
	outputFilename := os.Args[2]

	frequencyMap, err := CalculateCharFrequency(inputFilename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	huffmanTree := BuildHuffmanTree(frequencyMap)

	prefixCodeTable := make(map[rune]string)
	GeneratePrefixCodeTable(huffmanTree, "", prefixCodeTable)

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	err = WriteHeader(frequencyMap, outputFile)
	if err != nil {
		fmt.Printf("Error writing header: %v\n", err)
		return
	}

	err = CompressFile(inputFilename, prefixCodeTable, outputFile)
	if err != nil {
		fmt.Printf("Error compressing file: %v\n", err)
		return
	}

	fmt.Println("File compressed successfully.")
}
