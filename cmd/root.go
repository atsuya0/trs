package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func createTrash(trashPath string) error {
	if _, err := os.Stat(trashPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashPath, 0700); err != nil {
		return err
	}

	return nil
}

func createRootCmd(trashPath string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "trash",
		Short: "move files to trash",
	}

	cmd.AddCommand(createMoveCmd(trashPath))
	cmd.AddCommand(createListCmd(trashPath))
	cmd.AddCommand(createDeleteCmd(trashPath))
	cmd.AddCommand(createRestoreCmd(trashPath))
	cmd.AddCommand(createSizeCmd(trashPath))
	cmd.AddCommand(createAutoDeleteCmd(trashPath))

	return cmd
}

func Execute() {
	trashPath := os.Getenv("TRASH_PATH")
	if trashPath == "" {
		trashPath = os.Getenv("HOME") + "/.Trash"
	}
	if err := createTrash(trashPath); err != nil {
		log.Fatalln(err)
	}

	cmd := createRootCmd(trashPath)
	cmd.SetOutput(os.Stdout)

	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
