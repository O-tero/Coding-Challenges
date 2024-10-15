package main

import (
    "io/ioutil"
    "os"
)

func CompressFile(inputPath string, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    content, err := ioutil.ReadAll(inputFile)
    if err != nil {
        return err
    }

    frequencies := BuildFrequencyMap(string(content))
    huffmanTree := BuildHuffmanTree(frequencies)
    prefixTable := BuildPrefixTable(huffmanTree)

    packedBytes, err := PackBytes(string(content), prefixTable)
    if err != nil {
        return err
    }

    return WriteCompressedFile(outputPath, huffmanTree, packedBytes)
}

func DecompressFile(inputPath string, outputPath string) error {
    inputFile, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    huffmanTree, packedBytes, err := ReadCompressedFile(inputFile)
    if err != nil {
        return err
    }

    outputFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    return UnpackBytesAndDecode(packedBytes, outputFile, huffmanTree)
}
