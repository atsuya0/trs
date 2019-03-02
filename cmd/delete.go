package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
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

	date, err := chooseFile(path)
	if err != nil {
		return err
	}
	file, err := chooseFile(filepath.Join(path, date))
	if err != nil {
		return err
	}

	if confirmDel(file) {
		if err := os.RemoveAll(filepath.Join(path, date, file)); err != nil {
			return err
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
