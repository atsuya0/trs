package main

import (
	"os"
	"path/filepath"
)

// ゴミ箱のsizeを返す
func size(root string) (int64, error) {
	var sum int64 = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			sum += info.Size()
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return sum, nil
}
