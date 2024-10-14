// package main

// import (
// 	"container/heap"
// 	"fmt"
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

// func (pq PriorityQueue) Len() int { return len(pq) }
// func (pq PriorityQueue) Less(i, j int) bool {
// 	return pq[i].Frequency < pq[j].Frequency
// }
// func (pq PriorityQueue) Swap(i, j int) {
// 	pq[i], pq[j] = pq[j], pq[i]
// }
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

// // BuildHuffmanTree constructs the Huffman tree based on character frequencies
// func BuildHuffmanTree(freqMap map[rune]int) *Node { // Capitalized
// 	pq := make(PriorityQueue, 0)
// 	heap.Init(&pq)

// 	for char, freq := range freqMap {
// 		heap.Push(&pq, &Node{Char: char, Frequency: freq})
// 	}

// 	for pq.Len() > 1 {
// 		left := heap.Pop(&pq).(*Node)
// 		right := heap.Pop(&pq).(*Node)

// 		merged := &Node{
// 			Char:      0,
// 			Frequency: left.Frequency + right.Frequency,
// 			Left:      left,
// 			Right:     right,
// 		}
// 		heap.Push(&pq, merged)
// 	}

// 	return heap.Pop(&pq).(*Node)
// }

// // PrintHuffmanTree prints the Huffman tree structure
// func PrintHuffmanTree(node *Node, prefix string) { 
// 	if node == nil {
// 		return
// 	}
// 	if node.Char != 0 {
// 		fmt.Printf("%q: %d\n", node.Char, node.Frequency)
// 	} else {
// 		fmt.Printf("Node: %d\n", node.Frequency)
// 	}
// 	PrintHuffmanTree(node.Left, prefix+"0")
// 	PrintHuffmanTree(node.Right, prefix+"1")
// }
