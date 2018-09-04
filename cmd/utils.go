package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/tayusa/selector"
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

// Create a directory as a trash can.
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

// Fetch files and directories from the specified path.
func getFileNames(path string) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	defer func() {
		if err = fd.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	files, err := fd.Readdirnames(0)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	return files, err
}

// Select one from files and directories.
func selectFile(path string) (string, error) {
	files, err := getFileNames(path)
	if err != nil {
		return "", err
	}
	fileSelector := selector.NewSelector(files)
	return fileSelector.Run(), nil
}
