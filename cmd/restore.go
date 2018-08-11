package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type array []string

// ファイルの配列に、あるファイルが存在していのかどうか調べる。
func (a *array) contains(file string) bool {
	for _, v := range *a {
		if file == v {
			return true
		}
	}
	return false
}

// カレントディレクトリのファイル・ディレクトリ名の一覧
func currentDirNames() (array, error) {
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

// ゴミ箱からファイルを取り出す
func restore(_ *cobra.Command, args []string) error {
	trashPath, err := getSrc()
	if err != nil {
		return err
	}

	files, err := currentDirNames()
	if err != nil {
		return err
	}

	for _, fileName := range args {
		filePath := filepath.Join(trashPath, fileName)
		if _, err := os.Stat(filePath); err != nil {
			log.Println(err)
			continue
		}

		newFilePath :=
			fileName[:strings.LastIndex(fileName, "_")] +
				filepath.Ext(fileName)

		if files.contains(newFilePath) {
			log.Println("A file with the same name already exists.")
			continue
		}

		if err := os.Rename(filePath, newFilePath); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func cmdRestore() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "restore",
		Short: "Move files in the trash to the current directory",
		RunE:  restore,
	}

	return cmd
}
