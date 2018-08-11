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

func move(_ *cobra.Command, args []string) error {
	trashPath, err := getSrc()
	if err != nil {
		return err
	}

	prefix := "_" + strconv.FormatInt(time.Now().Unix(), 10)

	for _, fileName := range args {
		if _, err := os.Stat(fileName); err != nil {
			log.Println(err)
			continue
		}

		newFileName := trashPath + "/" +
			path.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))]) +
			prefix +
			filepath.Ext(fileName)

		if err := os.Rename(fileName, newFileName); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func cmdMove() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "move",
		Short: "Move files in the current directory to the trash",
		RunE:  move,
	}

	return cmd
}
