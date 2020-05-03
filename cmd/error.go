package cmd

import "fmt"

type fileNotFoundError struct {
	err  error
	path string
}

func (e *fileNotFoundError) Error() string {
	return fmt.Errorf("%s not found: %w", e.path, e.err).Error()
}

type alreadyFileExistsError struct {
	path string
}

func (e *alreadyFileExistsError) Error() string {
	return fmt.Sprintf("%s with the same name already exists.", e.path)
}

type dirNotFoundError struct {
	path string
}

func (e *dirNotFoundError) Error() string {
	return fmt.Sprintf("%s not found.", e.path)
}
