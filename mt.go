package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func createTrashCan(path string) error { // ゴミ箱が存在しないなら生成する。
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
	}

	return nil
}

func moveToTrashCan(files) { // ファイルをゴミ箱に移動させる
	// prefixの生成
	const layout = time.Now().Format("2006-01-02 15:04:05")
	const now string = strings.Replace(layout, " ", "_", 1)

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			log.Fatal(err)
		} else {
			if err := os.Rename(file, now+"_"+path.Base(file)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ls(path string) (files []string, err error) {
	files := make([]string, 0)

	fileInfo, err = ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
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

	const trashCanPath string = os.Getenv("HOME") + "/.Trash"

	if *list == true {
		files, err := ls(trashCanPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		fmt.Println(files)
	} else if *restore == true {
	} else if *size == true {
	} else if *delete == true {
	} else {
		moveToTrashCan(flag.Args())
	}

	if err := createTrashCan(trashCanPath); err != nil {
		log.Fatal(err) // [todo] log 種類調べる
		os.Exit(0)     // [todo] 番号を変える
	}
}
