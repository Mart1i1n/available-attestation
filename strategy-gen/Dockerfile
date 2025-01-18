# syntax = docker/dockerfile:1-experimental
FROM golang:1.21-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

WORKDIR /build

COPY . /build/strategy-gen


RUN --mount=type=cache,target=/go/pkg/mod \
    cd /build/strategy-gen && go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/strategy-gen && go build -o strategy

FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

WORKDIR /root

COPY  --from=build /build/strategy-gen/strategy /usr/bin/strategy
RUN chmod u+x /usr/bin/strategy

ENTRYPOINT [ "strategy" ]