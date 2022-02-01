.SILENT:

all: proto build

build:
	go build -o output/bin/bea cmd/main.go

proto:
	./scripts/generate_golang_protobuf

run:
	./output/bin/bea
