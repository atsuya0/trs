package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func moveToTrash(trashPath string, fileNames []string) [][]string { // ファイルをゴミ箱に移動させる
	prefix := "_" + strconv.FormatInt(time.Now().Unix(), 10)
	setFiles := make([][]string, 0)

	for _, fileName := range fileNames {
		if _, err := os.Stat(fileName); err != nil {
			log.Println(err)
			continue
		}

		newFileName := trashPath + "/" +
			path.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))]) +
			prefix +
			filepath.Ext(fileName)

		setFiles = append(setFiles, []string{fileName, newFileName})
	}

	return setFiles
}
