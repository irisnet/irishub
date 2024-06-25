#!/usr/bin/make -f

SIMAPP = ./simapp
BINDIR ?= $(GOPATH)/bin
CURRENT_DIR = $(shell pwd)

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

ldflags = -X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-zA-Z][a-zA-Z0-9/:]{2,127}

all: tools lint

# The below include contains the tools.
include contrib/devtools/Makefile

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

clean:
	rm -rf snapcraft-local.yaml build/

distclean: clean
	rm -rf vendor/

###############################################################################
###                                Protobuf                                 ###
###############################################################################
include scripts/build/protobuf.mk

########################################
### Testing

include scripts/build/testing.mk

include scripts/build/linting.mk

# lint: golangci-lint
# 	golangci-lint run
# 	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs gofmt -d -s
# 	go mod verify

# format:
# 	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs gofmt -w -s
# 	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" | xargs misspell -w
# 	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./lite/*/statik.go" -not -path "*.pb.go" -not -path "*.pulsar.go" | xargs goimports -w -local github.com/irisnet/irismod

# benchmark:
# 	@go test -mod=readonly -bench=. ./...

###############################################################################
###                        Compile Solidity Contracts                       ###
###############################################################################
include scripts/build/contract.mk


