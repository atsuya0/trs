package cmd

import (
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestMove(t *testing.T) {
	const workPath = "./test_data/work_dir"
	const trashPath = "./test_data/trash_dir"

	testFiles := []string{
		workPath + "/test1.txt",
		workPath + "/test2.txt",
		workPath + "/test3.txt",
	}

	for i, setFile := range move(trashPath, testFiles) {
		file := filepath.Base(setFile[1])
		testFile := filepath.Base(testFiles[i])

		prefix := testFile[:len(testFile)-len(filepath.Ext(testFile))]
		if !strings.HasPrefix(file, prefix) {
			t.Errorf("New file name %s, The desired prefix is %s", file, prefix)
		}
		ext := filepath.Ext(testFile)
		if !strings.HasSuffix(file, ext) {
			t.Errorf("New file name %s, The desired extension is %s", file, ext)
		}

		affix :=
			file[strings.Index(file, "_")+1 : len(file)-len(filepath.Ext(file))]
		_, err := strconv.Atoi(affix)
		if err != nil {
			t.Errorf("Affix of new file name is  = %s, The desired affix is integer",
				affix)
		}
	}

}
