.SILENT:

all: build

build:
	go build -o output/bin/bea cmd/main.go

run:
	./output/bin/bea
