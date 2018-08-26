package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func removeExt(fileName string) string {
	return path.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])
}

func addAffix(fileName string, affix string, trashPath string) string {
	return trashPath + "/" +
		strings.Replace(removeExt(fileName), " ", "", -1) +
		affix +
		filepath.Ext(fileName)
}

func move(_ *cobra.Command, args []string) error {
	trashPath, err := getSrc()
	if err != nil {
		return err
	}

	affix := "_" + strconv.FormatInt(time.Now().Unix(), 10)

	for _, fileName := range args {
		if _, err := os.Stat(fileName); err != nil {
			log.Println(err)
			continue
		}

		if err := os.Rename(fileName, addAffix(fileName, affix, trashPath)); err != nil {
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
