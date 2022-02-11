package eclipse

import (
	"encoding/xml"
)

type Project struct {
	XMLName   xml.Name        `xml:"projectDescription"`
	Name      string          `xml:"name"`
	Comment   string          `xml:"comment"`
	Projects  string          `xml:"projects"`
	BuildSpec []*BuildCommand `xml:"buildSpec>buildCommand"`
	Natures   []string        `xml:"natures>nature"`
}

type BuildCommand struct {
	XMLName   xml.Name `xml:"buildCommand"`
	Name      string   `xml:"name"`
	Arguments string   `xml:"arguments"`
}

func GenerateProject(name string) *Project {
	return &Project{
		Name: name,
		BuildSpec: []*BuildCommand{
			{Name: "org.eclipse.jdt.core.javabuilder"},
		},
		Natures: []string{
			"org.eclipse.jdt.core.javanature",
		},
	}
}
