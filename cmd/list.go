package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

const (
	executable os.FileMode = 0111
	header                 = "\x1b[7;39;49m%s\x1b[m\n"
	green                  = "\x1b[1;32m%s\x1b[m\n"
	blue                   = "\x1b[1;34m%s\x1b[m\n"
	cyan                   = "\x1b[1;36m%s\x1b[m\n"
	white                  = "\x1b[0;37m%s\x1b[m\n"
)

type listOption struct {
	days    int
	size    string
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

func convertSymbolsToNumbers(size string) int64 {
	for i, unit := range units {
		idx := strings.LastIndex(size, unit)
		if 0 < idx {
			num, err := strconv.Atoi(string([]rune(size)[:idx]))
			if err != nil {
				continue
			}
			return int64(num) * int64(math.Pow(1024, float64(i)))
		}
	}
	return 0
}

func printFiles(out io.Writer, path string, size int64) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, header, filepath.Base(path))

	for _, file := range files {
		if file.Size() < size {
			continue
		}
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

func list(option *listOption) error {
	path, err := getTrashCanPath()
	if err != nil {
		return err
	}

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	if option.reverse {
		sort.Sort(sort.Reverse(Dirs(dirs)))
	} else {
		sort.Sort(Dirs(dirs))
	}

	daysAgo := time.Now().AddDate(0, 0, -option.days)
	size := convertSymbolsToNumbers(option.size)

	for _, dir := range dirs {
		internalStat, ok := dir.Sys().(*syscall.Stat_t)
		if !ok {
			return xerrors.New("fileInfo.Sys(): cast error")
		}
		if option.days != 0 && internalStat.Ctim.Nano() < daysAgo.UnixNano() {
			continue
		}

		if err = printFiles(os.Stdout,
			filepath.Join(path, dir.Name()), size); err != nil {
			return err
		}
	}

	return nil
}

func listCmd() *cobra.Command {
	option := &listOption{}

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List the files in the trash can",
		RunE: func(_ *cobra.Command, _ []string) error {
			return list(option)
		},
	}
	cmd.Flags().IntVarP(
		&option.days, "days", "d", 0,
		"List the files moved to the trash can within [days] days.")
	cmd.Flags().StringVarP(
		&option.size, "size", "s", "0B",
		"List the files with size greater than [size].")
	cmd.Flags().BoolVarP(
		&option.reverse, "reverse", "r", false,
		"List in reverse order")

	return cmd
}
