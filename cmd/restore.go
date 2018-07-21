package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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
func restore(trashPath string, trashFiles []string) ([][]string, error) {
	setFiles := make([][]string, 0)

	files, err := currentDirNames()
	if err != nil {
		return setFiles, err
	}

	for _, fileName := range trashFiles {
		filePath := trashPath + "/" + fileName
		if _, err := os.Stat(filePath); err != nil {
			log.Println(err)
			continue
		}

		newFilePath :=
			fileName[:strings.LastIndex(fileName, "_")] +
				filepath.Ext(fileName)

		if contains(newFilePath, files) {
			log.Println("A file with the same name already exists.")
			continue
		}
		setFiles = append(setFiles, []string{filePath, newFilePath})
	}

	return setFiles, err
}

func createRestoreCmd(trashPath string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "restore",
		Short: "move files in the trash to the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			if setFiles, err := restore(trashPath, args); err == nil {
				for _, setFile := range setFiles {
					if err := os.Rename(setFile[0], setFile[1]); err != nil {
						log.Println(err)
					}
				}
			} else {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
