package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"

	"github.com/tayusa/go-choice"
)

func getTrashCanPath() (string, error) {
	if path := os.Getenv("TRASH_CAN_PATH"); path != "" {
		return path, nil
	}
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return filepath.Join(user.HomeDir, ".Trash"), nil
}

// Create a directory as a trash can.
func createTrashCan() error {
	trashCanPath, err := getTrashCanPath()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := os.Stat(trashCanPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashCanPath, 0700); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Get files and directories from the specified path.
func ls(path string) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return []string{}, fmt.Errorf("%w", err)
	}

	defer func() {
		if err = fd.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

	files, err := fd.Readdirnames(0)
	if err != nil {
		return []string{}, fmt.Errorf("%w", err)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, nil
}

func chooseFiles(path string) ([]string, error) {
	files, err := ls(path)
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	fileChooser, err := choice.NewChooser(files)
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
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
