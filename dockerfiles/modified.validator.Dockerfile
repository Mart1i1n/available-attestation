# syntax = docker/dockerfile:1-experimental
FROM golang:1.21-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

RUN go env -w CGO_ENABLED="1"

WORKDIR /build

COPY ./modified_prysm /build/prysm

RUN --mount=type=cache,target=/go/pkg/mod \
    cd /build/prysm && go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/prysm && go build -o /validator ./cmd/validator

FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

WORKDIR /root

COPY  --from=build /validator /usr/bin/validator
COPY ./entrypoint/validator.sh /usr/local/bin/validator.sh
RUN chmod u+x /usr/local/bin/validator.sh

ENTRYPOINT [ "/usr/local/bin/validator.sh" ]
