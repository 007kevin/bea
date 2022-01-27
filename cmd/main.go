package main

import (
	"bea"
	"log"
	"os"
)

func main() {
	err := bea.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
