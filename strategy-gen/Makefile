.PHONY: default strategy all clean docker

GOBIN = $(shell pwd)/build/bin
TAG ?= latest

default: strategy

all: strategy

strategy:
	go build $(BUILD_FLAGS) -o=${GOBIN}/$@ -gcflags "all=-N -l" .
	@echo "Done building."

clean:
	rm -fr build/*

docker:
	docker build -t strategy:${TAG} .
