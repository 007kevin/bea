package bazel

import (
	"errors"
	"os"
)

const FILE_MARKER = "BUILD"

type Adaptor struct {
}

func (ba *Adaptor) Applicable() (bool, error) {
	return dirHas(".", FILE_MARKER)
}

func (ba Adaptor) Identifier() string {
	return "Bazel Build Adaptor"
}

func dirHas(dir, marker string) (bool, error) {
	if _, err := os.Stat(marker); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
