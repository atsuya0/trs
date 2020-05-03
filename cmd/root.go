package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "trash",
		Short: "Move the files to the trash can.",
	}

	cmd.AddCommand(moveCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(removeCmd())
	cmd.AddCommand(restoreCmd())
	cmd.AddCommand(sizeCmd())
	cmd.AddCommand(autoRemoveCmd())

	return cmd
}

func Execute() error {
	if err := createTrashCan(); err != nil {
		return fmt.Errorf("Cannot create the trash can: %w", err)
	}

	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	return cmd.Execute()
}
