package cmd

import (
	"log"
	"testing"
)

func TestList(t *testing.T) {
	const trashPath = "./test_data/trash_dir"
	testFiles := []string{
		"test1_123456891.txt",
		"test2_123456891.txt",
		"test3_123456891.txt",
	}

	options := Options{}
	options.days = 0
	options.reverse = false
	files, err := list(options, trashPath)
	if err != nil {
		log.Fatalln(err)
	}
	for i, file := range files {
		if testFiles[i] != file[1] {
			t.Errorf("Return value is %s, The desired file name is %s", file[1], testFiles[i])
		}
	}

	options.reverse = true
	files, err = list(options, trashPath)
	if err != nil {
		log.Fatalln(err)
	}
	for i, file := range files {
		testFile := testFiles[len(testFiles)-i-1]
		if testFile != file[1] {
			t.Errorf("Return value is %s, The desired file name is %s", file[1], testFile)
		}
	}
}
