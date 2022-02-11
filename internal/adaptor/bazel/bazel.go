package bazel

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"internal/eclipse"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"third_party/bazel/analysis"
	"third_party/bazel/analysis_v2"

	"github.com/liyue201/gostl/ds/set"
	"github.com/ojizero/gofindup"
	"github.com/pterm/pterm"
	"google.golang.org/protobuf/proto"
)

const BUILD_MARKER = "BUILD"
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
	r := regexp.MustCompile("[0-9.]+")
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
	return dirHas(BUILD_MARKER)
}

func (ba *Adaptor) Identifier() string {
	return "Bazel Build Adaptor"
}

func (ba *Adaptor) Generate() (*eclipse.Project, *eclipse.Classpath, error) {
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
		return nil, nil, ss.err
	}
	srcDirs, tstDirs, err := bazelJavaDirs()
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("---Source---")
	for _, dep := range srcDirs {
		fmt.Println(dep)
	}
	fmt.Println("---Test---")
	for _, dep := range tstDirs {
		fmt.Println(dep)
	}
	return &eclipse.Project{}, &eclipse.Classpath{}, nil
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

func findBuildRoot() (string, error) {
	p, err := gofindup.Findup(BUILD_MARKER)
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

func addWorkspaceRoot(dependencies []string) []string {
	root, err := findWorkspaceRoot()
	if err != nil {
		pterm.Error.Println("Unable to determine workspace root: ", err)
		return nil
	}
	result := make([]string, len(dependencies))
	for i, dep := range dependencies {
		result[i] = path.Join(root, dep)
	}
	return result
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
	dependencies, err := aQueryResult.Dependencies()
	return addWorkspaceRoot(dependencies), err
}

func bazelJavaDeps() ([]string, error) {
	aQueryResult, err := bazelProtoAQuery("Javac", "--classpath", "java_library", "java_test", "java_binary")
	if err != nil {
		return nil, err
	}
	dependencies, err := aQueryResult.Dependencies()
	return addWorkspaceRoot(dependencies), err
}

func bazelJavaDirs() ([]string, []string, error) {
	root, err := findBuildRoot()
	if err != nil {
		return nil, nil, err
	}
	pterm.Info.Printf("Normalizing all java files in %s. This may take awhile...\n", root)
	javaFiles := find(root, ".java")
	var srcDirs []string
	var tstDirs []string
	for _, dir := range normalizeDirs(javaFiles) {
		if isSrcDirectory(dir) {
			srcDirs = append(srcDirs, dir)
		} else {
			tstDirs = append(tstDirs, dir)
		}
	}
	return srcDirs, tstDirs, nil
}

// Use heuristics to determine whether source or test directory
func isSrcDirectory(str string) bool {
	str = strings.ToLower(str)
	if strings.Contains(str, "src/main/java") {
		return true
	}
	if strings.Contains(str, "src/test/java") {
		return false
	}
	if strings.Contains(str, "source") {
		return true
	}
	if strings.Contains(str, "test") {
		return false
	}
	if strings.Contains(str, "src") {
		return true
	}
	if strings.Contains(str, "tst") {
		return false
	}
	pterm.Warning.Printf("Unable to determine source type for %s. Defaulting to true\n", str)
	return true
}

func find(root, ext string) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func normalizeDirs(javaFiles []string) []string {
	var dirs = set.New()
	buildRoot, err := findBuildRoot()
	if err != nil {
		pterm.Fatal.Println(err)
	}
	for _, file := range javaFiles {
		dir := filepath.Dir(file)
		pkg := extractPackage(file)
		dirs.Insert(strings.TrimPrefix(strings.TrimSuffix(dir, pkg), buildRoot))
	}
	var res []string
	for iter := dirs.Begin(); iter.IsValid(); iter.Next() {
		res = append(res, fmt.Sprintf("%v", iter.Value()))
	}
	return res
}

func extractPackage(javaFile string) string {
	file, err := os.Open(javaFile)
	if err != nil {
		pterm.Error.Println(err)
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	reg := regexp.MustCompile(`\s*package\s+(\S+);`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := reg.FindStringSubmatch(line)
		if len(matches) > 1 {
			return strings.ReplaceAll(matches[1], ".", string(os.PathSeparator))
		}
	}
	return ""
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
