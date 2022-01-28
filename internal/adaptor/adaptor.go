package adaptor

import (
	"errors"
	"internal/adaptor/bazel"
)

type Adaptor interface {
	Identifier() string
	Applicable() (bool, error)
	Run() error
}

var Adaptors = []Adaptor{
	&bazel.Adaptor{},
}

func Get() (Adaptor, error) {
	for _, adaptor := range Adaptors {
		if yes, err := adaptor.Applicable(); yes && err == nil {
			return adaptor, err
		} else if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("No suitable build adaptor found. Are you in a project directory?")
}
