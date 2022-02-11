package adaptor

import (
	"errors"
	"internal/adaptor/bazel"
	"internal/eclipse"
)

type Adaptor interface {
	Identifier() string
	Applicable() (bool, error)
	Generate() (*eclipse.Project, *eclipse.Classpath, error)
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
