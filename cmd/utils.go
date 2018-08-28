package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getTrashPath() string {
	path := os.Getenv("TRASH_PATH")

	if path != "" {
		return path
	} else {
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
			return ""
		}
		return filepath.Join(user.HomeDir, ".Trash")
	}
}

func createTrash() error {
	trashPath := getTrashPath()

	if _, err := os.Stat(trashPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashPath, 0700); err != nil {
		return err
	}

	return nil
}
