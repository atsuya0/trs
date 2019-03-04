package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "trash",
		Short: "Move the files to the trash can.",
	}

	cmd.AddCommand(moveCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(restoreCmd())
	cmd.AddCommand(sizeCmd())
	cmd.AddCommand(autoDeleteCmd())

	return cmd
}

func Execute() error {
	if err := createTrashCan(); err != nil {
		log.Fatalf("%+v\n", xerrors.Errorf("Cannot create the trash can: %w", err))
	}

	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	return cmd.Execute()
}
