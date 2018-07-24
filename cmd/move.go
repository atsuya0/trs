package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func move(trashPath string, fileNames []string) [][]string {
	prefix := "_" + strconv.FormatInt(time.Now().Unix(), 10)
	setFiles := make([][]string, 0, len(fileNames))

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

func createMoveCmd(trashPath string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "move",
		Short: "move files in the current directory to the trash",
		Run: func(cmd *cobra.Command, args []string) {
			for _, setFile := range move(trashPath, args) {
				if err := os.Rename(setFile[0], setFile[1]); err != nil {
					log.Println(err)
				}
			}
		},
	}

	return cmd
}
