.SILENT:

all: build

build:
	go build -o output/bea cmd/main.go

run:
	./output/main
