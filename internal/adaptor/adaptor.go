package adaptor

import (
	"errors"

	"github.com/007kevin/bea/internal/adaptor/bazel"
)

type Adaptor interface {
	Applicable() bool
}

var Adaptors = []Adaptor{
	&bazel.Adaptor{},
}

func Get() (Adaptor, error) {
	for _, adaptor := range Adaptors {
		if adaptor.Applicable() {
			return adaptor, nil
		}
	}
	return nil, errors.New("No suitable adaptor found")
}
