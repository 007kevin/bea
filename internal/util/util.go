package util

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

func MarshallXml(v interface{}) (string, error) {
	output, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", xml.Header, string(output)), nil

}

func WriteXml(xmlFileName string, xml interface{}) error {
	output, marshallErr := MarshallXml(xml)
	if marshallErr != nil {
		return marshallErr
	}
	file, createErr := os.Create(xmlFileName)
	if createErr != nil {
		return createErr
	}
	defer file.Close()
	_, writeErr := file.WriteString(output)
	if writeErr != nil {
		return writeErr
	}
	pterm.Info.Println("Successfully wrote to " + file.Name())
	return nil
}
