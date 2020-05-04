package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tayusa/go-choice"
)

type restoreOption struct {
	all bool
}

type targets struct {
	path      string
	fileNames []string
}

func (t targets) createPairs() (filePathPairs, error) {
	wd, err := os.Getwd()
	if err != nil {
		return make(filePathPairs, 0), fmt.Errorf("%w", err)
	}
	var pairs filePathPairs
	for _, v := range t.fileNames {
		filePathPair := filePathPair{
			oldPath: filepath.Join(t.path, v),
			newFile: removeAffix(v),
			newDir:  wd,
		}
		if err := pairs.add(filePathPair); err != nil {
			return make(filePathPairs, 0), fmt.Errorf("%w", err)
		}
	}
	return pairs, nil
}

// Remove a character string what given when moving to the trash can.
func removeAffix(fileName string) string {
	if index := strings.LastIndex(fileName, "_"); index >= 0 {
		return fileName[:index] + getExt(fileName)
	}
	return fileName
}

func getFilePathPairsInCorrespondingPath() (filePathPairs, error) {
	correspondingPath, fileNames, err := chooseFilesInCorrespondingPath()
	if err != nil {
		return make(filePathPairs, 0), fmt.Errorf("%w", err)
	}
	targets := targets{path: correspondingPath, fileNames: fileNames}
	filePathPairs, err := targets.createPairs()
	return filePathPairs, err
}

func getFilePaths() ([]string, error) {
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
	filePaths, err := getFilePaths()
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}

	fileChooser, err := choice.NewChooser(filePaths)
	if err != nil {
		return make([]string, 0), fmt.Errorf("%w", err)
	}
	return fileChooser.Run(), nil
}

func getFilePathPairs() (filePathPairs, error) {
	filePaths, err := chooseFilePaths()
	if err != nil {
		return make(filePathPairs, 0), fmt.Errorf("%w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return make(filePathPairs, 0), fmt.Errorf("%w", err)
	}

	pairs := make(filePathPairs, 0, len(filePaths))
	for _, filePath := range filePaths {
		filePathPair := filePathPair{
			oldPath: filePath,
			newDir:  wd,
			newFile: removeAffix(filepath.Base(filePath)),
		}
		if err := pairs.add(filePathPair); err != nil {
			return make(filePathPairs, 0), fmt.Errorf("%w", err)
		}
	}
	return pairs, err
}

func restore(option *restoreOption) error {
	var filePathPairs filePathPairs
	var err error
	if option.all {
		filePathPairs, err = getFilePathPairs()
	} else {
		filePathPairs, err = getFilePathPairsInCorrespondingPath()
	}
	if errors.Is(err, &dirNotFoundError{}) {
		fmt.Println("Never used the move command in this path.")
		return nil
	} else if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, pair := range filePathPairs {
		if err := pair.rename(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func restoreCmd() *cobra.Command {
	option := &restoreOption{}

	var cmd = &cobra.Command{
		Use:   "restore",
		Short: "Move the files in the trash can to the current directory.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return restore(option)
		},
	}

	cmd.Flags().BoolVarP(
		&option.all, "all", "a", false,
		"Target all the files.")

	return cmd
}
