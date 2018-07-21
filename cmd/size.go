package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func size(trashPath string) (int64, error) {
	var sum int64 = 0

	err := filepath.Walk(trashPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			sum += info.Size()
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func createSizeCmd(trashPath string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "size",
		Short: "The size of the trash directory",
		Run: func(cmd *cobra.Command, args []string) {
			if trashSize, err := size(trashPath); err == nil {
				fmt.Printf("%d MB", trashSize/(1024*1024))
			} else {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
