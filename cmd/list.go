package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type Options struct {
	days    int
	reverse bool
}

type ExFileInfo struct {
	info  os.FileInfo
	color string
}

type Files []os.FileInfo

func (f Files) Len() int {
	return len(f)
}

func (f Files) Less(i, j int) bool {
	return f[i].Sys().(*syscall.Stat_t).Ctim.Nano() <
		f[j].Sys().(*syscall.Stat_t).Ctim.Nano()
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// ゴミ箱の中のファイル一覧を表示
func list(options Options, path string) (files []ExFileInfo, err error) {
	files = make([]ExFileInfo, 0, len(files))

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	if options.reverse {
		sort.Sort(sort.Reverse(Files(fileInfo)))
	} else {
		sort.Sort(Files(fileInfo))
	}

	const executable os.FileMode = 0111
	const green = "\x1b[32m\x1b[1m%s\x1b[39m\x1b[0m\n"
	const blue = "\x1b[34m\x1b[1m%s\x1b[39m\x1b[0m\n"
	const cyan = "\x1b[36m\x1b[1m%s\x1b[39m\x1b[0m\n"
	const white = "\x1b[37m\x1b[0m%s\x1b[39m\x1b[0m\n"

	now := time.Now()
	daysAgo := now.AddDate(0, 0, -options.days)

	for _, info := range fileInfo {
		internalStat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			err = fmt.Errorf("fileInfo.Sys(): cast error")
			return
		}

		if options.days != 0 && internalStat.Ctim.Nano() < daysAgo.UnixNano() {
			continue
		}
		if info.IsDir() {
			files = append(files, ExFileInfo{info: info, color: blue})
		} else if info.Mode()&os.ModeSymlink != 0 {
			files = append(files, ExFileInfo{info: info, color: cyan})
		} else if info.Mode()&executable != 0 {
			files = append(files, ExFileInfo{info: info, color: green})
		} else {
			files = append(files, ExFileInfo{info: info, color: white})
		}
	}

	return
}

func createListCmd(trashPath string) *cobra.Command {
	options := &Options{}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "The list of the trash",
		Run: func(cmd *cobra.Command, args []string) {
			files, err := list(*options, trashPath)
			if err != nil {
				log.Fatalln(err)
			}
			for _, file := range files {
				fmt.Printf(file.color, file.info.Name())
			}
		},
	}
	cmd.Flags().IntVarP(&options.days, "days", "d", 0, "How many days ago")
	cmd.Flags().BoolVarP(&options.reverse, "reverse", "r", false, "display in reverse order")

	return cmd
}
