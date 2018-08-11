package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func size(_ *cobra.Command, _ []string) error {
	var sum int64 = 0

	trashPath, err := getSrc()
	if err != nil {
		return err
	}

	err = filepath.Walk(trashPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			sum += info.Size()
		}

		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("%d MB", sum/(1024*1024))
	return nil
}

func cmdSize() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "size",
		Short: "The size of the trash directory",
		RunE:  size,
	}

	return cmd
}
