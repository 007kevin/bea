package bea

import (
	"internal/adaptor"
	"internal/util"

	"github.com/pterm/pterm"
)

func Generate() error {
	adaptor, err := adaptor.Get()
	if err != nil {
		return err
	}
	pterm.Info.Println("Running " + adaptor.Identifier())
	project, classpath, error := adaptor.Generate()
	if error != nil {
		return error
	}
	projectErr := util.WriteXml(".project", project)
	if projectErr != nil {
		return projectErr
	}
	classpathErr := util.WriteXml(".classpath", classpath)
	if classpathErr != nil {
		return classpathErr
	}
	return nil
}
