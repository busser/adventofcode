.DEFAULT_TARGET = help

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## build: Build adventofcode binary
build: fmt vet
	go build -o bin/adventofcode main.go

## test: Run tests
test: fmt vet
	go test ./... -cover

## bench: Run tests & benchmarks
bench: fmt vet
	go test ./... -cover -bench . -benchmem -cpu 1,2,4,8

## run: Run adventofcode CLI
run: fmt vet
	go run ./main.go

## fmt: Run go fmt against code
fmt:
	go fmt ./...

## vet: Run go vet against code
vet:
	go vet ./...
