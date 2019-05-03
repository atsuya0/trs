// +build !mac

package cmd

import (
	"os"
	"syscall"
)

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

func greaterThenCtime(stat *syscall.Stat_t, sec int64) bool {
	if sec > stat.Ctim.Nano() {
		return true
	}
	return false
}
