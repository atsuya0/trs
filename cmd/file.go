package cmd

import (
	"fmt"
	"os"
)

type file struct {
	info os.FileInfo
	path string
}

func (f *file) removeEmptyDir() error {
	if !f.info.IsDir() {
		return nil
	}
	if childFiles, err := ls(f.path); err != nil {
		return fmt.Errorf("%w", err)
	} else if len(childFiles) == 0 {
		if err := os.RemoveAll(f.path); err != nil {
			return err
		}
	}
	return nil
}

type Files []file

func (f Files) Len() int {
	return len(f)
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
