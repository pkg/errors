#syntax=docker/dockerfile:1
#check=error=true

ARG GO_VERSION=1.26
ARG XX_VERSION=1.9.0

ARG COVER_FILENAME="cover.out"

FROM --platform=${BUILDPLATFORM} tonistiigi/xx:${XX_VERSION} AS xx

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine AS golang
COPY --link --from=xx / /
WORKDIR /src
ARG TARGETPLATFORM

FROM golang AS build
RUN --mount=target=/root/.cache,type=cache \
    --mount=type=bind xx-go build ./...

FROM golang AS runtest
ARG TESTFLAGS="-v"
ARG COVER_FILENAME
RUN --mount=target=/root/.cache,type=cache \
    --mount=type=bind \
    xx-go test -coverprofile=/tmp/${COVER_FILENAME} $TESTFLAGS ./...

FROM scratch AS test
ARG COVER_FILENAME
COPY --from=runtest /tmp/${COVER_FILENAME} /

FROM build
