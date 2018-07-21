package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"time"
)

// ゴミ箱に入っている、指定した一つのファイルを削除する。
func del(file string) bool {
	fmt.Printf("target: %s\n", file)
	fmt.Println("'yes' or 'no'")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	for scanner.Text() != "yes" && scanner.Text() != "no" {
		fmt.Println("'yes' or 'no'")
		scanner.Scan()
	}

	if scanner.Text() == "yes" {
		return true
	} else {
		return false
	}
}

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