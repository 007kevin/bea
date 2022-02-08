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
	"regexp"
	"strconv"
	"strings"
	"third_party/bazel/analysis"
	"third_party/bazel/analysis_v2"

	"github.com/ojizero/gofindup"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/proto"
)

const FILE_MARKER = "BUILD"
const WORKSPACE_MARKER = "WORKSPACE"

var bazelVersion *Version

func init() {
	output, err := CommandExec{Suppress: true}.runCommand("bazel", "--version")
	if err != nil {
		pterm.Fatal.Println(err)
	}
	pterm.Info.Println("Running " + output.String())
	bazelVersion, err = parseVersion(output.String())
	if err != nil {
		pterm.Fatal.Println(err)
	}
}

type Version struct {
	numeric []int
}

func (version *Version) useV2() bool {
	if len(version.numeric) == 0 {
		pterm.Warning.Println("Version not found. Defaulting to V2")
		return true
	}
	return version.numeric[0] >= 5 // Just need to check the major version.
}

func parseVersion(raw string) (*Version, error) {
	r, err := regexp.Compile("[0-9.]+")
	if err != nil {
		return nil, err
	}
	numeric := make([]int, 0)
	for _, value := range strings.Split(r.FindString(raw), ".") {
		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		numeric = append(numeric, number)
	}
	return &Version{numeric: numeric}, nil
}

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
	var ss StringSlice
	ss.append(bazelJavaProtos())
	ss.append(bazelJavaDeps())
	if ss.err != nil {
		return ss.err
	}
	for _, dep := range ss.slice {
		fmt.Println(dep)
	}
	// bazelProtoAQuery("Javac", "--source_jars", "proto_library")
	// bazelProtoAQuery("JavaSourceJar", "--sources", "java_library", "java_test", "java_binary")
	return nil
}

type StringSlice struct {
	err   error
	slice []string
}

func (ss *StringSlice) append(slice []string, err error) *StringSlice {
	if ss.err != nil {
		return ss
	}
	ss.slice = append(ss.slice, slice...)
	return ss
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
	output, err := CommandExec{}.runCommand("bazel", "query", "kind("+filter+",...)")
	return output.String(), err
}

func bazelJavaProtos() ([]string, error) {
	pterm.Info.Println("Building java protos")
	output, err := bazelQuery("java_proto_library")
	if err != nil {
		return nil, err
	}
	lines := splitLines(output)
	if len(lines) == 0 {
		pterm.Info.Println("No protos found. Skipping.")
		return nil, nil
	}
	args := append([]string{"bazel", "build", "--nobuild"}, lines...)
	_, err = CommandExec{}.runCommand(args...)
	aQueryResult, err := bazelProtoAQuery("Javac", "--output", "proto_library")
	if err != nil {
		return nil, err
	}
	return aQueryResult.Dependencies()
}

func bazelJavaDeps() ([]string, error) {
	aQueryResult, err := bazelProtoAQuery("Javac", "--classpath", "java_library", "java_test", "java_binary")
	if err != nil {
		return nil, err
	}
	return aQueryResult.Dependencies()
}

func bazelProtoAQuery(mnemonic string, filter string, kinds ...string) (AQueryResult, error) {
	var kindUnion []string
	for _, kind := range kinds {
		kindUnion = append(kindUnion, "kind("+kind+", ...)")
	}
	output, err := CommandExec{Suppress: true}.runCommand(
		"bazel",
		"aquery",
		"--output=proto",
		"--include_aspects",
		"--allow_analysis_failures",
		"mnemonic("+mnemonic+", "+strings.Join(kindUnion, " union ")+")",
	)
	if err != nil {
		return nil, err
	}
	return bazelParseProto(output)
}

func bazelParseProto(input *bytes.Buffer) (AQueryResult, error) {
	aQuery := newAQueryResult()
	err := proto.Unmarshal(input.Bytes(), aQuery.Result())
	if err != nil {
		pterm.Error.Println(err)
		return nil, err
	}
	return aQuery, nil
}

func newAQueryResult() AQueryResult {
	if bazelVersion.useV2() {
		return &AnalysisV2{result: &analysis_v2.ActionGraphContainer{}}
	} else {
		return &Analysis{result: &analysis.ActionGraphContainer{}}
	}
}

type CommandExec struct {
	// Suppress stdout when running the command
	Suppress bool
}

func (options CommandExec) runCommand(args ...string) (*bytes.Buffer, error) {
	pterm.Info.Printf("%s\n", strings.Join(args, " "))
	var cmd = exec.Command(args[0], args[1:]...)
	var buffer = &bytes.Buffer{}

	if options.Suppress {
		cmd.Stdout = buffer
	} else {
		// write output to also terminal so user can see.
		cmd.Stdout = io.MultiWriter(os.Stdout, buffer)
	}
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		pterm.Error.Println(err)
		return nil, err
	}
	return buffer, nil
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
