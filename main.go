package main

import (
	"log"

	"github.com/atsuya0/trs/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
