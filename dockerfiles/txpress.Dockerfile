# syntax = docker/dockerfile:1-experimental
FROM golang:1.21-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make build-base

WORKDIR /build

COPY ./txpress /build/txpress


RUN --mount=type=cache,target=/go/pkg/mod \
    cd /build/txpress && go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/txpress && make

FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add --no-cache build-base

WORKDIR /root

COPY  --from=build /build/txpress/build/bin/txpress /usr/bin/txpress
COPY ./config/txpress-app.json /root/app.json

ENTRYPOINT [ "txpress" ]
