package main

import (
	"os"
	"testing"
)

func TestCalculateCharFrequency(t *testing.T) {
	data := "hello world"
	expected := map[rune]int{
		'h': 1,
		'e': 1,
		'l': 3,
		'o': 2,
		'w': 1,
		'r': 1,
		'd': 1,
		' ': 1,
	}

	// Write data to a temporary file for testing
	file, _ := os.CreateTemp("", "testfile")
	defer os.Remove(file.Name())
	file.WriteString(data)
	file.Close()

	// Calculate character frequencies
	freq, err := calculateCharFrequency(file.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the results
	for k, v := range expected {
		if freq[k] != v {
			t.Errorf("Expected %q: %d, but got %d", k, v, freq[k])
		}
	}
}
