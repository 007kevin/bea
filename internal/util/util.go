package util

import (
	"encoding/xml"
	"fmt"
)

func MarshallXml(v interface{}) (string, error) {
	output, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n%s", xml.Header, string(output)), nil

}
