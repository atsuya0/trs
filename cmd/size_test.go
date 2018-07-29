package cmd

import (
	"log"
	"os"
	"testing"
)

func TestSize(t *testing.T) {
	const trashPath = "./test_data/trash_dir"
	const testSize = 15
	testFiles := []string{
		"test1_123456891.txt",
		"test2_123456891.txt",
		"test3_123456891.txt",
	}
	for _, testFile := range testFiles {
		file, err := os.Create(trashPath + "/" + testFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		file.WriteString("12345")
	}

	sum, err := size(trashPath)
	if err != nil {
		log.Fatalln(err)
	}

	if sum != testSize {
		t.Errorf("size() = %d, want %d", sum, testSize)
	}
}
