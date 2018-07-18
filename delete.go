package main

import (
	"bufio"
	"fmt"
	"os"
)

// ゴミ箱に入っている、指定した一つのファイルを削除する。
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
