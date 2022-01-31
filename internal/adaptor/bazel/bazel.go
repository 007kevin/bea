package bazel

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"internal/eclipse"
	"io"
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
	buildProtos()
	bazelAQuery("Javac", "--output", "proto_library")
	return nil
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

func bazelAQuery(mnemonic string, filter string, kinds ...string) (string, error) {
	var kindUnion []string
	for _, kind := range kinds {
		kindUnion = append(kindUnion, "kind("+kind+", ...)")
	}
	output, err := runCommand(
		"bazel",
		"aquery",
		"--include_aspects",
		"mnemonic("+mnemonic+", "+strings.Join(kindUnion, "union")+")",
	)
	if err != nil {
		return "", err
	}
	return output, nil
}

func buildProtos() error {
	pterm.Info.Println("Building java protos")
	output, err := bazelQuery("java_proto_library")
	if err != nil {
		return err
	}
	lines := splitLines(output)
	if len(lines) == 0 {
		pterm.Info.Println("No protos found. Skipping.")
		return nil
	}
	args := append([]string{"bazel", "build", "--nobuild"}, lines...)
	_, err = runCommand(args...)
	return err
}

func runCommand(args ...string) (string, error) {
	pterm.Info.Printf("Executing command: %v\n", args)
	var cmd = exec.Command(args[0], args[1:]...)
	var stdBuffer bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdBuffer) // write output to also terminal so user can see.
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		pterm.Error.Println(err)
		return "", err
	}
	return stdBuffer.String(), nil
}

func splitLines(input string) []string {
	var output []string
	for _, str := range strings.Split(input, "\n") {
		str = strings.TrimSpace(str)
		if str != "" {
			output = append(output, str)
		}
	}
	return output
}
