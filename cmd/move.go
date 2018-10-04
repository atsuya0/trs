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
	return path.Base(fileName[:len(fileName)-len(getExt(fileName))])
}

func insertAffix(fileName string, affix string, destination string) string {
	return filepath.Join(destination,
		strings.Replace(removeExt(fileName), " ", "", -1)+
			affix+
			getExt(fileName))
}

func getDestination(trashPath string) (string, error) {
	destination := filepath.Join(trashPath, time.Now().Format("2006-01-02"))
	if _, err := os.Stat(destination); err == nil {
		return destination, nil
	}

	if err := os.Mkdir(destination, 0700); err != nil {
		return "", err
	}

	return destination, nil
}

func move(_ *cobra.Command, args []string) error {
	destination, err := getDestination(getTrashPath())
	if err != nil {
		return err
	}

	affix := "_" + strconv.FormatInt(time.Now().Unix(), 10)

	for _, fileName := range args {
		if _, err := os.Stat(fileName); err != nil {
			log.Println(err)
			continue
		}

		if err := os.Rename(fileName, insertAffix(fileName, affix, destination)); err != nil {
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
