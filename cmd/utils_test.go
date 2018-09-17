package cmd

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestGetTrashPath(t *testing.T) {
	home := os.Getenv("HOME")
	patterns := []struct {
		env    string
		output string
	}{
		{env: filepath.Join(home, "TEST"), output: filepath.Join(home, "TEST")},
		{env: "", output: filepath.Join(home, ".Trash")},
	}

	for _, pattern := range patterns {
		if err := os.Setenv("TRASH_PATH", pattern.env); err != nil {
			log.Println(err)
			continue
		}
		if getTrashPath() != pattern.output {
			t.Errorf("%s != %s", getTrashPath(), pattern.output)
		}
	}
}

func TestGetFileNames(t *testing.T) {
	patterns := []struct {
		path  string
		files []string
	}{
		{path: "testdata/utils", files: []string{"test3.txt", "test2", "test1.txt"}},
	}

	for _, pattern := range patterns {
		files, err := getFileNames(pattern.path)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(files) < len(pattern.files) {
			t.Error("The array of execution result is too short.")
		} else if len(files) > len(pattern.files) {
			t.Error("The array of execution result is too long.")
		}
		for i, file := range pattern.files {
			if file != files[i] {
				t.Errorf("%s != %s", file, files[i])
			}
		}
	}
}
