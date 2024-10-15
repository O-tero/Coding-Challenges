package main

import (
    "container/heap"
)

type Node struct {
    Char      rune
    Frequency int
    Left      *Node
    Right     *Node
}

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
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}

func BuildFrequencyMap(content string) map[rune]int {
    frequencies := make(map[rune]int)
    for _, char := range content {
        frequencies[char]++
    }
    return frequencies
}

func BuildHuffmanTree(frequencies map[rune]int) *Node {
    pq := &PriorityQueue{}
    heap.Init(pq)

    for char, freq := range frequencies {
        heap.Push(pq, &Node{Char: char, Frequency: freq})
    }

    for pq.Len() > 1 {
        left := heap.Pop(pq).(*Node)
        right := heap.Pop(pq).(*Node)

        merged := &Node{
            Char:      0, // Internal node
            Frequency: left.Frequency + right.Frequency,
            Left:      left,
            Right:     right,
        }
        heap.Push(pq, merged)
    }
    return heap.Pop(pq).(*Node)
}

func BuildPrefixTable(tree *Node) map[rune]string {
    prefixTable := make(map[rune]string)
    buildPrefixTableHelper(tree, "", prefixTable)
    return prefixTable
}

func buildPrefixTableHelper(node *Node, prefix string, prefixTable map[rune]string) {
    if node == nil {
        return
    }
    if node.Left == nil && node.Right == nil {
        prefixTable[node.Char] = prefix
    }
    buildPrefixTableHelper(node.Left, prefix+"0", prefixTable)
    buildPrefixTableHelper(node.Right, prefix+"1", prefixTable)
}
