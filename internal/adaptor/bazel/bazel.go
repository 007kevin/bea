package bazel

import (
	"encoding/xml"
	"errors"
	"fmt"
	"internal/eclipse"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
)

const FILE_MARKER = "BUILD"

type Adaptor struct {
}

func (ba *Adaptor) Applicable() (bool, error) {
	return dirHas(FILE_MARKER)
}

func (ba Adaptor) Identifier() string {
	return "Bazel Build Adaptor"
}

func (ba Adaptor) Run() error {
	fmt.Print(xml.Header)
	out, _ := xml.MarshalIndent(&eclipse.Classpath{
		Entries: []*eclipse.ClasspathEntry{
			&eclipse.DefaultConEntry,
			{},
			{},
		},
	}, "", "    ")
	fmt.Println(string(out))
	output, err := bazelQuery("java_proto_library")
	if err != nil {
		return err
	}
	lines := splitLines(output)
	for _, v := range lines {
		fmt.Println(v)

	}
	return nil
}

func dirHas(marker string) (bool, error) {
	if _, err := os.Stat(marker); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func bazelQuery(filter string) (string, error) {
	return runCommand("bazel", "query", "kind("+filter+",//...)")
}

func runCommand(args ...string) (string, error) {
	pterm.Info.Printf("Executing command: %v\n", args)
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(out), nil
}

func splitLines(str string) []string {
	return strings.Split(str, "\n")
}
