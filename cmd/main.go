package main

import (
	"log"
	"os"

	"github.com/007kevin/bea/internal/cli"
)

func main() {
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
