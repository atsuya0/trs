package main

import (
	"log"
	"os"
)

func createTrashCan(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	trashCanPath := os.Getenv("HOME") + "/.Trash"

	err := createTrashCan(trashCanPath)
	if err != nil {
		log.Fatal(err)
	}
}
