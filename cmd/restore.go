package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type targets struct {
	path      string
	fileNames []string
}

func (t targets) createPairsToRestore() ([]filePathPair, error) {
	wd, err := os.Getwd()
	if err != nil {
		return make([]filePathPair, 0), fmt.Errorf("%w", err)
	}
	var filePathPairs []filePathPair
	for _, v := range t.fileNames {
		pair := filePathPair{
			oldPath: filepath.Join(t.path, v),
			newFile: removeAffix(v),
			newDir:  wd,
		}
		if err := pair.oldFileExists(); err != nil {
			return make([]filePathPair, 0), fmt.Errorf("%w", err)
		}
		if err := pair.newFileExists(); err != nil {
			return make([]filePathPair, 0), fmt.Errorf("%w", err)
		}
		filePathPairs = append(filePathPairs, pair)
	}
	return filePathPairs, nil
}

// Remove a character string what given when moving to the trash can.
func removeAffix(org string) string {
	return org[:strings.LastIndex(org, "_")] + getExt(org)
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

func getPairs() ([]filePathPair, error) {
	correspondingPath, err := getCorrespondingPath()
	if err != nil {
		return make([]filePathPair, 0), fmt.Errorf("%w", err)
	}

	fileNames, err := chooseFiles(correspondingPath)
	if err != nil {
		return make([]filePathPair, 0), fmt.Errorf("%w", err)
	}
	targets := targets{path: correspondingPath, fileNames: fileNames}
	filePathPairs, err := targets.createPairsToRestore()
	return filePathPairs, err
}

func restore(_ *cobra.Command, _ []string) error {
	filePathPairs, err := getPairs()
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
	var cmd = &cobra.Command{
		Use:   "restore",
		Short: "Move the files in the trash can to the current directory",
		RunE:  restore,
	}

	return cmd
}
