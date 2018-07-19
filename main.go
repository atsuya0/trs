package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func createTrashCan(trashCanPath string) error { // ゴミ箱が存在しないなら生成する。
	if _, err := os.Stat(trashCanPath); err == nil {
		return nil
	}

	if err := os.Mkdir(trashCanPath, 0700); err != nil {
		return err
	} else {
		return nil
	}

}

// optionの数が多いか調べる
func isDuplicatedOptions() bool {
	return flag.NFlag() > 1
}

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("引数が足りません")
	}
}

func main() {
	var (
		l       = flag.Bool("l", false, "list")
		r       = flag.Bool("r", false, "restore")
		s       = flag.Bool("s", false, "size")
		d       = flag.Bool("d", false, "delete")
		day     = flag.Int("day", 1, "[n] day ago")
		reverse = flag.Bool("reverse", false, "reverse")
	)
	flag.Parse()

	trashCanPath := os.Getenv("HOME") + "/.Trash"

	if err := createTrashCan(trashCanPath); err != nil {
		log.Fatalln(err)
	}

	if *l == true {
		files, err := list(trashCanPath, *day, *reverse)
		if err != nil {
			log.Fatalln(err)
		}
		for _, file := range files {
			fmt.Println(file)
		}
	} else if *r == true {
		if isDuplicatedOptions() {
			log.Fatalln("optionが不正です")
		}

		if err := restore(trashCanPath, flag.Args()); err != nil {
			log.Fatalln(err)
		}
	} else if *s == true {
		if isDuplicatedOptions() {
			log.Fatalln("optionが不正です")
		}

		trashCanSize, err := size(trashCanPath)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%d MB", trashCanSize/(1024*1024))
	} else if *d == true {
		if isDuplicatedOptions() {
			log.Fatalln("optionが不正です")
		}

		if err := del(trashCanPath, flag.Args()[0]); err != nil {
			log.Fatalln(err)
		}
	} else {
		moveToTrashCan(trashCanPath, flag.Args())
	}
}
