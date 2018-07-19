package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func currentDirNames() ([]string, error) { // カレントディレクトリのファイル・ディレクトリ名の一覧
	var files []string

	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return files, err
	}

	file, err := os.Open(wd)
	defer file.Close()

	if err != nil {
		log.Println(err)
		return files, err
	}

	files, err = file.Readdirnames(0)
	if err != nil {
		log.Println(err)
		return files, err
	}

	return files, err
}

// ファイルの配列に、あるファイルが存在していのかどうか調べる。
func contains(file string, files []string) bool {
	for _, v := range files {
		if file == v {
			return true
		}
	}
	return false
}

// ゴミ箱からファイルを取り出す
func restore(trashCanPath string, trashFiles []string) error {
	files, err := currentDirNames()
	if err != nil {
		return err
	}

	for _, fileName := range trashFiles {
		filePath := trashCanPath + "/" + fileName
		if _, err := os.Stat(filePath); err != nil {
			log.Println(err)
			continue
		}

		newFileName :=
			fileName[:strings.LastIndex(fileName, "_")] +
				filepath.Ext(fileName)

		if contains(newFileName, files) {
			log.Println("同じファイル名のファイルがあります")
			continue
		}
		if err := os.Rename(filePath, newFileName); err != nil {
			log.Println(err)
		}
	}

	return nil
}
