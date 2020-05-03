package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

type filePathPair struct {
	oldPath string
	newDir  string
	newFile string
}

func (p *filePathPair) join() string {
	return filepath.Join(p.newDir, p.newFile)
}

func (p *filePathPair) rename() error {
	return os.Rename(p.oldPath, p.join())
}

func (p *filePathPair) oldFileExists() error {
	if _, err := os.Stat(p.oldPath); os.IsNotExist(err) {
		return &fileNotFoundError{err: err, path: p.oldPath}
	} else if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (p *filePathPair) newFileExists() error {
	fileNames, err := ls(p.newDir)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, file := range fileNames {
		if p.newFile == file {
			return &alreadyFileExistsError{path: p.join()}
		}
	}
	return nil
}

type filePathPairs []filePathPair

func (p *filePathPairs) add(filePathPair filePathPair) error {
	if err := filePathPair.oldFileExists(); err != nil {
		return fmt.Errorf("%w", err)
	}
	if err := filePathPair.newFileExists(); err != nil {
		return fmt.Errorf("%w", err)
	}
	*p = append(*p, filePathPair)

	return nil
}
