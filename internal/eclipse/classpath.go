package eclipse

import (
	"encoding/xml"
)

type ClasspathOptions struct {
	SrcDirs []string
	TstDirs []string
	JarDirs []string
}

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
		{Name: "module", Value: "false"},
	},
}

var DefaultOutputEntry = ClasspathEntry{
	Kind: "output",
	Path: "bin",
}

func GenerateClasspath(options *ClasspathOptions) *Classpath {
	entries := []*ClasspathEntry{}
	entries = append(entries, &DefaultConEntry)
	entries = append(entries, SrcEntries(options.SrcDirs)...)
	entries = append(entries, TstEntries(options.TstDirs)...)
	entries = append(entries, JarEntries(options.JarDirs)...)
	entries = append(entries, &DefaultOutputEntry)
	return &Classpath{Entries: entries}
}

func SrcEntries(dirs []string) []*ClasspathEntry {
	var entries []*ClasspathEntry
	for _, dir := range dirs {
		entries = append(entries, &ClasspathEntry{
			Kind: "src",
			Path: dir,
		})
	}
	return entries
}

func TstEntries(dirs []string) []*ClasspathEntry {
	var entries []*ClasspathEntry
	for _, dir := range dirs {
		entries = append(entries, &ClasspathEntry{
			Kind: "src",
			Path: dir,
			Attributes: []*Attribute{
				{Name: "test", Value: "true"},
			},
		})
	}
	return entries
}

func JarEntries(dirs []string) []*ClasspathEntry {
	var entries []*ClasspathEntry
	for _, dir := range dirs {
		entries = append(entries, &ClasspathEntry{
			Kind: "lib",
			Path: dir,
		})
	}
	return entries
}

// <classpathentry kind="con" path="org.eclipse.jdt.launching.JRE_CONTAINER/org.eclipse.jdt.internal.debug.ui.launcher.StandardVMType/JavaSE-11">
// 	<attributes>
// 		<attribute name="module" value="true"/>
// 	</attributes>
// </classpathentry>
