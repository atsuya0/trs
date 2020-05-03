package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func generateAffix() string {
	return "_" + time.Now().Format("2006-01-02T15:04:05")
}

func removeExt(fileName string) string {
	return path.Base(fileName[:len(fileName)-len(getExt(fileName))])
}

func insertAffix(fileName string) string {
	return strings.Replace(removeExt(fileName), " ", "", -1) +
		generateAffix() +
		getExt(fileName)
}

func getDestination(filePath string) (string, string, error) {
	fileAbsolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", "", &fileNotFoundError{err: err, path: filePath}
	}
	trashCanPath, err := getTrashCanPath()
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	dirName := filepath.Join(trashCanPath, filepath.Dir(fileAbsolutePath))
	fileName := insertAffix(filepath.Base(fileAbsolutePath))
	if _, err := os.Stat(dirName); err == nil {
		return dirName, fileName, nil
	}
	if err := os.MkdirAll(dirName, 0700); err != nil {
		return "", "", fmt.Errorf("%w", err)
	}

	return dirName, fileName, nil
}

func move(_ *cobra.Command, args []string) error {
	filePathPairs := make(filePathPairs, 0, len(args))
	for _, filePath := range args {
		dirName, fileName, err := getDestination(filePath)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		filePathPair := filePathPair{
			oldPath: filePath,
			newDir:  dirName,
			newFile: fileName,
		}
		if err := filePathPairs.add(filePathPair); err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	for _, filePathPair := range filePathPairs {
		if err := filePathPair.rename(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func moveCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "move",
		Short: "Move the files to the trash can.",
		RunE:  move,
	}

	return cmd
}
