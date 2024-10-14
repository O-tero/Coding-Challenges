package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
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
	scanner.Split(bufio.ScanRunes) 
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

// PrintHuffmanTree prints the Huffman tree structure
func PrintHuffmanTree(node *Node, prefix string) {
	if node == nil {
		return
	}
	if node.Char != 0 {
		fmt.Printf("%q: %d (%s)\n", node.Char, node.Frequency, prefix)
	} else {
		fmt.Printf("Node: %d\n", node.Frequency)
	}
	PrintHuffmanTree(node.Left, prefix+"0")
	PrintHuffmanTree(node.Right, prefix+"1")
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename as input.")
		return
	}

	filename := os.Args[1]

	frequencyMap, err := CalculateCharFrequency(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	huffmanTree := BuildHuffmanTree(frequencyMap)

	fmt.Println("Huffman Tree:")
	PrintHuffmanTree(huffmanTree, "")

	prefixCodeTable := make(map[rune]string)
	GeneratePrefixCodeTable(huffmanTree, "", prefixCodeTable)

	fmt.Println("\nPrefix-Code Table:")
	for char, code := range prefixCodeTable {
		fmt.Printf("%q: %s\n", char, code)
	}
}
