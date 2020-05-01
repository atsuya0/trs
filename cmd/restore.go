package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

type targets struct {
	path      string
	fileNames []string
}

func (t targets) createPairs() ([]filePathPair, error) {
	var filePathPairs []filePathPair
	for _, v := range t.fileNames {
		pair := filePathPair{old: filepath.Join(t.path, v), new: removeAffix(v)}
		if err := pair.oldFileExists(); err != nil {
			return make([]filePathPair, 0), err
		}
		if err := pair.newFileExists(); err != nil {
			return make([]filePathPair, 0), err
		}
		filePathPairs = append(filePathPairs, pair)
	}
	return filePathPairs, nil
}

type filePathPair struct {
	old string
	new string
}

func (f filePathPair) oldFileExists() error {
	if _, err := os.Stat(f.old); err != nil {
		return xerrors.Errorf("The specified file does not exist: %w", err)
	}
	return nil
}

func (f filePathPair) newFileExists() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	filesInCurrentDir, err := ls(wd)
	if err != nil {
		return xerrors.Errorf("Cannot ls: %w", err)
	}
	for _, file := range filesInCurrentDir {
		if f.new == file {
			return xerrors.New("A file with the same name already exists")
		}
	}
	return nil
}

// Remove a character string what given when moving to the trash can.
func removeAffix(org string) string {
	return org[:strings.LastIndex(org, "_")] + getExt(org)
}

// Specify the files to restore.
func specifyTargets(trashCanPath string) (string, []string, error) {
	for {
		dates, err := chooseFiles(trashCanPath)
		date := dates[0]
		if err != nil {
			return "", make([]string, 0), xerrors.Errorf("Cannot choose the date: %w", err)
		} else if date == "" {
			return "", make([]string, 0), nil
		}

		fileNames, err := chooseFiles(filepath.Join(trashCanPath, date))
		if err != nil {
			return "", make([]string, 0), xerrors.Errorf("Cannot choose the file: %w", err)
		} else if len(fileNames) != 0 {
			return date, fileNames, nil
		}
	}
}

func getTargets() ([]filePathPair, error) {
	path, err := getTrashCanPath()
	if err != nil {
		return make([]filePathPair, 0), xerrors.Errorf("Cannot get the path of the trash can: %w", err)
	}

	date, fileNames, err := specifyTargets(path)
	if err != nil {
		return make([]filePathPair, 0), xerrors.Errorf("Cannot specify the files to restore: %w", err)
	}
	targets := targets{path: filepath.Join(path, date), fileNames: fileNames}
	filePathPairs, err := targets.createPairs()
	return filePathPairs, err
}

func restore(_ *cobra.Command, _ []string) error {
	filePathPairs, err := getTargets()
	if err != nil {
		return err
	}

	for _, v := range filePathPairs {
		if err := os.Rename(v.old, v.new); err != nil {
			return err
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
