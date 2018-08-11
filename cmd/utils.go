package cmd

import (
	"os"
	"os/user"
	"path/filepath"
)

func getSrc() (string, error) {
	path := os.Getenv("TRASH_PATH")

	if path != "" {
		return path, nil
	} else {
		user, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(user.HomeDir, ".Trash"), nil
	}
}

func createTrash() error {
	trashPath, err := getSrc()
	if err != nil {
		return err
	}

	if _, err := os.Stat(trashPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashPath, 0700); err != nil {
		return err
	}

	return nil
}
