package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const (
	executable os.FileMode = 0111
	header                 = "\x1b[4;30;47m%s\x1b[0m\n"
	green                  = "\x1b[32m\x1b[1m%s\x1b[39m\x1b[0m\n"
	blue                   = "\x1b[34m\x1b[1m%s\x1b[39m\x1b[0m\n"
	cyan                   = "\x1b[36m\x1b[1m%s\x1b[39m\x1b[0m\n"
	white                  = "\x1b[37m\x1b[0m%s\x1b[39m\x1b[0m\n"
)

type Options struct {
	days    int
	reverse bool
}

type Dirs []os.FileInfo

func (d Dirs) Len() int {
	return len(d)
}

func (d Dirs) Less(i, j int) bool {
	return d[i].Sys().(*syscall.Stat_t).Ctim.Nano() <
		d[j].Sys().(*syscall.Stat_t).Ctim.Nano()
}

func (d Dirs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func printFiles(out io.Writer, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, header, filepath.Base(path))

	for _, file := range files {
		if file.IsDir() {
			fmt.Fprintf(out, blue, file.Name())
		} else if file.Mode()&os.ModeSymlink != 0 {
			fmt.Fprintf(out, cyan, file.Name())
		} else if file.Mode()&executable != 0 {
			fmt.Fprintf(out, green, file.Name())
		} else {
			fmt.Fprintf(out, white, file.Name())
		}
	}

	return nil
}

// ゴミ箱の中のファイル一覧を表示
func list(options *Options) error {
	trashPath := getTrashPath()

	dirs, err := ioutil.ReadDir(trashPath)
	if err != nil {
		return err
	}

	if options.reverse {
		sort.Sort(sort.Reverse(Dirs(dirs)))
	} else {
		sort.Sort(Dirs(dirs))
	}

	daysAgo := time.Now().AddDate(0, 0, -options.days)

	for _, dir := range dirs {
		internalStat, ok := dir.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("fileInfo.Sys(): cast error")
		}
		if options.days != 0 && internalStat.Ctim.Nano() < daysAgo.UnixNano() {
			continue
		}

		if err =
			printFiles(os.Stdout, filepath.Join(trashPath, dir.Name())); err != nil {
			return err
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
			return list(options)
		},
	}
	cmd.Flags().IntVarP(&options.days, "days", "d", 0, "How many days ago")
	cmd.Flags().BoolVarP(&options.reverse, "reverse", "r", false, "display in reverse order")

	return cmd
}
