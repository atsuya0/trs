package cmd

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	executable os.FileMode = 0111
	green                  = "\x1b[1;32m%s\x1b[m\n"
	blue                   = "\x1b[1;34m%s\x1b[m\n"
	cyan                   = "\x1b[1;36m%s\x1b[m\n"
	white                  = "\x1b[0;37m%s\x1b[m\n"
)

type listOption struct {
	days    int
	size    string
	reverse bool
	path    bool
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

func printFile(file file, path bool) {
	var name string
	if path {
		name = file.path
	} else {
		name = file.info.Name()
	}

	if file.info.IsDir() {
		fmt.Printf(blue, name)
	} else if file.info.Mode()&os.ModeSymlink != 0 {
		fmt.Printf(cyan, name)
	} else if file.info.Mode()&executable != 0 {
		fmt.Printf(green, name)
	} else {
		fmt.Printf(white, name)
	}
}

func getFiles() (Files, error) {
	root, err := getTrashCanPath()
	if err != nil {
		return make(Files, 0), fmt.Errorf("%w", err)
	}

	var files Files
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if !info.IsDir() {
			files = append(files, file{info: info, path: path})
		}

		return nil
	}); err != nil {
		return make(Files, 0), fmt.Errorf("%w", err)
	}
	return files, nil
}

func list(option *listOption) error {
	files, err := getFiles()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if option.reverse {
		sort.Sort(sort.Reverse(files))
	} else {
		sort.Sort(Files(files))
	}

	days := time.Now().AddDate(0, 0, -option.days).UnixNano()
	size := convertSymbolsToNumbers(option.size)

	for _, file := range files {
		if option.days != 0 {
			if bool, err := file.withoutPeriod(days); bool {
				continue
			} else if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		if file.info.Size() < size {
			continue
		}
		printFile(file, option.path)
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
		"Show the files moved to the trash can within [days] days.")
	cmd.Flags().StringVarP(
		&option.size, "size", "s", "0B",
		"Show the files with size greater than [size].")
	cmd.Flags().BoolVarP(
		&option.reverse, "reverse", "r", false,
		"Show in reverse order")
	cmd.Flags().BoolVarP(
		&option.path, "path", "p", false,
		"Show the file paths")

	return cmd
}
