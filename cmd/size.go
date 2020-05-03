package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	units = [6]string{"B", "kB", "MB", "GB", "TB", "PB"}
)

func getSize(trashCanPath string) (int64, error) {
	var size int64 = 0

	err := filepath.Walk(trashCanPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return size, nil
}

func convertNumbersToSymbols(size float64, cnt int) string {
	result := math.Pow(1024, float64(cnt))
	if size < result*1024 {
		if cnt >= len(units) {
			return fmt.Sprintf("%0.1f %s", size, units[0])
		}
		return fmt.Sprintf("%0.1f %s", size/result, units[cnt])
	}
	return convertNumbersToSymbols(size, cnt+1)
}

func sizeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "size",
		Short: "The size of the trash can directory",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, err := getTrashCanPath()
			if err != nil {
				return err
			}
			size, err := getSize(path)
			if err != nil {
				return fmt.Errorf("Don't get the the allocated size of trash can: %w", err)
			}

			fmt.Println(convertNumbersToSymbols(float64(size), 0))

			return nil
		},
	}

	return cmd
}
