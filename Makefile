IMAGES_TAG   = ${shell git describe --tags --match 'v[0-9]*\.[0-9]*\.[0-9]*' 2> /dev/null || echo 'latest'}
GIT_COMMIT   ?= $(shell git rev-parse HEAD)
GIT_TAG      ?= $(shell git tag --points-at HEAD)
DIST_TYPE    ?= snapshot
BRANCH       ?= $(shell git rev-parse --abbrev-ref HEAD)
REPO         ?= $(shell echo $(JOB_NAME) | cut -d/ -f2)
DATE         ?= $(shell date -u +%FT%T%z)

PROJECT_NAME        := gameserver-websocket
PKG_ORG             := github.com/jtbonhomme/$(PROJECT_NAME)
SERVER_BINARY_NAME  := server
CLIENT_BINARY_NAME  := client
SERVER_CMD          := cmd/$(SERVER_BINARY_NAME)
CLIENT_CMD          := cmd/$(CLIENT_BINARY_NAME)
SERVER_PKG 		    := $(PKG_ORG)/$(SERVER_CMD)
CLIENT_PKG 		    := $(PKG_ORG)/$(CLIENT_CMD)
#PKG_LIST 	        := $(shell go list ${PKG}/...)
GO_FILES 	        := $(shell find . -name '*.go' -not -path "./vendor/*" | grep -v _test.go)

GO			 := go
GOLANGCILINT := golangci-lint
GORELEASER	 := goreleaser
GOFMT		 := gofmt
GOIMPORTS	 := goimports
GCOV2LCOV    := gcov2lcov
GOCOV        := gocov
GOCOVXML     := gocov-xml
GOCOVHTML    := gocov-html
OS			 := $(shell uname -s)
GOOS		 ?= $(shell echo $OS | tr '[:upper:]' '[:lower:]')
GOARCH		 ?= amd64
DOCKER		 ?= docker


# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ; $(info $(M) Display makefile targets…) @ ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

linter: ; $(info $(M) Lint go source code…) @ ### check by golangci linter.
	@which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1)
	$(GOLANGCILINT) -v --deadline 100s --skip-dirs docs run ./...
.PHONY: linter

test: ; $(info $(M) Executing tests…)@ ### run tests.
	@which  $(GCOV2LCOV) || (go install github.com/jandelgado/gcov2lcov@latest)
	$(GO) test -race -cover -coverprofile=coverage.out ./... && \
		$(GCOV2LCOV) -infile=coverage.out -outfile=coverage.lcov
.PHONY: test

cover: test ; $(info $(M) Test coverage…)@ ## Measure the test coverage.
	which gocov || (go install github.com/axw/gocov/gocov@latest)
	which gocov-xml || (go install github.com/AlekSi/gocov-xml@latest)
	which gocov-html || (go install github.com/matm/gocov-html@latest)
	gocov convert coverage.out | gocov-xml > cover.xml
	gocov convert coverage.out | gocov-html > cover.html
.PHONY: cover

server: ; $(info $(M) Running server program…) @ ### run server program.
	$(GO) run $(SERVER_PKG)
.PHONY: server

client: ; $(info $(M) Running client program…) @ ### run client program.
	$(GO) run $(CLIENT_PKG)
.PHONY: client

download: ; $(info $(M) Downloading go dependencies…) @ ### downloads go dependencies.
	$(GO) mod download
.PHONY: download
