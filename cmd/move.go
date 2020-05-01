package cmd

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func generateAffix() string {
	return "_" + time.Now().Format("15:04:05")
}

func generateDestination(path string) (string, bool) {
	destination := filepath.Join(path, time.Now().Format("2006-01-02"))
	if _, err := os.Stat(destination); err == nil {
		return destination, true
	}
	return "", false
}

func removeExt(fileName string) string {
	return path.Base(fileName[:len(fileName)-len(getExt(fileName))])
}

func insertAffix(fileName string, affix string, destination string) string {
	return filepath.Join(destination,
		strings.Replace(removeExt(fileName), " ", "", -1)+
			affix+
			getExt(fileName))
}

func getDestination(trashCanPath string) (string, error) {
	destination, isGenerated := generateDestination(trashCanPath)
	if !isGenerated {
		return destination, nil
	}

	if err := os.Mkdir(destination, 0700); err != nil {
		return "", err
	}

	return destination, nil
}

func move(_ *cobra.Command, args []string) error {
	path, err := getTrashCanPath()
	if err != nil {
		return err
	}

	destination, err := getDestination(path)
	if err != nil {
		return err
	}

	affix := generateAffix()
	for _, fileName := range args {
		if _, err := os.Stat(fileName); err != nil {
			log.Printf("%+v\n", err)
			continue
		}

		if err := os.Rename(fileName, insertAffix(fileName, affix, destination)); err != nil {
			log.Printf("%+v\n", err)
		}
	}

	return nil
}

func moveCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "move",
		Short: "Move the files to the trash can",
		RunE:  move,
	}

	return cmd
}
