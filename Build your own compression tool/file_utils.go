package main

import (
    "encoding/binary"
    "encoding/gob"
    "os"
)

func WriteCompressedFile(outputPath string, huffmanTree *Node, packedBytes []byte) error {
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the Huffman tree to the file
    if err := gob.NewEncoder(file).Encode(huffmanTree); err != nil {
        return err
    }

    // Write the packed bytes to the file
    packedSize := int64(len(packedBytes))
    if err := binary.Write(file, binary.BigEndian, packedSize); err != nil {
        return err
    }
    if _, err := file.Write(packedBytes); err != nil {
        return err
    }

    return nil
}

func ReadCompressedFile(inputFile *os.File) (*Node, []byte, error) {
    huffmanTree := new(Node)
    if err := gob.NewDecoder(inputFile).Decode(huffmanTree); err != nil {
        return nil, nil, err
    }

    var packedSize int64
    if err := binary.Read(inputFile, binary.BigEndian, &packedSize); err != nil {
        return nil, nil, err
    }

    packedBytes := make([]byte, packedSize)
    if _, err := inputFile.Read(packedBytes); err != nil {
        return nil, nil, err
    }

    return huffmanTree, packedBytes, nil
}

func UnpackBytesAndDecode(packedBytes []byte, outputFile *os.File, huffmanTree *Node) error {
    currentNode := huffmanTree
    var decodedRunes []rune

    for _, b := range packedBytes {
        for i := 0; i < 8; i++ {
            bit := (b >> (7 - i)) & 1
            if bit == 0 {
                currentNode = currentNode.Left
            } else {
                currentNode = currentNode.Right
            }

            if currentNode.Left == nil && currentNode.Right == nil {
                decodedRunes = append(decodedRunes, currentNode.Char)
                currentNode = huffmanTree
            }
        }
    }

    _, err := outputFile.WriteString(string(decodedRunes))
    return err
}
