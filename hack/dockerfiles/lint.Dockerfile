#syntax=docker/dockerfile:1
#check=error=true

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.23
ARG XX_VERSION=1.9.0
ARG GOLANGCI_LINT_VERSION=2.11.4

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS golang-base
FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx

FROM golang-base AS base
ENV GOFLAGS="-buildvcs=false"
ARG GOLANGCI_LINT_VERSION
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_LINT_VERSION}
COPY --link --from=xx / /
WORKDIR /go/src/github.com/pkg/errors

FROM base AS golangci-lint
ARG TARGETPLATFORM
RUN --mount=type=bind \
    --mount=target=/root/.cache,type=cache \
  xx-go --wrap && \
  GOROOT=$(xx-go env GOROOT) golangci-lint run && \
  touch /golangci-lint.done

FROM scratch
COPY --link --from=golangci-lint /golangci-lint.done /
