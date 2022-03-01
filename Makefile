MAKEGO := make/go
MAKEGO_REMOTE := https://github.com/bufbuild/makego.git
PROJECT := connect
GO_MODULE := github.com/bufbuild/connect
GO_MOD_VERSION := 1.18
GO_BINS := $(GO_BINS) cmd/protoc-gen-connect-go
GO_ALL_REPO_PKGS := ./...
GIT_FILE_IGNORES := $(GIT_FILE_IGNORES) \
	cover.out \
	*.pprof \
	*.svg
# TODO: remove when golangci-lint works with 1.18
SKIP_GOLANGCI_LINT := 1

LICENSE_HEADER_LICENSE_TYPE := apache
LICENSE_HEADER_COPYRIGHT_HOLDER := Buf Technologies, Inc.
LICENSE_HEADER_YEAR_RANGE := 2021-2022
LICENSE_HEADER_IGNORES := \/testdata

BUF_LINT_INPUT := .

include make/go/bootstrap.mk
include make/go/go.mk
include make/go/buf.mk
include make/go/license_header.mk
include make/go/dep_protoc_gen_go.mk
include make/go/dep_protoc_gen_go_grpc.mk

bufgeneratedeps:: $(BUF) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC) installprotoc-gen-connect-go

.PHONY: bufgeneratecleango
bufgeneratecleango:
	rm -rf internal/gen/proto

bufgenerateclean:: bufgeneratecleango

.PHONY: bufgeneratego
bufgeneratego:
	buf generate

bufgeneratesteps:: bufgeneratego
