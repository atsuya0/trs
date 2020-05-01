package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type autoDelOption struct {
	period int
}

func autoDel(option *autoDelOption) error {
	path, err := getTrashCanPath()
	if err != nil {
		return err
	}

	dirs, err := ls(path)
	if err != nil {
		return err
	}

	period := time.Now().AddDate(0, 0, -option.period).UnixNano()

	for _, dir := range dirs {
		date, err := time.Parse("2006-01-02", dir)
		if err != nil {
			return err
		}
		if date.UnixNano() < period {
			if err := os.RemoveAll(filepath.Join(path, dir)); err != nil {
				return err
			}
		}
	}

	return nil
}

func autoDeleteCmd() *cobra.Command {
	option := &autoDelOption{}

	var cmd = &cobra.Command{
		Use:   "auto-delete",
		Short: "Delete the files if the date and time that the file moved in the trash can exceed the specified period.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return autoDel(option)
		},
	}

	cmd.Flags().IntVarP(
		&option.period, "period", "p", 30,
		"Delete the files moved in the trash can [days] days ago and later.")

	return cmd
}
