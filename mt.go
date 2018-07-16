package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func createTrashCan(trashCanPath string) error { // ゴミ箱が存在しないなら生成する。
	if _, err := os.Stat(trashCanPath); err != nil {
		if err := os.Mkdir(trashCanPath, 0700); err != nil {
			return err
		}
	}

	return nil
}

func moveToTrashCan(trashCanPath string, files []string) { // ファイルをゴミ箱に移動させる
	prefix := trashCanPath + "/" + time.Now().Format("2006-01-02_15:04:05") + "_"

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			log.Fatal(err)
		} else {
			if err := os.Rename(file, prefix+path.Base(file)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func restore(trashCanPath string, files []string) {
	for _, file := range files {
	}
}

func ls(path string) (files []string, err error) {
	files = make([]string, 0)

	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	const executable os.FileMode = 0111
	const green = "\x1b[32m\x1b[1m%s"
	const blue = "\x1b[34m\x1b[1m%s"
	const cyan = "\x1b[36m\x1b[1m%s"
	const white = "\x1b[37m\x1b[0m%s"

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, fmt.Sprintf(blue, file.Name()))
		} else if file.Mode()&os.ModeSymlink != 0 {
			files = append(files, fmt.Sprintf(cyan, file.Name()))
		} else if file.Mode()&executable != 0 {
			files = append(files, fmt.Sprintf(green, file.Name()))
		} else {
			files = append(files, fmt.Sprintf(white, file.Name()))
		}
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("引数が足りません")
		os.Exit(0)
	}

	var (
		list    = flag.Bool("l", false, "list")
		restore = flag.Bool("r", false, "restore")
		size    = flag.Bool("s", false, "size")
		delete  = flag.Bool("d", false, "delete")
	)
	flag.Parse()
	if flag.NFlag() > 1 {
		fmt.Println("optionが多すぎます")
		os.Exit(0)
	}

	trashCanPath := os.Getenv("HOME") + "/.Trash"

	if *list == true {
		files, err := ls(trashCanPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		for _, file := range files {
			fmt.Println(file)
		}
	} else if *restore == true {
		fmt.Println("restore")
	} else if *size == true {
		fmt.Println("size")
	} else if *delete == true {
		fmt.Println("delete")
	} else {
		moveToTrashCan(trashCanPath, flag.Args())
	}

	if err := createTrashCan(trashCanPath); err != nil {
		log.Fatal(err) // [todo] log 種類調べる
		os.Exit(0)     // [todo] 番号を変える
	}
}
