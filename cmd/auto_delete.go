package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func loadStoragePeriod() (int64, error) {
	const defaultPeriod = -30

	strPeriod := os.Getenv("STORAGE_PERIOD_OF_THE_TRASH")
	if strPeriod == "" {
		return time.Now().AddDate(0, 0, defaultPeriod).UnixNano(), nil
	}
	period, err := strconv.Atoi(strPeriod)
	if err != nil {
		return 0, err
	}
	return time.Now().AddDate(0, 0, -period).UnixNano(), nil
}

func autoDel(_ *cobra.Command, _ []string) error {
	path := getTrashPath()

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	period, err := loadStoragePeriod()
	if err != nil {
		return err
	}

	for _, info := range fileInfo {
		internalStat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("fileInfo.Sys(): cast error")
		}
		if period > internalStat.Ctim.Nano() {
			if err := os.RemoveAll(filepath.Join(path, info.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func cmdAutoDelete() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auto-delete",
		Short: "Delete files that passed one month after moving to the trash can",
		RunE:  autoDel,
	}

	return cmd
}
