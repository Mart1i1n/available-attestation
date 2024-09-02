.PHONY: default txpress all clean fmt docker

GOBIN = $(shell pwd)/build/bin
TAG ?= latest
GOFILES_NOVENDOR := $(shell go list -f "{{.Dir}}" ./...)

VERSION := $(shell git describe --tags)
COMMIT_SHA1 := $(shell git rev-parse HEAD)
AppName := txpress

default: txpress

all: txpress

BUILD_FLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT_SHA1}"

txpress:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" ./
	@echo "Done building."


clean:
	rm -fr build/*

docker:
	docker build -t txpress:${TAG} .
