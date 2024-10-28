// // package main

// import (
// 	"bufio"
// 	"container/heap"
// 	"encoding/gob"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"
// )

// // Node represents a node in the Huffman tree
// type Node struct {
// 	Char      rune
// 	Frequency int
// 	Left      *Node
// 	Right     *Node
// }

// // PriorityQueue implements a priority queue for Nodes
// type PriorityQueue []*Node

// func (pq PriorityQueue) Len() int           { return len(pq) }
// func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Frequency < pq[j].Frequency }
// func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
// func (pq *PriorityQueue) Push(x interface{}) {
// 	*pq = append(*pq, x.(*Node))
// }
// func (pq *PriorityQueue) Pop() interface{} {
// 	old := *pq
// 	n := len(old)
// 	x := old[n-1]
// 	*pq = old[0 : n-1]
// 	return x
// }

// // CalculateCharFrequency reads a file and returns a map of character frequencies
// func CalculateCharFrequency(filename string) (map[rune]int, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	frequencyMap := make(map[rune]int)
// 	scanner := bufio.NewScanner(file)
// 	scanner.Split(bufio.ScanRunes)
// 	for scanner.Scan() {
// 		char := []rune(scanner.Text())[0]
// 		frequencyMap[char]++
// 	}

// 	return frequencyMap, scanner.Err()
// }

// // BuildHuffmanTree constructs the Huffman tree based on character frequencies
// func BuildHuffmanTree(freqMap map[rune]int) *Node {
// 	pq := make(PriorityQueue, 0, len(freqMap))
// 	heap.Init(&pq)

// 	for char, freq := range freqMap {
// 		heap.Push(&pq, &Node{Char: char, Frequency: freq})
// 	}

// 	for pq.Len() > 1 {
// 		left := heap.Pop(&pq).(*Node)
// 		right := heap.Pop(&pq).(*Node)
// 		merged := &Node{Frequency: left.Frequency + right.Frequency, Left: left, Right: right}
// 		heap.Push(&pq, merged)
// 	}

// 	return heap.Pop(&pq).(*Node)
// }

// // GeneratePrefixCodeTable recursively generates the prefix-code table
// func GeneratePrefixCodeTable(node *Node, prefix string, table map[rune]string) {
// 	if node == nil {
// 		return
// 	}
// 	if node.Char != 0 {
// 		table[node.Char] = prefix
// 	} else {
// 		GeneratePrefixCodeTable(node.Left, prefix+"0", table)
// 		GeneratePrefixCodeTable(node.Right, prefix+"1", table)
// 	}
// }

// // WriteHeader writes the frequency map as a header to the output file
// func WriteHeader(freqMap map[rune]int, outputFile *os.File) error {
// 	encoder := gob.NewEncoder(outputFile)
// 	return encoder.Encode(freqMap)
// }

// // ReadHeader reads the frequency map from the header of the encoded file
// func ReadHeader(inputFile *os.File) (map[rune]int, error) {
// 	var frequencyMap map[rune]int
// 	decoder := gob.NewDecoder(inputFile)
// 	if err := decoder.Decode(&frequencyMap); err != nil {
// 		return nil, err
// 	}
// 	return frequencyMap, nil
// }

// // PackBitsIntoBytes converts a bit string into bytes and writes them to the output file
// func PackBitsIntoBytes(bitString string, outputFile *os.File) error {
// 	writer := bufio.NewWriter(outputFile)
// 	var currentByte byte
// 	bitCount := 0

// 	for _, bit := range bitString {
// 		if bit == '1' {
// 			currentByte = (currentByte << 1) | 1
// 		} else {
// 			currentByte = currentByte << 1
// 		}
// 		bitCount++

// 		// If we have accumulated 8 bits, write the byte to the file
// 		if bitCount == 8 {
// 			if err := writer.WriteByte(currentByte); err != nil {
// 				return err
// 			}
// 			currentByte = 0
// 			bitCount = 0
// 		}
// 	}

// 	if bitCount > 0 {
// 		currentByte <<= (8 - bitCount)
// 		if err := writer.WriteByte(currentByte); err != nil {
// 			return err
// 		}
// 	}

// 	return writer.Flush()
// }

// // RebuildHuffmanTreeFromHeader reads the header from the file and rebuilds the Huffman tree
// func RebuildHuffmanTreeFromHeader(inputFile *os.File) (*Node, error) {
// 	frequencyMap, err := ReadHeader(inputFile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return BuildHuffmanTree(frequencyMap), nil
// }

// func CompressFile(inputFilename, outputFilename string) error {
//     frequencyMap, err := CalculateCharFrequency(inputFilename)
//     if err != nil {
//         return err
//     }

//     huffmanTree := BuildHuffmanTree(frequencyMap)

//     prefixCodeTable := make(map[rune]string)
//     GeneratePrefixCodeTable(huffmanTree, "", prefixCodeTable)

//     outputFile, err := os.Create(outputFilename)
//     if err != nil {
//         return err
//     }
//     defer outputFile.Close()

//     // Ensure consistent type: map[rune]int
//     enc := gob.NewEncoder(outputFile)
//     if err := enc.Encode(frequencyMap); err != nil {
//         return fmt.Errorf("error encoding frequency map: %v", err)
//     }

//     inputFile, err := os.Open(inputFilename)
//     if err != nil {
//         return err
//     }
//     defer inputFile.Close()

//     var bitString strings.Builder
//     scanner := bufio.NewScanner(inputFile)
//     scanner.Split(bufio.ScanRunes)

//     for scanner.Scan() {
//         char := []rune(scanner.Text())[0]
//         bitString.WriteString(prefixCodeTable[char])
//     }

//     if err := scanner.Err(); err != nil {
//         return err
//     }

//     return PackBitsIntoBytes(bitString.String(), outputFile)
// }


// func DecompressFile(inputFilename, outputFilename string) error {
//     inputFile, err := os.Open(inputFilename)
//     if err != nil {
//         return err
//     }
//     defer inputFile.Close()

//     // Ensure consistent type: map[rune]int
//     var frequencyMap map[rune]int
//     dec := gob.NewDecoder(inputFile)
//     if err := dec.Decode(&frequencyMap); err != nil {
//         return fmt.Errorf("error decoding frequency map: %v", err)
//     }

//     huffmanTree := BuildHuffmanTree(frequencyMap)

//     outputFile, err := os.Create(outputFilename)
//     if err != nil {
//         return err
//     }
//     defer outputFile.Close()

//     err = UnpackBytesAndDecode(inputFile, outputFile, huffmanTree)
//     if err != nil {
//         return fmt.Errorf("error during decompression: %v", err)
//     }

//     fmt.Println("Decompression completed successfully")
//     return nil
// }

// // UnpackBytesAndDecode reads the packed bytes from the input file and decodes them using the Huffman tree.
// // It writes the decompressed characters to the output file.
// func UnpackBytesAndDecode(inputFile *os.File, outputFile *os.File, huffmanTree *Node) error {
//     var bitString strings.Builder

//     // Read the packed bytes from the input file
//     byteBuffer := make([]byte, 1)
//     for {
//         n, err := inputFile.Read(byteBuffer)
//         if err != nil {
//             if err == io.EOF {
//                 break
//             }
//             return err
//         }
//         if n > 0 {
//             // Convert each byte to a bit string (e.g., "10101010")
//             bitString.WriteString(fmt.Sprintf("%08b", byteBuffer[0]))
//         }
//     }

//     // Now we have the full bit string, we need to decode it using the Huffman Tree
//     currentNode := huffmanTree
//     for _, bit := range bitString.String() {
//         if bit == '0' {
//             currentNode = currentNode.Left
//         } else {
//             currentNode = currentNode.Right
//         }

//         // If we hit a leaf node, we have found a character
//         if currentNode.Left == nil && currentNode.Right == nil {
//             // Write the decoded character to the output file
// 			_, err := outputFile.WriteString(string(currentNode.Char))
//             if err != nil {
//                 return err
//             }
//             // Reset to the root of the tree for the next character
//             currentNode = huffmanTree
//         }
//     }

//     return nil
// }

// func main() {
// 	if len(os.Args) < 4 {
// 		fmt.Println("Usage: [compress|decompress] <input> <output>")
// 		return
// 	}

// 	mode := os.Args[1]
// 	inputFilename := os.Args[2]
// 	outputFilename := os.Args[3]

// 	if mode == "compress" {
// 		err := CompressFile(inputFilename, outputFilename)
// 		if err != nil {
// 			fmt.Printf("Compression error: %v\n", err)
// 		} else {
// 			fmt.Println("File compressed successfully.")
// 		}
// 	} else if mode == "decompress" {
// 		err := DecompressFile(inputFilename, outputFilename)
// 		if err != nil {
// 			fmt.Printf("Decompression error: %v\n", err)
// 		} else {
// 			fmt.Println("File decompressed successfully.")
// 		}
// 	} else {
// 		fmt.Println("Unknown mode. Use 'compress' or 'decompress'.")
// 	}
// }
