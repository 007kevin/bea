package bea

import (
	"internal/adaptor"

	"github.com/pterm/pterm"
)

func Generate() error {
	adaptor, err := adaptor.Get()
	if err != nil {
		return err
	}
	pterm.Info.Println("Running " + adaptor.Identifier())
	return adaptor.Run()
}
