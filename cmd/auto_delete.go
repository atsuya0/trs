package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func autoDel(path string) (files []string, err error) {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	now := time.Now()
	oneMonthAgo := now.AddDate(0, -1, 0)

	for _, info := range fileInfo {
		internalStat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			err = fmt.Errorf("fileInfo.Sys(): cast error")
			return
		}
		if (internalStat.Ctim.Nano() - oneMonthAgo.UnixNano()) < 0 {
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
