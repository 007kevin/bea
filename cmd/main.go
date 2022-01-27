package main

import (
	"internal/bea"
	"log"
	"os"
)

func main() {
	err := bea.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
