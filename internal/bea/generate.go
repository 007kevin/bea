package bea

import (
	"adaptor"
	"fmt"
)

func generate() error {
	adaptor, err := adaptor.Get()
	if err != nil {
		return err
	}
	fmt.Println(adaptor.Applicable())
	return nil
}
