# Note

This directory holds protobuf definitions for https://github.com/bazelbuild/bazel to work with
bazel actions graphs with proto outputs. These files are slightly modified to work with golang.

analysis.proto is for parsing proto outputs for bazel versions < 5.0.0, analysis_v2.proto for all else.

See also:
- https://github.com/bazelbuild/bazel/tree/master/src/main/protobuf
- scripts/generate_golang_protobuf
