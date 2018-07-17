package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

func moveToTrashCan(trashCanPath string, files []string) { // ファイルをゴミ箱に移動させる
	prefix := trashCanPath + "/" + time.Now().Format("2006-01-02_15:04:05") + "_"

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			log.Println(err)
			continue
		}

		if err := os.Rename(file, prefix+path.Base(file)); err != nil {
			log.Println(err)
		}
	}
}

func currentDirNames() ([]string, error) { // カレントディレクトリのファイル・ディレクトリ名の一覧
	var files []string

	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return files, err
	}

	file, err := os.Open(wd)
	defer file.Close()

	if err != nil {
		log.Println(err)
		return files, err
	}

	files, err = file.Readdirnames(0)
	if err != nil {
		log.Println(err)
		return files, err
	}

	return files, err
}

func contains(file string, files []string) bool {
	for _, v := range files {
		if file == v {
			return true
		}
	}
	return false
}

// ゴミ箱からファイルを取り出す
func restore(trashCanPath string, trashFiles []string) error {
	files, err := currentDirNames()
	if err != nil {
		return err
	}

	for _, fileName := range trashFiles {
		filePath := trashCanPath + "/" + fileName
		if _, err := os.Stat(filePath); err != nil {
			log.Println(err)
			continue
		}

		index1 := strings.Index(fileName, "_")
		index2 := strings.Index(fileName[index1+1:], "_")
		newFileName := fileName[index1+index2+2:]

		if contains(newFileName, files) {
			log.Println("同じファイル名のファイルがあります")
			continue
		}
		if err := os.Rename(filePath, newFileName); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func list(path string) (files []string, err error) { // ゴミ箱の中のファイル一覧を表示
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

func size(root string) (int64, error) {
	var sum int64 = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			sum += info.Size()
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func del(path string, file string) error {
	fmt.Printf("target: %s", file)
	fmt.Println("'yes' or 'no'")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	for scanner.Text() != "yes" && scanner.Text() != "no" {
		fmt.Println("'yes' or 'no'")
		scanner.Scan()
	}

	if scanner.Text() == "yes" {
		if err := os.Remove(path + "/" + file); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("引数が足りません")
	}

	var (
		l = flag.Bool("l", false, "list")
		r = flag.Bool("r", false, "restore")
		s = flag.Bool("s", false, "size")
		d = flag.Bool("d", false, "delete")
	)
	flag.Parse()
	if flag.NFlag() > 1 {
		log.Fatalln("optionが多すぎます")
	}

	trashCanPath := os.Getenv("HOME") + "/.Trash"

	if err := createTrashCan(trashCanPath); err != nil {
		log.Fatalln(err)
	}

	if *l == true {
		files, err := list(trashCanPath)
		if err != nil {
			log.Fatalln(err)
		}
		for _, file := range files {
			fmt.Println(file)
		}
	} else if *r == true {
		if err := restore(trashCanPath, flag.Args()); err != nil {
			log.Fatalln(err)
		}
	} else if *s == true {
		trashCanSize, err := size(trashCanPath)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%d MB", trashCanSize/(1024*1024))
	} else if *d == true {
		if err := del(trashCanPath, flag.Args()[0]); err != nil {
			log.Fatalln(err)
		}
	} else {
		moveToTrashCan(trashCanPath, flag.Args())
	}
}
