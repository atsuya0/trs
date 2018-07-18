package main

import (
	"log"
	"os"
	"path"
	"time"
)

func moveToTrashCan(trashCanPath string, files []string) { // ファイルをゴミ箱に移動させる
	prefix := trashCanPath + "/" + time.Now().Format("2006-01-02_15:04:05") + "_"

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			log.Println(err)
			continue
		}

		if err := os.Rename(file, prefix+path.Base(file)); err != nil {
			log.Println(err)
		}
	}
}
