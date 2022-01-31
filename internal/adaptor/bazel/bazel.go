package bazel

import (
	"encoding/xml"
	"errors"
	"fmt"
	"internal/eclipse"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ojizero/gofindup"
	"github.com/pterm/pterm"
)

const FILE_MARKER = "BUILD"
const WORKSPACE_MARKER = "WORKSPACE"

type Adaptor struct {
}

func (ba *Adaptor) Applicable() (bool, error) {
	return dirHas(FILE_MARKER)
}

func (ba *Adaptor) Identifier() string {
	return "Bazel Build Adaptor"
}

func (ba *Adaptor) Run() error {
	fmt.Print(xml.Header)
	out, _ := xml.MarshalIndent(&eclipse.Classpath{
		Entries: []*eclipse.ClasspathEntry{
			&eclipse.DefaultConEntry,
			{},
			{},
		},
	}, "", "    ")
	fmt.Println(string(out))
	fmt.Println(findWorkspaceRoot())
	return buildProtos()
}

func findWorkspaceRoot() (string, error) {
	p, err := gofindup.Findup(WORKSPACE_MARKER)
	if err != nil {
		return "", err
	}
	return path.Dir(p), nil
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
	return runCommand("bazel", "query", "kind("+filter+",...)")
}

func buildProtos() error {
	output, err := bazelQuery("java_proto_library")
	if err != nil {
		return err
	}
	lines := splitLines(output)
	args := append([]string{"bazel", "build", "--nobuild"}, lines...)
	_, err = runCommand(args...)
	pterm.Info.Println("HERE")
	pterm.Info.Println(err)
	return err
}

func runCommand(args ...string) (string, error) {
	pterm.Info.Printf("Executing command: %v\n", args)
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		pterm.Error.Println(err)
		return "", err
	}
	return string(out), nil
}

func splitLines(str string) []string {
	return strings.Split(str, "\n")
}
