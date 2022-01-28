package bea

import (
	"fmt"

	"github.com/007kevin/bea/internal/adaptor"
)

func Generate() error {
	adaptor, err := adaptor.Get()
	if err != nil {
		return err
	}
	fmt.Println(adaptor.Applicable())
	return nil
}
