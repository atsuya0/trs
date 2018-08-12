package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
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
func list(options *Options, out io.Writer) error {
	path, err := getSrc()
	if err != nil {
		return err
	}

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
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

	daysAgo := time.Now().AddDate(0, 0, -options.days)

	for _, info := range fileInfo {
		internalStat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("fileInfo.Sys(): cast error")
		}

		if options.days != 0 && internalStat.Ctim.Nano() < daysAgo.UnixNano() {
			continue
		}
		if info.IsDir() {
			fmt.Fprintf(out, blue, info.Name())
		} else if info.Mode()&os.ModeSymlink != 0 {
			fmt.Fprintf(out, cyan, info.Name())
		} else if info.Mode()&executable != 0 {
			fmt.Fprintf(out, green, info.Name())
		} else {
			fmt.Fprintf(out, white, info.Name())
		}
	}

	return nil
}

func cmdList() *cobra.Command {
	options := &Options{}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "The list of the trash",
		RunE: func(_ *cobra.Command, _ []string) error {
			return list(options, os.Stdout)
		},
	}
	cmd.Flags().IntVarP(&options.days, "days", "d", 0, "How many days ago")
	cmd.Flags().BoolVarP(&options.reverse, "reverse", "r", false, "display in reverse order")

	return cmd
}
