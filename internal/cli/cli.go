package cli

import (
	"github.com/007kevin/bea/internal/bea"
	"github.com/urfave/cli/v2"
)

func Run(args []string) error {
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
					return bea.Generate()
				},
			},
		},
	}
	return app.Run(args)
}
