package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "bea",
		Version: "0.1",
		Usage:   "Bazel Eclipse Adapter",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage: `Generates Eclipse artifacts (i.e .project / .classpath) for a
                Bazel workspace to work with Eclipse IDEs (including the language server).`,
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
