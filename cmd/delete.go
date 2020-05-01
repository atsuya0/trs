package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func confirmDel(file string) bool {
	scanner := bufio.NewScanner(os.Stdin)
	const message = "target: %s\n'yes' or 'no' >>> "

	fmt.Printf(message, file)
	scanner.Scan()
	for scanner.Text() != "yes" && scanner.Text() != "no" {
		fmt.Print("\x1b[2A\x1b[1G\x1b[J")
		fmt.Printf(message, file)
		scanner.Scan()
	}

	if scanner.Text() == "yes" {
		return true
	}
	return false
}

func del(_ *cobra.Command, args []string) error {
	path, err := getTrashCanPath()
	if err != nil {
		return err
	}

	date, fileNames, err := specifyTargets(path)
	if err != nil {
		return xerrors.Errorf("Cannot specify the files to restore: %w", err)
	}

	for _, fileName := range fileNames {
		if confirmDel(fileName) {
			if err := os.RemoveAll(filepath.Join(path, date, fileName)); err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a file in the trash can",
		RunE:  del,
	}

	return cmd
}
