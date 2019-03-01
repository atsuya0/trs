package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "trash",
		Short: "move files to trash can",
	}

	cmd.AddCommand(moveCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(restoreCmd())
	cmd.AddCommand(sizeCmd())
	cmd.AddCommand(autoDeleteCmd())

	return cmd
}

func Execute() {
	if err := createTrashCan(); err != nil {
		log.Fatalln(err)
	}

	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	cmd.Execute()
}
