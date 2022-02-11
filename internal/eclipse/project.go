package eclipse

import "encoding/xml"

type Project struct {
	XMLName xml.Name `xml:"projectDescription"`
}
