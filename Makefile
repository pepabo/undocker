VERSION ?= $(shell git describe --tag --abbrev=0)

export GO111MODULE=on

default: test

ci: depsdev build test sec

build:
	go build --ldflags "-s -w -X main.version=$(VERSION)" ./cmd/undocker

test:
	go test ./... -v -coverprofile=coverage.txt -covermode=count

test_on_docker_on_mac:
	docker run --add-host=localhost:`ipconfig getifaddr en0` --rm -it -v "$(PWD)":/go/src/github.com/pepabo/undocker -w /go/src/github.com/pepabo/undocker golang:latest go test ./... -v

sec:
	gosec -exclude=G110 ./...

depsdev:
	go get golang.org/x/tools/cmd/cover
	go get github.com/securego/gosec/cmd/gosec

.PHONY: default
