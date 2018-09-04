package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tayusa/selector"
)

type fileNames []string

func (f fileNames) contains(file string) bool {
	for _, v := range f {
		if file == v {
			return true
		}
	}
	return false
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

// A file of the same name exists in the current directory.
func fileExistsCurrentDir(name string) (bool, error) {
	wd, err := os.Getwd()
	if err != nil {
		return true, err
	}
	filesInCurrentDir, err := getFileNames(wd)
	if err != nil {
		return true, err
	}
	if fileNames(filesInCurrentDir).contains(name) {
		return true, nil
	}
	return false, nil
}

// Remove a character string what given when moving to the trash can.
func removeAffix(org string) string {
	return org[:strings.LastIndex(org, "_")] + filepath.Ext(org)
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

// Select the file to restore.
func selectTarget() (string, string, error) {
	trashPath := getTrashPath()

	date, err := selectFile(trashPath)
	if err != nil {
		return "", "", err
	} else if date == "" {
		return "", "", fmt.Errorf("Cannot get date")
	}

	fileName, err := selectFile(filepath.Join(trashPath, date))
	if err != nil {
		return "", "", err
	} else if fileName == "" {
		return "", "", fmt.Errorf("Cannot get file name")
	}

	filePath := filepath.Join(trashPath, date, fileName)
	if _, err := os.Stat(filePath); err != nil {
		return "", "", err
	}

	return filePath, removeAffix(fileName), nil
}

// Restore selected file or directory.
func restore(_ *cobra.Command, _ []string) error {
	filePath, newFilePath, err := selectTarget()
	if err != nil {
		return err
	}

	if exists, err := fileExistsCurrentDir(newFilePath); err != nil {
		return err
	} else if exists {
		return fmt.Errorf("A file with the same name already exists.")
	}

	if err := os.Rename(filePath, newFilePath); err != nil {
		return err
	}

	return nil
}

func cmdRestore() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "restore",
		Short: "Move files in the trash to the current directory",
		RunE:  restore,
	}

	return cmd
}
