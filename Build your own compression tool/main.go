package main

import (
    "flag"
    "fmt"
)

func main() {
    compressPtr := flag.Bool("compress", false, "Compress the input file")
    decompressPtr := flag.Bool("decompress", false, "Decompress the input file")
    inputFilePtr := flag.String("input", "", "Input file path")
    outputFilePtr := flag.String("output", "", "Output file path")

    flag.Parse()

    if *compressPtr {
        if err := CompressFile(*inputFilePtr, *outputFilePtr); err != nil {
            fmt.Println("Error compressing file:", err)
            return
        }
        fmt.Println("File compressed successfully!")
    } else if *decompressPtr {
        if err := DecompressFile(*inputFilePtr, *outputFilePtr); err != nil {
            fmt.Println("Error decompressing file:", err)
            return
        }
        fmt.Println("File decompressed successfully!")
    } else {
        fmt.Println("Please specify -compress or -decompress")
    }
}
