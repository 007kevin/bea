package main

import (
	"internal/cli"
	"os"

	"github.com/pterm/pterm"
)

func main() {
	err := cli.Run(os.Args)
	if err != nil {
		pterm.Error.Println(err)
	}
}
