package bazel

import (
	"fmt"
	"path/filepath"
	"strings"
	"third_party/bazel/analysis"
	"third_party/bazel/analysis_v2"

	"github.com/liyue201/gostl/ds/set"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type AQueryResult interface {
	JavaDependencies() ([]string, error)
	Result() protoreflect.ProtoMessage
}

type AnalysisV2 struct {
	result *analysis_v2.ActionGraphContainer
}

func (analysis *AnalysisV2) Result() protoreflect.ProtoMessage {
	return analysis.result
}

func (analysis *AnalysisV2) JavaDependencies() ([]string, error) {
	var argFilter = "--classpath"
	var argPaths = set.New()
	var outputIds = set.New()
	for _, action := range analysis.result.Actions {
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
	pathFragments := map[uint32]*analysis_v2.PathFragment{}
	for _, pathFragent := range analysis.result.PathFragments {
		pathFragments[pathFragent.Id] = pathFragent
	}
	for _, artifact := range analysis.result.Artifacts {
		var pathFragment = pathFragments[artifact.PathFragmentId]
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
			fmt.Printf("INFO: found bazel dependency [%s]\n", relative)
			artifactPaths = append(artifactPaths, relative)
		}
	}

	return artifactPaths, nil
}

func expandPathFragment(pathFragment *analysis_v2.PathFragment, pathFragments map[uint32]*analysis_v2.PathFragment) (string, error) {
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

type Analysis struct {
	result *analysis.ActionGraphContainer
}

func (analysis *Analysis) Result() protoreflect.ProtoMessage {
	return analysis.result
}

func (analysis *Analysis) JavaDependencies() ([]string, error) {
	var argFilter = "--classpath"
	var argPaths = set.New()
	var outputIds = set.New()
	for _, action := range analysis.result.Actions {
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
	for _, artifact := range analysis.result.Artifacts {
		var relative = artifact.ExecPath
		if relative != "" {
			if !argPaths.Contains(relative) {
				// pterm.Warning.Println("...artifact was not specified by --filterArgument: '" + relative + "'")
				continue
			}
			if outputIds.Contains(artifact.Id) && argFilter != "--output" {
				// pterm.Warning.Println("...artifact is the output of another java action: '" + artifact.Id + "'")
				continue
			}
			fmt.Printf("INFO: found bazel dependency [%s]\n", relative)
			artifactPaths = append(artifactPaths, relative)
		}
	}
	return artifactPaths, nil
}
