# See https://tech.davis-hansson.com/p/make/
SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory
BIN := .tmp/bin
COPYRIGHT_YEARS := 2021-2022
# Set to use a different compiler. For example, `GO=go1.18rc1 make test`.
GO ?= go

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.PHONY: all
all: ## Build, test, and lint (default)
	$(MAKE) test
	$(MAKE) lint

.PHONY: clean
clean: ## Delete intermediate build artifacts
	@# -X only removes untracked files, -d recurses into directories, -f actually removes files/dirs
	git clean -Xdf

.PHONY: test
test: build ## Run unit tests
	$(GO) test -vet=off -race -cover ./...

.PHONY: build
build: generate ## Build all packages
	$(GO) build ./...

.PHONY: lint
lint: $(BIN)/gofmt $(BIN)/buf ## Lint Go and protobuf
	test -z "$$($(BIN)/gofmt -s -l . | tee /dev/stderr)"
	test -z "$$($(BIN)/buf format -d . | tee /dev/stderr)"
	@# TODO: replace vet with golangci-lint when it supports 1.18
	@# Configure staticcheck to target the correct Go version and enable
	@# ST1020, ST1021, and ST1022.
	$(GO) vet ./...
	$(BIN)/buf lint

.PHONY: lintfix
lintfix: $(BIN)/gofmt $(BIN)/buf ## Automatically fix some lint errors
	$(BIN)/gofmt -s -w .
	$(BIN)/buf format -w .

.PHONY: generate
generate: $(BIN)/buf $(BIN)/protoc-gen-go $(BIN)/protoc-gen-connect-go $(BIN)/license-header ## Regenerate code and licenses
	rm -rf internal/gen
	PATH=$(BIN) $(BIN)/buf generate
	@$(BIN)/license-header \
		--license-type apache \
		--copyright-holder "Buf Technologies, Inc." \
		--year-range "$(COPYRIGHT_YEARS)" \
		$(shell git ls-files) $(shell git ls-files --exclude-standard --others)

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	go get -u -t ./... && go mod tidy -v

.PHONY: checkgenerate
checkgenerate:
	@# Used in CI to verify that `make generate` doesn't produce a diff.
	test -z "$$(git status --porcelain | tee /dev/stderr)"

$(BIN)/gofmt:
	@mkdir -p $(@D)
	$(GO) build -o $(@) cmd/gofmt

.PHONY: $(BIN)/protoc-gen-connect-go
$(BIN)/protoc-gen-connect-go:
	@mkdir -p $(@D)
	$(GO) build -o $(@) ./cmd/protoc-gen-connect-go

$(BIN)/buf: Makefile
	@mkdir -p $(@D)
	GOBIN=$(abspath $(@D)) $(GO) install github.com/bufbuild/buf/cmd/buf@v1.3.0

$(BIN)/license-header: Makefile
	@mkdir -p $(@D)
	GOBIN=$(abspath $(@D)) $(GO) install \
		  github.com/bufbuild/buf/private/pkg/licenseheader/cmd/license-header@v1.3.0

$(BIN)/protoc-gen-go: Makefile go.mod
	@mkdir -p $(@D)
	@# The version of protoc-gen-go is determined by the version in go.mod
	GOBIN=$(abspath $(@D)) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go

