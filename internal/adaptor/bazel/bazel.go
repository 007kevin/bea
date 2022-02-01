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
	"path/filepath"
	"strings"
	"third_party/bazel/analysis"

	"github.com/liyue201/gostl/ds/set"
	"github.com/ojizero/gofindup"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/proto"
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
	bazelProtoAQuery("Javac", "--output", "proto_library")
	bazelProtoAQuery("Javac", "--classpath", "java_library", "java_test", "java_binary")
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
	output, err := CommandExec{}.runCommand("bazel", "query", "kind("+filter+",...)")
	return output.String(), err
}

func bazelProtoAQuery(mnemonic string, filter string, kinds ...string) (string, error) {
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
		return "", err
	}
	aqueryResult, err := bazelParseProto(output)
	if err != nil {
		return "", err
	}
	bazelReadDependencies(aqueryResult, filter)
	return output.String(), nil
}

func bazelParseProto(input *bytes.Buffer) (*analysis.ActionGraphContainer, error) {
	aqueryResult := &analysis.ActionGraphContainer{}
	err := proto.Unmarshal(input.Bytes(), aqueryResult)
	if err != nil {
		pterm.Error.Println(err)
		return nil, err
	}
	return aqueryResult, nil
}

func bazelReadDependencies(aqueryResult *analysis.ActionGraphContainer, argFilter string) []string {
	var argPaths = set.New()
	var outputIds = set.New()
	for _, action := range aqueryResult.Actions {
		var isArgFilter = false
		for _, argument := range action.Arguments {
			if isArgFilter && strings.HasPrefix(argument, "-") {
				isArgFilter = false
				continue
			}
			if !isArgFilter {
				isArgFilter = argument == argFilter
				continue
			}
			argPaths.Insert(argument)
		}
		for _, outputId := range action.OutputIds {
			outputIds.Insert(outputId)
		}
	}
	var artifactPaths []string
	pathFragments := map[uint32]*analysis.PathFragment{}
	for _, pathFragent := range aqueryResult.PathFragments {
		pathFragments[pathFragent.Id] = pathFragent
	}
	for _, artifact := range aqueryResult.Artifacts {
		var pathFragment = pathFragments[artifact.Id]
		if pathFragment != nil {
			var relative, err = expandPathFragment(pathFragment, pathFragments)
			if err != nil {
				continue
			}
			if !argPaths.Contains(relative) {
				// pterm.Warning.Println("...artifact was not specified by --filterArgument: '" + relative + "'")
				continue
			}
			if outputIds.Contains(artifact.Id) && argFilter != "--output" {
				// pterm.Warning.Println("...artifact is the output of another java action: '" + strconv.Itoa(int(artifact.Id)) + "'")
				continue
			}
			pterm.Info.Println("...found bazel dependency " + relative)
			artifactPaths = append(artifactPaths, relative)
		}
	}

	return artifactPaths
}

func expandPathFragment(pathFragment *analysis.PathFragment, pathFragments map[uint32]*analysis.PathFragment) (string, error) {
	labels := []string{}
	currId := pathFragment.Id
	// Only positive IDs are valid for path fragments. An ID of zero indicates a terminal node.
	for currId > 0 {
		currFragment, ok := pathFragments[currId]
		if !ok {
			return "", fmt.Errorf("undefined path fragment id %d", currId)
		}
		labels = append([]string{currFragment.Label}, labels...)
		currId = currFragment.ParentId
	}
	return filepath.Join(labels...), nil
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
	_, err = CommandExec{}.runCommand(args...)
	return err
}

type CommandExec struct {
	// Suppress stdout when running the command
	Suppress bool
}

func (options CommandExec) runCommand(args ...string) (*bytes.Buffer, error) {
	pterm.Info.Printf("Executing command: %v\n", args)
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
