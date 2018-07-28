package cmd

import (
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestRestore(t *testing.T) {
	const trashPath = "./test_data/trash_dir"
	testFiles := []string{
		"test1_123456891.txt",
		"test2_123456891.txt",
		"test3_123456891.txt",
	}

	if setFiles, err := restore(trashPath, testFiles); err == nil {
		for i, setFile := range setFiles {
			desiredFileName := testFiles[i][:strings.Index(testFiles[i], "_")] + filepath.Ext(testFiles[i])
			if setFile[1] != desiredFileName {
				t.Errorf("new file name %s, The desired file name is %s", setFile[1], desiredFileName)
			}
		}
	} else {
		log.Fatalln(err)
	}
}
