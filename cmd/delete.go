package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func del(_ *cobra.Command, args []string) error {
	fmt.Printf("target: %s\n", args[0])
	fmt.Println("'yes' or 'no'")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	for scanner.Text() != "yes" && scanner.Text() != "no" {
		fmt.Println("'yes' or 'no'")
		scanner.Scan()
	}

	if scanner.Text() == "yes" {
		path, err := getSrc()
		if err != nil {
			return err
		}
		if err := os.RemoveAll(filepath.Join(path, args[0])); err != nil {
			return err
		}
	}

	return nil
}

func cmdDelete() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a file in the trash",
		Args:  cobra.MinimumNArgs(1),
		RunE:  del,
	}

	return cmd
}
