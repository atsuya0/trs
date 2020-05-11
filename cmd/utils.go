package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/tayusa/go-chooser"
)

const (
	logFileName = "trs.log"
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

// Filter files that have an affix.
func filterFilesHaveAffix(fileNames []string) ([]string, error) {
	r, err := regexp.Compile("_[0-9]*-[0-9]*-[0-9]*T[0-9]*:[0-9]*:[0-9]*.*$")
	if err != nil {
		return make([]string, 0), err
	}
	var filteredFiles []string
	for _, fileName := range fileNames {
		if r.MatchString(fileName) {
			filteredFiles = append(filteredFiles, fileName)
		}
	}
	return filteredFiles, nil
}

func getCorrespondingPath() (string, error) {
	root, err := getTrashCanPath()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	path := filepath.Join(root, wd)
	if _, err := os.Stat(path); err != nil {
		return "", &dirNotFoundError{path: path}
	}

	return path, nil
}

func chooseFilesInCorrespondingPath() (string, []string, error) {
	correspondingPath, err := getCorrespondingPath()
	if err != nil {
		return "", make([]string, 0), fmt.Errorf("%w", err)
	}
	files, err := ls(correspondingPath)
	if err != nil {
		return "", make([]string, 0), fmt.Errorf("%w", err)
	}
	filteredFiles, err := filterFilesHaveAffix(files)
	if err != nil {
		return "", make([]string, 0), fmt.Errorf("%w", err)
	}
	fileChooser, err := chooser.NewChooser(filteredFiles)
	if err != nil {
		return "", make([]string, 0), fmt.Errorf("%w", err)
	}
	return correspondingPath, fileChooser.Run(), nil
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

func getFilesInTrash() (Files, Files, error) {
	root, err := getTrashCanPath()
	if err != nil {
		return make(Files, 0), make(Files, 0), fmt.Errorf("%w", err)
	}

	var files Files
	var dirs Files
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if info.IsDir() {
			dirs = append(dirs, file{info: info, path: path})
		} else {
			files = append(files, file{info: info, path: path})
		}

		return nil
	}); err != nil {
		return make(Files, 0), make(Files, 0), fmt.Errorf("%w", err)
	}
	return files, dirs, nil
}

func getFilePathsInTrash() ([]string, error) {
	root, err := getTrashCanPath()
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}

	var paths []string
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}

		return nil
	}); err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	return paths, nil
}

func chooseFilePaths() ([]string, error) {
	filePaths, err := getFilePathsInTrash()
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}

	fileChooser, err := chooser.NewChooser(filePaths)
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	return fileChooser.Run(), nil
}

func writeLog(msg string) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	f, err := os.Create(filepath.Join(dir, logFileName))
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if _, err := f.WriteString(msg); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
