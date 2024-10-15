package main

import (
    "bytes"
)

func PackBytes(content string, prefixTable map[rune]string) ([]byte, error) {
    var packedBits bytes.Buffer

    for _, char := range content {
        prefix := prefixTable[char]
        packedBits.WriteString(prefix)
    }

    return bitsToBytes(packedBits.String()), nil
}

func bitsToBytes(bits string) []byte {
    var bytesBuffer bytes.Buffer
    for i := 0; i < len(bits); i += 8 {
        var b byte
        for j := 0; j < 8 && i+j < len(bits); j++ {
            if bits[i+j] == '1' {
                b |= 1 << (7 - j)
            }
        }
        bytesBuffer.WriteByte(b)
    }
    return bytesBuffer.Bytes()
}
