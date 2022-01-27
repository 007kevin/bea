module github.com/007kevin/bea

go 1.17

require (
	adaptor v0.0.0-00010101000000-000000000000 // indirect
	bazel v0.0.0-00010101000000-000000000000 // indirect
	bea v0.0.0-00010101000000-000000000000 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/urfave/cli/v2 v2.3.0 // indirect
)

replace (
	adaptor => ./internal/adaptor
	bazel => ./internal/adaptor/bazel
	bea => ./internal/bea
)
