package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func moveToTrashCan(trashCanPath string, fileNames []string) { // ファイルをゴミ箱に移動させる
	prefix := "_" + strconv.FormatInt(time.Now().Unix(), 10)

	for _, fileName := range fileNames {
		if _, err := os.Stat(fileName); err != nil {
			log.Println(err)
			continue
		}
		newFileName := trashCanPath + "/" +
			path.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))]) +
			prefix +
			filepath.Ext(fileName)

		if err := os.Rename(fileName, newFileName); err != nil {
			log.Println(err)
		}
	}
}
