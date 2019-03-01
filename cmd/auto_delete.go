package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func loadStoragePeriod() (int64, error) {
	const defaultPeriod = 30

	strPeriod := os.Getenv("TRASH_CAN_PERIOD")
	if strPeriod == "" {
		return time.Now().AddDate(0, 0, -defaultPeriod).UnixNano(), nil
	}
	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		return 0, err
	}
	return time.Now().AddDate(0, 0, -period).UnixNano(), nil
}

func autoDel(_ *cobra.Command, _ []string) error {
	path := getTrashCanPath()

	dirs, err := getFileNames(path)
	if err != nil {
		return err
	}

	period, err := loadStoragePeriod()
	if err != nil {
		return err
	}

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
	var cmd = &cobra.Command{
		Use:   "auto-delete",
		Short: "Delete files that passed one month after moving to the trash can",
		RunE:  autoDel,
	}

	return cmd
}
