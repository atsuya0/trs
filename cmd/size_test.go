package cmd

import (
	"log"
	"testing"
)

func TestSize(t *testing.T) {
	const trashPath = "./test_data"
	const testSize = 15

	sum, err := size(trashPath)
	if err != nil {
		log.Fatalln(err)
	}

	if sum != testSize {
		t.Errorf("size() = %d, want %d", sum, testSize)
	}
}
