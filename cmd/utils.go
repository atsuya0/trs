package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"

	"github.com/tayusa/go-choice"
)

func getTrashPath() string {
	if path := os.Getenv("TRASH_PATH"); path != "" {
		return path
	}
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	return filepath.Join(user.HomeDir, ".Trash")
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
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, err
}

// Choose one from files and directories.
func chooseFile(path string) (string, error) {
	files, err := getFileNames(path)
	if err != nil {
		return "", err
	}
	fileChooser := choice.NewChooser(files)
	return fileChooser.Run(), nil
}
