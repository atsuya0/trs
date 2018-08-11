package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "trash",
		Short: "move files to trash",
	}

	cmd.AddCommand(cmdMove())
	cmd.AddCommand(cmdList())
	cmd.AddCommand(cmdDelete())
	cmd.AddCommand(cmdRestore())
	cmd.AddCommand(cmdSize())
	cmd.AddCommand(cmdAutoDelete())

	return cmd
}

func Execute() {
	if err := createTrash(); err != nil {
		log.Fatalln(err)
	}

	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	cmd.Execute()
}
