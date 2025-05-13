.DEFAULT_GOAL := build

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "1.0.1")
BUILD_DATE := $(shell go run scripts/getbuilddate.go 2>/dev/null || echo "unknown")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X github.com/acheddir/git-contrib/cmd.Version=$(VERSION) -X github.com/acheddir/git-contrib/cmd.BuildDate=$(BUILD_DATE) -X github.com/acheddir/git-contrib/cmd.CommitHash=$(COMMIT_HASH)"

.PHONY:tidy fmt vet build
tidy:
	go mod tidy

fmt: tidy
	go fmt ./...

vet: fmt
	go vet ./...

build: tidy
	go build $(LDFLAGS)

clean:
	go clean
