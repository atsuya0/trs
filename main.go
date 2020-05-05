package main

import (
	"log"

	"github.com/tayusa/trs/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
