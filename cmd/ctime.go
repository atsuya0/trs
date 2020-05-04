// +build !mac

package cmd

import (
	"errors"
	"syscall"
)

func (f Files) Less(i, j int) bool {
	return f[i].info.Sys().(*syscall.Stat_t).Ctim.Nano() <
		f[j].info.Sys().(*syscall.Stat_t).Ctim.Nano()
}

func (f *file) withoutPeriod(sec int64) (bool, error) {
	internalStat, ok := f.info.Sys().(*syscall.Stat_t)
	if !ok {
		return false, errors.New("fileInfo.Sys(): cast error")
	}
	if internalStat.Ctim.Nano() < sec {
		return true, nil
	}
	return false, nil
}
