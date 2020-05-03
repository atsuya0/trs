package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type autoRemoveOption struct {
	period int
}

func getFilesAndDirs() (Files, error) {
	root, err := getTrashCanPath()
	if err != nil {
		return make(Files, 0), fmt.Errorf("%w", err)
	}

	var files Files
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		files = append(files, file{info: info, path: path})

		return nil
	}); err != nil {
		return make(Files, 0), fmt.Errorf("%w", err)
	}
	return files, nil
}

func autoRemove(option *autoRemoveOption) error {
	files, err := getFilesAndDirs()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	days := time.Now().AddDate(0, 0, -option.period).UnixNano()

	for _, file := range files {
		if err := file.removeEmptyDir(); err != nil {
			return fmt.Errorf("%w", err)
		}
		if bool, err := file.withinPeriod(days); err != nil {
			return fmt.Errorf("%w", err)
		} else if bool || file.info.IsDir() {
			continue
		}
		if err := os.RemoveAll(file.path); err != nil {
			return err
		}
	}

	return nil
}

func autoRemoveCmd() *cobra.Command {
	option := &autoRemoveOption{}

	var cmd = &cobra.Command{
		Use:   "auto-remove",
		Short: "Remove the files if the date and time that the file moved in the trash can exceed the specified period.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return autoRemove(option)
		},
	}

	cmd.Flags().IntVarP(
		&option.period, "period", "p", 30,
		"Remove the files moved in the trash can [days] days ago and later.")

	return cmd
}