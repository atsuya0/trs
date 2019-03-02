package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"

	choice "github.com/tayusa/go-choice"
)

func getTrashCanPath() (string, error) {
	if path := os.Getenv("TRASH_CAN_PATH"); path != "" {
		return path, nil
	}
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(user.HomeDir, ".Trash"), nil
}

// Create a directory as a trash can.
func createTrashCan() error {
	trashCanPath, err := getTrashCanPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(trashCanPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashCanPath, 0700); err != nil {
		return err
	}

	return nil
}

// Fetch files and directories from the specified path.
func getFileNames(path string) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}

	defer func() {
		if err = fd.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	files, err := fd.Readdirnames(0)
	if err != nil {
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

// If it is a hidden file with no extension, it returns an empty string.
func getExt(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) == len(fileName) {
		return ""
	} else {
		return ext
	}
}
