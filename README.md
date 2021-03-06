# Bazel Eclipse Adaptor - Generates Eclipse Artifacts For Bazel Projects

Generates Eclipse artifacts (i.e .project / .classpath) for a Bazel workspace to work with Eclipse IDEs (including the language server).

## Resources

- Standard Golang Project Layout Guide
  https://github.com/golang-standards/project-layout
- Go Playgound struct->xml example
  https://stackoverflow.com/questions/61603334/trying-to-use-golang-encoding-xml-to-encode-a-struct-to-xml-a-structure-with-an
  https://go.dev/play/p/Tl5o7s0-voW
- Sample eclipse projects for modelling artifacts
  https://github.com/eclipse/eclipse.jdt.ls/tree/master/org.eclipse.jdt.ls.tests/projects/eclipse
  https://github.com/eclipse/eclipse.jdt.ls/blob/master/org.eclipse.jdt.ls.tests/projects/eclipse/reference/.classpath
  https://github.com/eclipse/eclipse.jdt.ls/blob/master/org.eclipse.jdt.ls.tests/projects/eclipse/hello/.classpath
  explanation: https://programmer.group/detailed-explanation-of-classpath-file-in-eclipse-project.html
- Executing command and writing to file:
  https://yourbasic.org/golang/temporary-file-directory/
  https://stackoverflow.com/questions/18986943/in-golang-how-can-i-write-the-stdout-of-an-exec-cmd-to-a-file
- Investigation: https://www.diffchecker.com/fhrsdsLR

## TODO

- Newer bazel versions require analysis_v2.proto whereas old versions require analysis.proto. See:
  https://github.com/bazelbuild/bazel/blob/master/src/main/protobuf/analysis_v2.proto
  https://github.com/bazelbuild/bazel/blob/43bcbb623e241e6381149c08e2be653433cc9407/src/main/protobuf/analysis.proto
  Need to determine which bazel versions correspond to the different protos.
