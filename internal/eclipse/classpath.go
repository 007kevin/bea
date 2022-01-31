package eclipse

import "encoding/xml"

type Classpath struct {
	XMLName xml.Name          `xml:"classpath"`
	Entries []*ClasspathEntry `xml:"classpathentry"`
}

type ClasspathEntry struct {
	XMLName    xml.Name     `xml:"classpathentry"`
	Kind       string       `xml:"kind,attr,omitempty"`
	Path       string       `xml:"path,attr,omitempty"`
	Attributes []*Attribute `xml:"attributes,omitempty>attribute,omitempty"`
}

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	Name    string   `xml:"name,attr,omitempty"`
	Value   string   `xml:"value,attr,omitempty"`
}

var DefaultConEntry = ClasspathEntry{
	Kind: "con",
	Path: "org.eclipse.jdt.launching.JRE_CONTAINER/org.eclipse.jdt.internal.debug.ui.launcher.StandardVMType/JavaSE-11",
	Attributes: []*Attribute{
		{Name: "module", Value: "true"},
	},
}

// <classpathentry kind="con" path="org.eclipse.jdt.launching.JRE_CONTAINER/org.eclipse.jdt.internal.debug.ui.launcher.StandardVMType/JavaSE-11">
// 	<attributes>
// 		<attribute name="module" value="true"/>
// 	</attributes>
// </classpathentry>
