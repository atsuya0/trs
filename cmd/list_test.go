package cmd

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	const trashPath = "./test_data/trash_dir"
	testFiles := []string{
		"test1_123456891.txt",
		"test2_123456891.txt",
		"test3_123456891.txt",
	}
	for _, testFile := range testFiles {
		time.Sleep(1 * time.Second)
		if file, err := os.Create(trashPath + "/" + testFile); err != nil {
			log.Fatalln(err)
		} else {
			file.Close()
		}
	}

	options := Options{}
	options.days = 0
	options.reverse = false
	files, err := list(options, trashPath)
	if err != nil {
		log.Fatalln(err)
	}
	for i, file := range files {
		if testFiles[i] != file.info.Name() {
			t.Errorf("Return value is %s, The desired file name is %s", file.info.Name(), testFiles[i])
		}
	}

	options.reverse = true
	files, err = list(options, trashPath)
	if err != nil {
		log.Fatalln(err)
	}
	for i, file := range files {
		testFile := testFiles[len(testFiles)-i-1]
		if testFile != file.info.Name() {
			t.Errorf("Return value is %s, The desired file name is %s", file.info.Name(), testFile)
		}
	}
}
