package main

import (
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func createTrashCan(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
	}

	return nil
}

func moveToTrashCan(files) {
	const now string = strings.Replace(time.Now().Format("2006-01-02 15:04:05"), " ", "_", 1)

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("引数が足りません")
		os.Exit(0)
	}

	const trashCanPath string = os.Getenv("HOME") + "/.Trash"

	err := createTrashCan(trashCanPath)
	if err != nil {
		log.Fatal(err) // [todo] log 種類調べる
		os.Exit(0)     // [todo] 番号を変える
	}

	moveToTrashCan(os.Args)
}
