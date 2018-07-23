package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func autoDel(path string) (files []string, err error) {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	period, err := loadStoragePeriod()
	if err != nil {
		return
	}

	for _, info := range fileInfo {
		internalStat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			err = fmt.Errorf("fileInfo.Sys(): cast error")
			return
		}
		if period > internalStat.Ctim.Nano() {
			files = append(files, path+"/"+info.Name())
		}
	}

	return
}

func createAutoDeleteCmd(trashPath string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auto-delete",
		Short: "delete files that passed one month after moving to the trash can",
		Run: func(cmd *cobra.Command, args []string) {
			if files, err := autoDel(trashPath); err == nil {
				for _, file := range files {
					if err := os.RemoveAll(file); err != nil {
						log.Fatalln(err)
					}
				}
			} else {
				log.Fatalln(err)
			}
		},
	}

	return cmd
}
